/*
@asiainfo.com
*/

package servicebroker

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/kubernetes-incubator/apiserver-builder/pkg/builders"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"

	prdutil "github.com/asiainfoldp/apiserver-servicebroker/pkg/apis/prd/util"
	"github.com/asiainfoldp/apiserver-servicebroker/pkg/apis/prd/v1"
	clientset "github.com/asiainfoldp/apiserver-servicebroker/pkg/client/clientset_generated/clientset"
	prdclient_v1 "github.com/asiainfoldp/apiserver-servicebroker/pkg/client/clientset_generated/clientset/typed/prd/v1"
	listers "github.com/asiainfoldp/apiserver-servicebroker/pkg/client/listers_generated/prd/v1"
	"github.com/asiainfoldp/apiserver-servicebroker/pkg/controller/sharedinformers"
)

// +controller:group=prd,version=v1,kind=ServiceBroker,resource=servicebrokers
type ServiceBrokerControllerImpl struct {
	builders.DefaultControllerFns

	// lister indexes properties about ServiceBroker
	lister listers.ServiceBrokerLister
	// todo: no need to export
	Clientset clientset.Interface
	Client    prdclient_v1.PrdV1Interface
	// todo: this type is not essential
	ServiceBrokerClient Interface
}

// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *ServiceBrokerControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
	// Use the lister for indexing servicebrokers labels
	c.lister = arguments.GetSharedInformers().Factory.Prd().V1().ServiceBrokers().Lister()

	// ...
	c.Clientset = clientset.NewForConfigOrDie(arguments.GetRestConfig())
	c.Client = c.Clientset.PrdV1() // prdclient_v1.NewForConfigOrDie(config)
	c.ServiceBrokerClient = NewServiceBrokerClient()
}

// Reconcile handles enqueued messages
func (c *ServiceBrokerControllerImpl) Reconcile(sb *v1.ServiceBroker) error {
	// Implement controller logic here
	log.Printf("Running reconcile ServiceBroker for %s\n", sb.Name)

	if sb.Spec.Url == "" {
		return nil
	}

	// The sb is deleted from rest api, but there may be many dependents of the sb.
	// We will do the real deletion only if all the dependents have been deleted.
	if sb.DeletionTimestamp != nil {
		needUpdate := sb.Status.Phase != v1.ServiceBroker_PhaseDeleting
		if needUpdate {
			// todo: disable all related backing service resources
		}

		for finalizers := sb.GetFinalizers(); len(finalizers) > 0; {
			if finalizers[0] != prdutil.DataFoundryFinalizer() {
				break
			}

			// ...
			if numDependents, err := countWorkingBackingServiceInstance(sb.Name, c.Client.BackingServices(), c.Client.BackingServiceInstances(meta_v1.NamespaceAll)); err != nil {
				glog.Errorln("countWorkingBackingServiceInstance err ", err)
				break
			} else if numDependents > 0 {
				break
			}

			// ...
			sb.SetFinalizers(finalizers[1:])
			needUpdate = true
			if len(sb.GetFinalizers()) == 0 {
				sb.DeletionTimestamp = nil // to enter the following switch block next time
			}
			break
		}

		if needUpdate {
			// Will delete the sb in the following swtich block.
			c.Client.ServiceBrokers().Update(sb)
		}

		return nil
	}

	switch sb.Status.Phase {
	case v1.ServiceBroker_PhaseNew:
		if getRetryTime(sb) <= 3 {
			if timeUp(v1.ServiceBroker_PingTimer, sb, 10) {
				setRetryTime(sb)

				services, err := c.ServiceBrokerClient.Catalog(sb.Spec.Url, sb.Spec.UserName, sb.Spec.Password)
				if err != nil {
					c.Client.ServiceBrokers().Update(sb)
					return err
				}

				errs := []error{}
				for _, v := range services.Services {

					if err := backingServiceHandler(c.Client, newBackingService(sb.Name, v)); err != nil {
						errs = append(errs, err)
					}
				}
				if len(errs) == 0 {
					removeRetryTime(sb)
					sb.Status.Phase = v1.ServiceBroker_PhaseActive
				}

				c.Client.ServiceBrokers().Update(sb)
				//return nil
			}
		} else {
			sb.Status.Phase = v1.ServiceBroker_PhaseFailed
			c.Client.ServiceBrokers().Update(sb)

			c.inActiveBackingService(sb.Name)
			//return nil
		}
	case v1.ServiceBroker_PhaseActive:
		if timeUp(v1.ServiceBroker_PingTimer, sb, 60) {
			services, err := c.ServiceBrokerClient.Catalog(sb.Spec.Url, sb.Spec.UserName, sb.Spec.Password)
			if err != nil {
				sb.Status.Phase = v1.ServiceBroker_PhaseFailed
				c.Client.ServiceBrokers().Update(sb)

				c.inActiveBackingService(sb.Name)
				return err
			}

			if timeUp(v1.ServiceBroker_RefreshTimer, sb, 300) {
				for _, v := range services.Services {

					backingServiceHandler(c.Client, newBackingService(sb.Name, v))
				}
			}

			c.Client.ServiceBrokers().Update(sb)
			//return nil
		}
	case v1.ServiceBroker_PhaseFailed:
		if timeUp(v1.ServiceBroker_PingTimer, sb, 60) {
			_, err := c.ServiceBrokerClient.Catalog(sb.Spec.Url, sb.Spec.UserName, sb.Spec.Password)
			if err != nil {
				c.Client.ServiceBrokers().Update(sb)
				return err
			}

			sb.Status.Phase = v1.ServiceBroker_PhaseActive
			c.Client.ServiceBrokers().Update(sb)

			c.ActiveBackingService(sb.Name)
			//return nil
		}

	case v1.ServiceBroker_PhaseDeleting:
		glog.Info("Inavtinging Bs", sb.Name)
		c.inActiveBackingService(sb.Name)
		c.Client.ServiceBrokers().Delete(sb.Name, &meta_v1.DeleteOptions{})
		//return nil

	}

	return nil
}

func (c *ServiceBrokerControllerImpl) recoverBackingService(backingService *v1.BackingService) error {
	_, err := c.Client.BackingServices( /*v1.BackingService_NS*/ ).Get(backingService.Name, meta_v1.GetOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			if _, err := c.Client.BackingServices( /*v1.BackingService_NS*/ ).Create(backingService); err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

func (c *ServiceBrokerControllerImpl) inActiveBackingService(serviceBrokerName string) {
	selector, _ := labels.Parse(v1.ServiceBroker_Label + "=" + serviceBrokerName)

	bsList, err := c.Client.BackingServices( /*v1.BackingService_NS*/ ).List(meta_v1.ListOptions{LabelSelector: selector.String()})
	if err == nil {
		for _, bsvc := range bsList.Items {
			if bsvc.Status.Phase != v1.BackingService_PhaseInactive {
				bsvc.Status.Phase = v1.BackingService_PhaseInactive
				c.Client.BackingServices( /*v1.BackingService_NS*/ ).Update(&bsvc)
			}
		}
	} else {
		glog.Error("can't find bs of sb", serviceBrokerName)
	}
}

func (c *ServiceBrokerControllerImpl) ActiveBackingService(serviceBrokerName string) {
	selector, _ := labels.Parse(v1.ServiceBroker_Label + "=" + serviceBrokerName)

	bsList, err := c.Client.BackingServices( /*v1.BackingService_NS*/ ).List(meta_v1.ListOptions{LabelSelector: selector.String()})
	if err == nil {
		for _, bsvc := range bsList.Items {
			if bsvc.Status.Phase != v1.BackingService_PhaseActive {
				bsvc.Status.Phase = v1.BackingService_PhaseActive
				c.Client.BackingServices( /*v1.BackingService_NS*/ ).Update(&bsvc)
			}
		}
	}
}

func timeUp(timeKind string, sb *v1.ServiceBroker, intervalSecond int64) bool {
	if sb.Annotations == nil {
		sb.Annotations = map[string]string{}
	}

	lastTimeStr := sb.Annotations[timeKind]
	if len(lastTimeStr) == 0 {
		sb.Annotations[timeKind] = fmt.Sprintf("%d", time.Now().UnixNano())
		return true
	}

	lastPing, err := strconv.Atoi(lastTimeStr)
	if err != nil {
		sb.Annotations[timeKind] = fmt.Sprintf("%d", time.Now().UnixNano())
		return false
	}

	if (time.Now().UnixNano()-int64(lastPing))/1e9 < intervalSecond {
		return false
	}

	sb.Annotations[timeKind] = fmt.Sprintf("%d", time.Now().UnixNano())
	return true
}

func getRetryTime(sb *v1.ServiceBroker) int {
	retries := sb.Annotations[v1.ServiceBroker_NewRetryTimes]
	if len(retries) == 0 {
		return 0
	}

	i, err := strconv.Atoi(retries)
	if err != nil {
		return 0
	}

	return i
}

func setRetryTime(sb *v1.ServiceBroker) {
	if len(sb.Annotations[v1.ServiceBroker_NewRetryTimes]) == 0 {
		sb.Annotations[v1.ServiceBroker_NewRetryTimes] = fmt.Sprintf("%d", 0)
	}

	i, err := strconv.Atoi(sb.Annotations[v1.ServiceBroker_NewRetryTimes])
	if err != nil {
		sb.Annotations[v1.ServiceBroker_NewRetryTimes] = fmt.Sprintf("%d", 1)
		return
	}

	sb.Annotations[v1.ServiceBroker_NewRetryTimes] = fmt.Sprintf("%d", i+1)

	return
}

func removeRetryTime(sb *v1.ServiceBroker) {
	delete(sb.Annotations, v1.ServiceBroker_NewRetryTimes)
}

//================================

func newBackingService(name string, spec v1.BackingServiceSpec) *v1.BackingService {
	bs := new(v1.BackingService)
	bs.Spec = v1.BackingServiceSpec(spec)
	bs.Annotations = make(map[string]string)
	bs.Name = spec.Name
	bs.GenerateName = name
	bs.Labels = map[string]string{
		v1.ServiceBroker_Label: name,
	}

	return bs
}

func backingServiceHandler(client prdclient_v1.PrdV1Interface, backingService *v1.BackingService) error {
	newBs, err := client.BackingServices( /*v1.BackingService_NS*/ ).Get(backingService.Name, meta_v1.GetOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			if _, err := client.BackingServices( /*v1.BackingService_NS*/ ).Create(backingService); err != nil {
				glog.Errorln("servicebroker create backingservice err ", err)
				return err
			}
		}
	} else {
		newBs.Spec = backingService.Spec
		newBs.Status.Phase = v1.BackingService_PhaseActive
		if _, err := client.BackingServices( /*v1.BackingService_NS*/ ).Update(newBs); err != nil {
			glog.Errorln("servicebroker update backingservice err ", err)
			return err
		}
	}

	return nil
}

//================================

type ServiceList struct {
	Services []v1.BackingServiceSpec `json:"services"`
}

type Interface interface {
	Catalog(Url string, credential ...string) (ServiceList, error)
}

func NewServiceBrokerClient() Interface {
	return &httpClient{
		Get:  httpGet,
		Post: httpPostJson,
	}
}

type httpClient struct {
	Get  func(getUrl string, credential ...string) ([]byte, error)
	Post func(getUrl string, body []byte, credential ...string) ([]byte, error)
}

func (c *httpClient) Catalog(Url string, credential ...string) (ServiceList, error) {
	services := new(ServiceList)
	b, err := c.Get(Url+"/v2/catalog", credential...)
	if err != nil {
		fmt.Printf("httpclient catalog err %s", err.Error())
		return *services, err
	}

	if err := json.Unmarshal(b, services); err != nil {
		return *services, err
	}

	return *services, nil
}

//todo 支持多种自定义认证方式
func httpGet(getUrl string, credential ...string) ([]byte, error) {
	var resp *http.Response
	var err error
	if len(credential) == 2 {
		tr := &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		req, err := http.NewRequest("GET", getUrl, nil)
		if err != nil {
			return nil, fmt.Errorf("[servicebroker http client] err %s, %s\n", getUrl, err)
		}
		req.Close = true

		basic := fmt.Sprintf("Basic %s", string(base64Encode([]byte(fmt.Sprintf("%s:%s", credential[0], credential[1])))))
		req.Header.Set("Authorization", basic)

		resp, err = client.Do(req)
		if err != nil {
			fmt.Errorf("http get err:%s", err.Error())
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("[servicebroker http client] status err %s, %d\n", getUrl, resp.StatusCode)
		}
	} else {
		resp, err = http.Get(getUrl)
		if err != nil {
			fmt.Errorf("servicebroker http client get err:%s", err.Error())
			return nil, err
		}
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("[http get] status err %s, %d\n", getUrl, resp.StatusCode)
		}
	}

	glog.Infof("GET %s returns http code %v", getUrl, resp.StatusCode)
	return ioutil.ReadAll(resp.Body)
}

func httpPostJson(postUrl string, body []byte, credential ...string) ([]byte, error) {
	var resp *http.Response
	var err error
	req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("[http] err %s, %s\n", postUrl, err)
	}
	req.Header.Set("Content-Type", "application/json")
	if len(credential) == 2 {
		basic := fmt.Sprintf("Basic %s", string(base64Encode([]byte(fmt.Sprintf("%s:%s", credential[0], credential[1])))))
		req.Header.Set("Authorization", basic)
	}
	resp, err = http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("[http] err %s, %s\n", postUrl, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[http] status err %s, %d\n", postUrl, resp.StatusCode)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[http] read err %s, %s\n", postUrl, err)
	}
	return b, nil
}

func base64Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}

//=========================================

func listBackingServiceByServiceBrokerName(name string, bsClient prdclient_v1.BackingServiceInterface) (*v1.BackingServiceList, error) {
	selector, _ := labels.Parse(v1.ServiceBroker_Label + "=" + name)
	return bsClient.List(meta_v1.ListOptions{LabelSelector: selector.String()})
}

func listBackingServiceInstanceByBackingServiceName(name string, bsiClient prdclient_v1.BackingServiceInstanceInterface) (*v1.BackingServiceInstanceList, error) {
	selector, _ := fields.ParseSelector("spec.provisioning.backingservice_name=" + name)
	return bsiClient.List(meta_v1.ListOptions{FieldSelector: selector.String()})
}

func countWorkingBackingServiceInstance(name string, bsClient prdclient_v1.BackingServiceInterface, bsiClient prdclient_v1.BackingServiceInstanceInterface) (int, error) {
	total := 0

	bsList, err := listBackingServiceByServiceBrokerName(name, bsClient)
	if err != nil {
		return total, err
	}

	for _, v := range bsList.Items {
		bsiList, err := listBackingServiceInstanceByBackingServiceName(v.Name, bsiClient)
		if err != nil && kerrors.IsNotFound(err) {
			continue
		}

		total += len(bsiList.Items)
	}

	return total, nil
}

//=========================================

func (c *ServiceBrokerControllerImpl) Get(namespace, name string) (*v1.ServiceBroker, error) {
	return c.lister.Get(name)
}
