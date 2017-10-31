/*
@asiainfo.com
*/

package backingserviceinstance

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/kubernetes-incubator/apiserver-builder/pkg/builders"
	kapi "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	//"k8s.io/client-go/tools/record"

	"github.com/asiainfoldp/apiserver-servicebroker/pkg/apis/prd/v1"
	clientset "github.com/asiainfoldp/apiserver-servicebroker/pkg/client/clientset_generated/clientset"
	prdclient_v1 "github.com/asiainfoldp/apiserver-servicebroker/pkg/client/clientset_generated/clientset/typed/prd/v1"
	listers "github.com/asiainfoldp/apiserver-servicebroker/pkg/client/listers_generated/prd/v1"
	"github.com/asiainfoldp/apiserver-servicebroker/pkg/controller/sharedinformers"
)

// +controller:group=prd,version=v1,kind=BackingServiceInstance,resource=backingserviceinstances
type BackingServiceInstanceControllerImpl struct {
	builders.DefaultControllerFns

	// lister indexes properties about BackingServiceInstance
	lister listers.BackingServiceInstanceLister
	// ...
	//recorder record.EventRecorder
	// todo: no need to export
	Clientset clientset.Interface
	Client    prdclient_v1.PrdV1Interface
}

// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *BackingServiceInstanceControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
	// Use the lister for indexing backingserviceinstances labels
	c.lister = arguments.GetSharedInformers().Factory.Prd().V1().BackingServiceInstances().Lister()

	// ...
	c.Clientset = clientset.NewForConfigOrDie(arguments.GetRestConfig())
	c.Client = c.Clientset.PrdV1() // prdclient_v1.NewForConfigOrDie(config)

	//si.KubernetesFactory
	//si.KubernetesClientSet.Deployments("").Get(...)

	//eventBroadcaster := record.NewBroadcaster()
	//eventBroadcaster.StartRecordingToSink(factory.KubeClient.Events(""))
	//recorder:   eventBroadcaster.NewRecorder(kapi.EventSource{Component: "bsi"}),
}

// Reconcile handles enqueued messages
func (c *BackingServiceInstanceControllerImpl) Reconcile(bsi *v1.BackingServiceInstance) (result error) {
	// Implement controller logic here
	log.Printf("Running reconcile BackingServiceInstance for %s\n", bsi.Name)

	////c.recorder.Eventf(bsi, "Debug", "bsi handler called.%s", bsi.Name)

	changed := false
	bs := &v1.BackingService{}
	if bsi.Annotations[v1.UPS] != "true" &&
		bsi.Status.Phase != v1.BackingServiceInstancePhaseFailure {
		//bsp, err := c.Client.BackingServices("openshift").Get(bsi.Spec.BackingServiceName, meta_v1.GetOptions{})
		bsp, err := c.Client.BackingServices().Get(bsi.Spec.BackingServiceName, meta_v1.GetOptions{})
		if err != nil {
			//c.recorder.Eventf(bsi, kapi.EventTypeWarning, "New", err.Error())
			bsi.Status.Phase = v1.BackingServiceInstancePhaseFailure
			bsi.Status.Reason = err.Error()
			c.Client.BackingServiceInstances(bsi.Namespace).Update(bsi)
			return err
		} else {
			bs = bsp
		}
	}
	if bsi.Annotations == nil {
		bsi.Annotations = make(map[string]string)
	}
	bsi.Annotations["datafoundry.io/servicebroker"] = bs.GenerateName

	switch bsi.Status.Phase {
	case v1.BackingServiceInstancePhaseFailure:
		if bsi.Status.Action == v1.BackingServiceInstanceActionToDelete {
			return c.Client.BackingServiceInstances(bsi.Namespace).Delete(bsi.Name, &meta_v1.DeleteOptions{})
		}
		break
	default:

		result = fmt.Errorf("unknown phase: %s", bsi.Status.Phase)

	case v1.BackingServiceInstancePhaseDeleted:

		glog.Infoln("bsi delete etcd ", bsi.Name)

		result = c.Client.BackingServiceInstances(bsi.Namespace).Delete(bsi.Name, &meta_v1.DeleteOptions{})

	case "":

		bsi.Status.Phase = v1.BackingServiceInstancePhaseProvisioning

		changed = true

		fallthrough

	case v1.BackingServiceInstancePhaseProvisioning:
		if bsi.Status.Action == v1.BackingServiceInstanceActionToDelete {
			return c.Client.BackingServiceInstances(bsi.Namespace).Delete(bsi.Name, &meta_v1.DeleteOptions{})
		}

		glog.Infoln("bsi provisioning ", bsi.Name)
		////c.recorder.Eventf(bsi, "Provisioning", "bsi:%s, service:%s", bsi.Name, bsi.Spec.BackingServiceName)

		plan_found := false
		for _, plan := range bs.Spec.Plans {
			if bsi.Spec.BackingServicePlanGuid == plan.Id {
				bsi.Spec.BackingServicePlanName = plan.Name
				plan_found = true
				break
			}
		}

		if !plan_found {
			result = fmt.Errorf("plan (%s) in bs(%s) for bsi (%s) not found",
				bsi.Spec.BackingServicePlanGuid, bsi.Spec.BackingServiceName, bsi.Name)
			//c.recorder.Eventf(bsi, kapi.EventTypeWarning, "Provisioning", result.Error())
			bsi.Status.Phase = v1.BackingServiceInstancePhaseFailure
			bsi.Status.Reason = result.Error()
			changed = true
			break
		}

		// ...

		glog.Infoln("bsi provisioning servicebroker_load, ", bsi.Name)
		////c.recorder.Eventf(bsi, "Provisioning", "bsi %s provisioning servicebroker_load", bsi.Name)
		bsInstanceID := string(uuid.NewUUID())

		bsi.Spec.InstanceID = bsInstanceID

		servicebroker, err := servicebroker_load(c.Client, bs.GenerateName)
		if err != nil {
			result = err
			break
		}

		serviceinstance := &ServiceInstance{}
		serviceinstance.ServiceId = bs.Spec.Id
		serviceinstance.PlanId = bsi.Spec.BackingServicePlanGuid
		serviceinstance.OrganizationGuid = bsi.Namespace
		serviceinstance.SpaceGuid = bsi.Namespace
		//serviceinstance.Parameters = bsi.Spec.InstanceProvisioning.Parameters
		// serviceinstance.Parameters = make(map[string]string)
		// for k, v := range bsi.Spec.InstanceProvisioning.Parameters {
		// 	serviceinstance.Parameters[k] = v
		// }
		serviceinstance.Parameters = bsi.Spec.InstanceProvisioning.Parameters

		glog.Infoln("bsi provisioning servicebroker_create_instance, ", bsi.Name)

		svcinstance, err := servicebroker_create_instance(serviceinstance, bsInstanceID, servicebroker)
		if err != nil {
			result = err
			//c.recorder.Eventf(bsi, kapi.EventTypeWarning, "Provisioning", err.Error())
			bsi.Status.Phase = v1.BackingServiceInstancePhaseFailure
			bsi.Status.Reason = err.Error()
			changed = true
			break
		} else {
			//c.recorder.Eventf(bsi, kapi.EventTypeNormal, "Provisioning", "bsi provisioning done, instanceid: %s", bsInstanceID)
			glog.Infoln("bsi provisioning servicebroker_create_instance done, ", bsi.Name)
		}

		bsi.Spec.DashboardUrl = svcinstance.DashboardUrl
		bsi.Spec.BackingServiceSpecID = bs.Spec.Id
		bsi.Spec.Creds = svcinstance.Credentials
		if bsi.Spec.Parameters == nil {
			bsi.Spec.Parameters = make(map[string]string)
		}
		bsi.Spec.Parameters["instance_id"] = bsInstanceID

		bsi.Status.Phase = v1.BackingServiceInstancePhaseUnbound

		changed = true

		glog.Infoln("bsi inited. ", bsi.Name)

	case v1.BackingServiceInstancePhaseUnbound:
		switch bsi.Status.Action {
		case v1.BackingServiceInstanceActionToDelete:
			if bsi.Annotations[v1.UPS] == "true" {
				glog.Infoln(v1.UPS, " bsi deleted ", bsi.Name)

				bsi.Status.Phase = v1.BackingServiceInstancePhaseDeleted

				bsi.Status.Action = ""
				changed = true
			} else {
				if result = c.deleteInstance(bs, bsi); result == nil {
					changed = true
				}
			}
			//c.recorder.Eventf(bsi, kapi.EventTypeNormal, "Deleting", "instance:%s [%v]", bsi.Name, changed)
		case v1.BackingServiceInstanceActionToBind:
			if hadoopUser, ok := bsi.Annotations[v1.BindKind_HadoopUser]; !ok {
				dcname := c.get_deploymentconfig_name(bsi, v1.BindDeploymentConfigBinding)
				if bsi.Annotations[v1.UPS] == "true" {
					if result = c.bindInstanceUPS(dcname, bsi); result == nil {
						changed = true
					}
				} else {
					if result = c.bindInstance(dcname, bs, bsi); result == nil {
						changed = true
					}
				}
				//c.recorder.Eventf(bsi, kapi.EventTypeNormal, "Binding", "instance: %s, dc: %s [%v]", bsi.Name, dcname, changed)
			} else {
				if result = c.bindInstanceHadoop(hadoopUser, bs, bsi); result == nil {
					changed = true
				}
				//c.recorder.Eventf(bsi, kapi.EventTypeNormal, "Binding", "instance: %s, hadoopuser: %s [%v]", bsi.Name, hadoopUser, changed)
			}
		}
	case v1.BackingServiceInstancePhaseBound:
		switch bsi.Status.Action {
		case v1.BackingServiceInstanceActionToUnbind:
			if hadoopUser, ok := bsi.Annotations[v1.BindKind_HadoopUser]; !ok {
				dcname := c.get_deploymentconfig_name(bsi, v1.BindDeploymentConfigUnbinding)
				if bsi.Annotations[v1.UPS] == "true" {
					if result = c.unbindInstanceUPS(dcname, bsi); result == nil {
						changed = true
					}
				} else {
					if result = c.unbindInstance(dcname, bs, bsi); result == nil {
						changed = true
					}
				}
				//c.recorder.Eventf(bsi, kapi.EventTypeNormal, "Unbinding", "instance: %s, dc: %s [%v]", bsi.Name, dcname, changed)
			} else {
				if result = c.unbindInstanceHadoop(hadoopUser, bs, bsi); result == nil {
					changed = true
				}
				//c.recorder.Eventf(bsi, kapi.EventTypeNormal, "Binding", "instance: %s, hadoopuser: %s [%v]", bsi.Name, hadoopUser, changed)
			}
		case v1.BackingServiceInstanceActionToBind:
			if hadoopUser, ok := bsi.Annotations[v1.BindKind_HadoopUser]; !ok {
				dcname := c.get_deploymentconfig_name(bsi, v1.BindDeploymentConfigBinding)
				if bsi.Annotations[v1.UPS] == "true" {
					if result = c.bindInstanceUPS(dcname, bsi); result == nil {
						changed = true
					}
				} else {
					if result = c.bindInstance(dcname, bs, bsi); result == nil {
						changed = true
					}
				}
				//c.recorder.Eventf(bsi, kapi.EventTypeNormal, "Binding", "instance: %s, dc: %s [%v]", bsi.Name, dcname, changed)
			} else {
				if result = c.bindInstanceHadoop(hadoopUser, bs, bsi); result == nil {
					changed = true
				}
				//c.recorder.Eventf(bsi, kapi.EventTypeNormal, "Binding", "instance: %s, hadoopuser: %s [%v]", bsi.Name, hadoopUser, changed)
			}
		default:
			c.check_dc_healthy(bsi)

		}

	}

	//patch
	if bsi.Status.Phase == v1.BackingServiceInstancePhaseBound ||
		bsi.Status.Phase == v1.BackingServiceInstancePhaseUnbound {
		switch bsi.Status.Patch {
		case v1.BackingServiceInstancePatchUpdate:
			// do patch api./
			// bsi.Status.Patch = v1.BackingServiceInstancePatchUpdating
			if result = c.updateInstance(bs, bsi); result == nil {
				changed = true
				bsi.Status.Patch = v1.BackingServiceInstancePatchUpdated
				// bsi.Spec.Accesses = nil
			} else {
				bsi.Status.Patch = v1.BackingServiceInstancePatchFailure
				changed = true
			}
		case v1.BackingServiceInstancePatchUpdated:
			bsi.Status.Patch = ""
			changed = true
		// do remove patch phase.
		case v1.BackingServiceInstancePatchUpdating:
		case v1.BackingServiceInstancePatchFailure:
		//prevent updating.
		default:
			// glog.Info(bsi.Name, " no update due to unknown patch phase.", bsi.Status.Patch)
			bsi.Status.Patch = ""
			changed = true
		}

	}

	if result != nil {

		err_msg := result.Error()
		/*
			if err_msg != bsi.Status.Error {
				glog.Info("#:", err_msg, "#:", bsi.Status.Error)
				changed = true
				bsi.Status.Error = err_msg
			}
		*/

		glog.Infoln("bsi controller error. ", err_msg)
		//c.recorder.Eventf(bsi, kapi.EventTypeWarning, "Error", err_msg)
	}

	if changed {
		// glog.Infoln("bsi etc changed and update. ", bsi.Name)

		c.Client.BackingServiceInstances(bsi.Namespace).Update(bsi)
	}

	return //nil
}

//=====================================
//
//=====================================

type fatalError string

func (e fatalError) Error() string {
	return "fatal error handling BackingServiceInstanceController: " + string(e)
}

//=====================================
//
//=====================================

func has_action_word(text, word v1.BackingServiceInstanceAction) bool {
	return strings.Index(string(text), string(word)) >= 0
}

//func remove_action_word(text, word v1.BackingServiceInstanceAction) v1.BackingServiceInstanceAction {
func remove_action_word(text, word string) string {
	for {
		index := strings.Index(string(text), string(word))
		if index >= 0 {
			text = text[:index] + text[index+len(word):]
			continue
		}
		break
	}

	return text
}

func servicebroker_load(c prdclient_v1.PrdV1Interface /*osclient.Interface*/, name string) (*ServiceBroker, error) {
	servicebroker := &ServiceBroker{}
	if sb, err := c.ServiceBrokers().Get(name, meta_v1.GetOptions{}); err != nil {
		return nil, err
	} else {
		servicebroker.Url = sb.Spec.Url
		servicebroker.UserName = sb.Spec.UserName
		servicebroker.Password = sb.Spec.Password
		servicebroker.Name = name

		return servicebroker, nil
	}
}

func checkIfPlanidExist(client prdclient_v1.PrdV1Interface /*osclient.Interface*/, planId string) (bool, *v1.BackingService, error) {

	//items, err := client.BackingServices("openshift").List(meta_v1.ListOptions{})
	items, err := client.BackingServices().List(meta_v1.ListOptions{})

	if err != nil {
		return false, nil, err
	}

	for _, bs := range items.Items {
		for _, plans := range bs.Spec.Plans {
			if planId == plans.Id {
				glog.Info("we found plan id at plan", bs.Spec.Name)

				return true, &bs, nil
			}
		}
	}
	return false, nil, fatalError(fmt.Sprintf("Can't find plan id %s", planId))

}

func commToServiceBroker(method, path string, jsonData []byte, header map[string]string) (resp *http.Response, err error) {

	glog.Infoln(method, path, string(jsonData))

	tr := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{

		Transport: tr,
		Timeout:   time.Second * 120,
	}

	req, err := http.NewRequest(strings.ToUpper(method) /*SERVICE_BROKER_API_SERVER+*/, path, bytes.NewBuffer(jsonData))

	if len(header) > 0 {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}
	return client.Do(req)
}

type ServiceBroker struct {
	Url      string
	UserName string
	Password string
	Name     string
}

type ServiceInstance struct {
	ServiceId        string `json:"service_id"`
	PlanId           string `json:"plan_id"`
	OrganizationGuid string `json:"organization_guid"`
	SpaceGuid        string `json:"space_guid"`
	//Incomplete       bool        `json:"accepts_incomplete, omitempty"`
	//Parameters interface{} `json:"parameters, omitempty"`
	Parameters map[string]string `json:"parameters, omitempty"`
}
type LastOperation struct {
	State                    string `json:"state"`
	Description              string `json:"description"`
	AsyncPollIntervalSeconds int    `json:"async_poll_interval_seconds, omitempty"`
}
type CreateServiceInstanceResponse struct {
	DashboardUrl  string            `json:"dashboard_url"`
	LastOperation *LastOperation    `json:"last_operation, omitempty"`
	Credentials   map[string]string `json:"credentials, omitempty"`
}

type ServiceBinding struct {
	ServiceId       string            `json:"service_id"`
	PlanId          string            `json:"plan_id"`
	AppGuid         string            `json:"app_guid,omitempty"`
	BindResource    map[string]string `json:"bind_resource,omitempty"`
	Parameters      map[string]string `json:"parameters,omitempty"`
	svc_instance_id string
}

type ServiceBindingResponse struct {
	Credentials     map[string]string `json:"credentials"`
	SyslogDrainUrl  string            `json:"syslog_drain_url"`
	RouteServiceUrl string            `json:"route_service_url"`
}

type Credential struct {
	Uri      string `json:"uri"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Vhost    string `json:"vhost"`
	//Database string `json:"database"`
}

func servicebroker_create_instance(param *ServiceInstance, instance_guid string, sb *ServiceBroker) (*CreateServiceInstanceResponse, error) {
	jsonData, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Authorization"] = basicAuthStr(sb.UserName, sb.Password)

	// resp, err := commToServiceBroker("PUT", sb.Url+"/v2/service_instances/"+instance_guid+"?accepts_incomplete=true", jsonData, header)
	resp, err := commToServiceBroker("PUT", sb.Url+"/v2/service_instances/"+instance_guid, jsonData, header)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	glog.Infof("respcode from /v2/service_instances/%s: %v", instance_guid, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	svcinstance := &CreateServiceInstanceResponse{}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated ||
		resp.StatusCode == http.StatusAccepted {

		if len(body) > 0 {
			err = json.Unmarshal(body, svcinstance)

			if err != nil {
				glog.Error(err)
				return nil, err
			}
		}
	} else {
		glog.Error("Error:", string(body))
		return nil, fmt.Errorf("%d returned from broker %s: %s", resp.StatusCode, sb.Name, string(body))
	}
	glog.Infof("%v,%+v\n", string(body), svcinstance)
	return svcinstance, nil
}

func servicebroker_update_instance(bsi *v1.BackingServiceInstance, sb *ServiceBroker) (interface{}, error) {

	serviceinstance := &ServiceInstance{}
	serviceinstance.ServiceId = bsi.Spec.BackingServiceSpecID
	serviceinstance.PlanId = bsi.Spec.BackingServicePlanGuid
	// serviceinstance.OrganizationGuid = bsi.Namespace
	// serviceinstance.SpaceGuid = bsi.Namespace
	serviceinstance.Parameters = bsi.Spec.InstanceProvisioning.Parameters
	// serviceinstance.Parameters = make(map[string]interface{})
	// for k, v := range bsi.Spec.InstanceProvisioning.Parameters {
	// 	serviceinstance.Parameters[k] = v
	// }
	// serviceinstance.Parameters["accesses"] = bsi.Spec.Accesses

	jsonBody, err := json.Marshal(serviceinstance)
	if err != nil {
		return nil, err
	}

	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Authorization"] = basicAuthStr(sb.UserName, sb.Password)

	resp, err := commToServiceBroker("PATCH", sb.Url+"/v2/service_instances/"+bsi.Spec.InstanceID, jsonBody, header)
	if err != nil {

		glog.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	glog.Infof("respcode from PATCH /v2/service_instances/%s: %v", bsi.Spec.InstanceID, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	type patchResp struct {
		Response interface{}
	}
	svcPatch := &patchResp{}

	if resp.StatusCode == http.StatusOK {
		if len(body) > 0 {
			err = json.Unmarshal(body, svcPatch)

			if err != nil {
				glog.Error(err)
				return nil, err
			}
		}
	}
	glog.Infof("%v,%+v\n", string(body), svcPatch)

	return svcPatch, nil

}
func servicebroker_binding(param *ServiceBinding, binding_guid string, sb *ServiceBroker) (*ServiceBindingResponse, error) {
	jsonData, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Authorization"] = basicAuthStr(sb.UserName, sb.Password)

	resp, err := commToServiceBroker("PUT", sb.Url+"/v2/service_instances/"+param.svc_instance_id+"/service_bindings/"+binding_guid, jsonData, header)
	if err != nil {

		glog.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	glog.Infof("respcode from PUT /v2/service_instances/%s/service_bindings/%s: %v", param.svc_instance_id, binding_guid, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	svcBinding := &ServiceBindingResponse{}

	glog.Infof("%v,%+v\n", string(body), svcBinding)
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		if len(body) > 0 {
			err = json.Unmarshal(body, svcBinding)

			if err != nil {
				glog.Error(err)
				return nil, err
			}
		}
	}

	return svcBinding, nil
}

func servicebroker_unbinding(bindId string, svc *ServiceBinding, sb *ServiceBroker) (interface{}, error) {

	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Authorization"] = basicAuthStr(sb.UserName, sb.Password)

	resp, err := commToServiceBroker("DELETE", sb.Url+"/v2/service_instances/"+svc.svc_instance_id+"/service_bindings/"+bindId+"?service_id="+svc.ServiceId+"&plan_id="+svc.PlanId, nil, header)
	if err != nil {

		glog.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	glog.Infof("respcode from DELETE /v2/service_instances/%s/service_bindings/%s: %v", svc.svc_instance_id, bindId, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	type UnBindindResp struct {
		Response interface{}
	}
	svcUnbinding := &UnBindindResp{}

	if resp.StatusCode == http.StatusOK {
		if len(body) > 0 {
			err = json.Unmarshal(body, svcUnbinding)

			if err != nil {
				glog.Error(err)
				return nil, err
			}
		}
	}
	glog.Infof("%v,%+v\n", string(body), svcUnbinding)
	return svcUnbinding, nil
}

func servicebroker_deprovisioning(bsi *v1.BackingServiceInstance, sb *ServiceBroker) (interface{}, error) {

	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Authorization"] = basicAuthStr(sb.UserName, sb.Password)

	resp, err := commToServiceBroker("DELETE", sb.Url+"/v2/service_instances/"+bsi.Spec.InstanceID+"?service_id="+bsi.Spec.BackingServiceSpecID+"&plan_id="+bsi.Spec.BackingServicePlanGuid, nil, header)
	if err != nil {

		glog.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	glog.Infof("respcode from DELETE /v2/service_instances/%s: %v", bsi.Spec.InstanceID, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	type DeprovisioningResp struct {
		Response interface{}
	}
	svcDeprovisioning := &DeprovisioningResp{}

	if resp.StatusCode == http.StatusOK {
		if len(body) > 0 {
			err = json.Unmarshal(body, svcDeprovisioning)

			if err != nil {
				glog.Error(err)
				return nil, err
			}
		}
	}
	glog.Infof("%v,%+v\n", string(body), svcDeprovisioning)
	return svcDeprovisioning, nil
}

func basicAuthStr(username, password string) string {
	auth := username + ":" + password
	authstr := base64.StdEncoding.EncodeToString([]byte(auth))
	return "Basic " + authstr
}

var InvalidCharFinder = regexp.MustCompile("[^a-zA-Z0-9]")

func deploymentconfig_env_prefix(bsName, bsiName string) string {
	return strings.ToUpper(fmt.Sprintf("BSI_%s_%s_", InvalidCharFinder.ReplaceAllLiteralString(bsName, ""), InvalidCharFinder.ReplaceAllLiteralString(bsiName, "")))
}

func deploymentconfig_env_name(prefix string, envName string) string {
	return strings.ToUpper(fmt.Sprintf("%s%s", prefix, InvalidCharFinder.ReplaceAllLiteralString(envName, "_")))
}

func (c *BackingServiceInstanceControllerImpl) deploymentconfig_inject_envs(dc string, bsi *v1.BackingServiceInstance, b *v1.InstanceBinding) error {
	return c.deploymentconfig_modify_envs(dc, bsi, b, true)
}

func (c *BackingServiceInstanceControllerImpl) deploymentconfig_clear_envs(dc string, bsi *v1.BackingServiceInstance, b *v1.InstanceBinding) error {
	return c.deploymentconfig_modify_envs(dc, bsi, b, false)
}

func (c *BackingServiceInstanceControllerImpl) deploymentconfig_inject_envs_ups(dc string, bsi *v1.BackingServiceInstance, b *v1.InstanceBinding) error {
	return c.deploymentconfig_modify_envs_ups(dc, bsi, b, true)
}
func (c *BackingServiceInstanceControllerImpl) deploymentconfig_clear_envs_ups(dc string, bsi *v1.BackingServiceInstance, b *v1.InstanceBinding) error {
	return c.deploymentconfig_modify_envs_ups(dc, bsi, b, false)
}

func (c *BackingServiceInstanceControllerImpl) check_dc_healthy(bsi *v1.BackingServiceInstance) {
	if bsi.Spec.Bound < 1 {
		return
	}
	for _, binding := range bsi.Spec.Binding {
		if len(binding.BindDeploymentConfig) == 0 {
			// glog.Info("not a dc.")
			continue
		}
		glog.Infof("check dc '%s' healthy. (disabled now)")
		/*
			glog.Infof("check dc '%s' healthy.", binding.BindDeploymentConfig)
			dc, err := c.Client.DeploymentConfigs(bsi.Namespace).Get(binding.BindDeploymentConfig, meta_v1.GetOptions{})
			if err != nil {
				if kerrors.IsNotFound(err) {
					glog.Infof("dc %s is not found.", binding.BindDeploymentConfig)
				} else {
					glog.Error(err.Error())
				}
			} else {
				if dc.Annotations["backingservice.instance/"+bsi.Name] != "bound" {
					glog.Infof("rebind envs in to dc %s", binding.BindDeploymentConfig)
					if bsi.Annotations[v1.UPS] == "true" {
						err = c.deploymentconfig_inject_envs_ups(binding.BindDeploymentConfig, bsi, &binding)
						if err != nil {
							glog.Error(err.Error())
						}
					} else {
						err = c.deploymentconfig_inject_envs(binding.BindDeploymentConfig, bsi, &binding)
						if err != nil {
							glog.Error(err.Error())
						}
					}
				}
			}
		*/
	}
	return
}

// return exists or not
func env_get(envs []kapi.EnvVar, envName string) (bool, string) {
	for i := len(envs) - 1; i >= 0; i-- {
		if envs[i].Name == envName {
			return true, envs[i].Value
		}
	}

	return false, ""
}

// return overritten or not
func env_set(envs []kapi.EnvVar, envName, envValue string) (bool, []kapi.EnvVar) {
	if envs == nil {
		envs = []kapi.EnvVar{}
	}

	for i := len(envs) - 1; i >= 0; i-- {
		if envs[i].Name == envName {
			envs[i] = kapi.EnvVar{Name: envName, Value: envValue}
			return true, envs
		}
	}

	envs = append(envs, kapi.EnvVar{Name: envName, Value: envValue})
	return false, envs
}

// return unset or not
func env_unset(envs []kapi.EnvVar, envName string) (bool, []kapi.EnvVar) {
	if envs == nil {
		return false, nil
	}

	n := len(envs)
	index := 0
	for i := 0; i < n; i++ {
		if envs[i].Name != envName {
			if index < i {
				envs[index] = envs[i]
			}
			index++
		}
	}

	return index < n, envs[:index]
}

func (c *BackingServiceInstanceControllerImpl) deploymentconfig_modify_envs_ups(dcname string, bsi *v1.BackingServiceInstance, binding *v1.InstanceBinding, toInject bool) error {
	/*
		dc, err := c.Client.DeploymentConfigs(bsi.Namespace).Get(dcname, meta_v1.GetOptions{})
		if err != nil {
			return err
		}

		if dc.Spec.Template == nil {
			return nil
		}

		env_prefix := deploymentconfig_env_prefix(bsi.Spec.BackingServiceName, bsi.Name)
		containers := dc.Spec.Template.Spec.Containers
		if toInject {

			vsp := &VcapServiceParameters{
				Name:        bsi.Name,
				Label:       "",
				Plan:        v1.UPS,
				Credentials: binding.Credentials,
			}

			for i := range containers {
				for k, v := range binding.Credentials {
					_, containers[i].Env = env_set(containers[i].Env, deploymentconfig_env_name(env_prefix, k), v)
				}

				if vsp != nil {
					_, containers[i].Env = modifyVcapServicesEnvNameEnv(containers[i].Env, v1.UPS, vsp, "")
				}
			}
			if dc.Annotations == nil {
				dc.Annotations = make(map[string]string)
			}
			dc.Annotations["backingservice.instance/"+bsi.Name] = "bound"
		} else {
			for i := range containers {
				for k := range binding.Credentials {
					_, containers[i].Env = env_unset(containers[i].Env, deploymentconfig_env_name(env_prefix, k))
				}

				_, containers[i].Env = modifyVcapServicesEnvNameEnv(containers[i].Env, v1.UPS, nil, bsi.Name)
			}
			delete(dc.Annotations, "backingservice.instance/"+bsi.Name)
		}

		if _, err := c.Client.DeploymentConfigs(bsi.Namespace).Update(dc); err != nil {
			return err
		}

		c.deploymentconfig_print_envs(bsi.Namespace, binding)
	*/
	return nil
}

func (c *BackingServiceInstanceControllerImpl) deploymentconfig_modify_envs(dcname string, bsi *v1.BackingServiceInstance, binding *v1.InstanceBinding, toInject bool) error {
	/*
		dc, err := c.Client.DeploymentConfigs(bsi.Namespace).Get(dcname, meta_v1.GetOptions{})
		if err != nil {
			return err
		}

		if dc.Spec.Template == nil {
			return nil
		}

		env_prefix := deploymentconfig_env_prefix(bsi.Spec.BackingServiceName, bsi.Name)
		containers := dc.Spec.Template.Spec.Containers

		if toInject {
			//bs, err := c.Client.BackingServices("openshift").Get(bsi.Spec.BackingServiceName, meta_v1.GetOptions{})
			bs, err := c.Client.BackingServices().Get(bsi.Spec.BackingServiceName, meta_v1.GetOptions{})
			if err != nil {
				return err
			}

			var plan = (*v1.ServicePlan)(nil)
			for k := range bs.Spec.Plans {
				if bsi.Spec.BackingServicePlanGuid == bs.Spec.Plans[k].Id {
					plan = &(bs.Spec.Plans[k])
				}
			}

			var vsp *VcapServiceParameters = nil
			if plan != nil {
				vsp = &VcapServiceParameters{
					Name:        bsi.Name,
					Label:       "",
					Plan:        plan.Name,
					Credentials: binding.Credentials,
				}
			}

			for i := range containers {
				for k, v := range binding.Credentials {
					_, containers[i].Env = env_set(containers[i].Env, deploymentconfig_env_name(env_prefix, k), v)
				}

				if vsp != nil {
					_, containers[i].Env = modifyVcapServicesEnvNameEnv(containers[i].Env, bs.Name, vsp, "")
				}
			}
			if dc.Annotations == nil {
				dc.Annotations = make(map[string]string)
			}
			dc.Annotations["backingservice.instance/"+bsi.Name] = "bound"
		} else {
			for i := range containers {
				for k := range binding.Credentials {
					_, containers[i].Env = env_unset(containers[i].Env, deploymentconfig_env_name(env_prefix, k))
				}

				_, containers[i].Env = modifyVcapServicesEnvNameEnv(containers[i].Env, bsi.Spec.BackingServiceName, nil, bsi.Name)
			}
			delete(dc.Annotations, "backingservice.instance/"+bsi.Name)
		}

		if _, err := c.Client.DeploymentConfigs(bsi.Namespace).Update(dc); err != nil {
			return err
		}

		c.deploymentconfig_print_envs(bsi.Namespace, binding)
	*/

	return nil
}

func modifyVcapServicesEnvNameEnv(env []kapi.EnvVar, bsName string, vsp *VcapServiceParameters, bsiName string) (bool, []kapi.EnvVar) {
	_, json_env := env_get(env, VcapServicesEnvName)

	vs := VcapServices{}
	if len(strings.TrimSpace(json_env)) > 0 {
		err := json.Unmarshal([]byte(json_env), &vs)
		if err != nil {
			glog.Warningln("unmarshalVcapServices error: ", err.Error())
		}
	}

	if vsp != nil {
		vs = addVcapServiceParameters(vs, bsName, vsp)
	}
	if bsiName != "" {
		vs = removeVcapServiceParameters(vs, bsName, bsiName)
	}

	if len(vs) == 0 {
		return env_unset(env, VcapServicesEnvName)
	}
	json_data, err := json.Marshal(vs)
	if err != nil {
		glog.Warningln("marshalVcapServices error: ", err.Error())
		return false, env
	}

	json_env = string(json_data)

	glog.Info("new ", VcapServicesEnvName, " = ", json_env)

	return env_set(env, VcapServicesEnvName, json_env)
}

const VcapServicesEnvName = "VCAP_SERVICES"

type VcapServices map[string][]*VcapServiceParameters

type VcapServiceParameters struct {
	Name        string            `json:"name, omitempty"`
	Label       string            `json:"label, omitempty"`
	Plan        string            `json:"plan, omitempty"`
	Credentials map[string]string `json:"credentials, omitempty"`
}

func addVcapServiceParameters(vs VcapServices, serviceName string, vsParameters *VcapServiceParameters) VcapServices {
	if vs == nil {
		vs = VcapServices{}
	}

	if vsParameters == nil {
		return vs
	}

	removeVcapServiceParameters(vs, serviceName, vsParameters.Name)

	vsp_list := vs[serviceName]
	if vsp_list == nil {
		vsp_list = []*VcapServiceParameters{}
	}
	vsp_list = append(vsp_list, vsParameters)
	vs[serviceName] = vsp_list

	return vs
}

func removeVcapServiceParameters(vs VcapServices, serviceName, instanceName string) VcapServices {
	if vs == nil {
		vs = VcapServices{}
	}

	vsp_list := vs[serviceName]
	if len(vsp_list) == 0 {
		return vs
	}
	num := len(vsp_list)
	vsp_list2 := make([]*VcapServiceParameters, 0, num)
	for i := 0; i < num; i++ {
		vsp := vsp_list[i]
		if vsp != nil && vsp.Name != instanceName {
			vsp_list2 = append(vsp_list2, vsp)
		}
	}
	if len(vsp_list2) == 0 {
		delete(vs, serviceName)
	} else {
		vs[serviceName] = vsp_list2
	}

	return vs
}

/*
func (c *BackingServiceInstanceControllerImpl) deploymentconfig_print_envs(ns string, binding *v1.InstanceBinding) {
	dc, err := c.Client.DeploymentConfigs(ns).Get(binding.BindDeploymentConfig, meta_v1.GetOptions{})
	if err != nil {
		fmt.Println("dc not found: ", binding.BindDeploymentConfig)
		return
	}

	if dc.Spec.Template == nil {
		fmt.Println("dc.Spec.Template is nil")
		return
	}

	containers := dc.Spec.Template.Spec.Containers

	for _, c := range containers {
		fmt.Println("**********  envs in container")

		for _, env := range c.Env {
			fmt.Println("     env[", env.Name, ",] = ", env.Value)
		}
	}
}
*/

func (c *BackingServiceInstanceControllerImpl) deleteInstance(bs *v1.BackingService, bsi *v1.BackingServiceInstance) (result error) {
	glog.Infoln("bsi to delete ", bsi.Name)

	servicebroker, err := servicebroker_load(c.Client, bs.GenerateName)
	if err != nil {
		return err

	}

	glog.Infoln("deleting ", bsi.Name)
	if _, err = servicebroker_deprovisioning(bsi, servicebroker); err != nil {
		return err

	}

	glog.Infoln("bsi deleted ", bsi.Name)

	bsi.Status.Phase = v1.BackingServiceInstancePhaseDeleted

	bsi.Status.Action = remove_action_word(bsi.Status.Action, v1.BackingServiceInstanceActionToDelete)
	return

}

func (c *BackingServiceInstanceControllerImpl) bindInstanceUPS(dc string, bsi *v1.BackingServiceInstance) (err error) {
	glog.Infoln(v1.UPS, "bsi to bind ", bsi.Name, " and ", dc)

	instanceBinding := v1.InstanceBinding{}
	now := meta_v1.Now()             // unversioned.Now()
	instanceBinding.BoundTime = &now //&unversioned.Now()
	instanceBinding.BindUuid = v1.UPS
	instanceBinding.BindDeploymentConfig = dc
	instanceBinding.Credentials = bsi.Spec.Credentials

	/*
		instanceBinding.Credentials = make(map[string]string)
		for k, v := range bsi.Spec.Credentials{
			instanceBinding.Credentials[k] = v
		}
	*/

	glog.Infoln("deploymentconfig_inject_envs")

	err = c.deploymentconfig_inject_envs_ups(dc, bsi, &instanceBinding)
	if err != nil {
		return err
	} else {
		bsi.Spec.Binding = append(bsi.Spec.Binding, instanceBinding)
	}

	glog.Infoln("bsi bound. ", bsi.Name)

	bsi.Spec.Bound += 1

	bsi.Status.Phase = v1.BackingServiceInstancePhaseBound
	bsi.Annotations[dc] = v1.BindDeploymentConfigBound

	bsi.Status.Action = "" //remove_action_word(bsi.Status.Action, v1.BackingServiceInstanceActionToBind)
	return

}

func (c *BackingServiceInstanceControllerImpl) bindInstanceHadoop(user string, bs *v1.BackingService, bsi *v1.BackingServiceInstance) (result error) {
	glog.Infoln("bsi to bind hadoop user ", bsi.Name, " and ", user)

	servicebroker, err := servicebroker_load(c.Client, bs.GenerateName)
	if err != nil {
		return err
	}

	bind_uuid := string(uuid.NewUUID())

	servicebinding := &ServiceBinding{
		ServiceId: bs.Spec.Id,
		PlanId:    bsi.Spec.BackingServicePlanGuid,
		AppGuid:   bsi.Namespace,
		//BindResource: ,
		Parameters:      bsi.Spec.Parameters,
		svc_instance_id: bsi.Spec.InstanceID,
	}

	servicebinding.Parameters["user_name"] = user

	glog.Infoln("bsi to bind", bsi.Name)

	bindingresponse, err := servicebroker_binding(servicebinding, bind_uuid, servicebroker)
	if err != nil {
		return err
	}

	instanceBinding := v1.InstanceBinding{}
	now := meta_v1.Now()             // unversioned.Now()
	instanceBinding.BoundTime = &now //&unversioned.Now()
	instanceBinding.BindUuid = bind_uuid
	instanceBinding.BindHadoopUser = user
	instanceBinding.Credentials = bindingresponse.Credentials
	// instanceBinding.Credentials = make(map[string]string)
	// instanceBinding.Credentials["Uri"] = bindingresponse.Credentials.Uri
	// instanceBinding.Credentials["Name"] = bindingresponse.Credentials.Name
	// instanceBinding.Credentials["Username"] = bindingresponse.Credentials.Username
	// instanceBinding.Credentials["Password"] = bindingresponse.Credentials.Password
	// instanceBinding.Credentials["Host"] = bindingresponse.Credentials.Host
	// instanceBinding.Credentials["Port"] = bindingresponse.Credentials.Port
	// instanceBinding.Credentials["Vhost"] = bindingresponse.Credentials.Vhost

	bsi.Spec.Binding = append(bsi.Spec.Binding, instanceBinding)

	glog.Infoln("bsi bound. ", bsi.Name)

	bsi.Spec.Bound += 1

	bsi.Status.Phase = v1.BackingServiceInstancePhaseBound
	delete(bsi.Annotations, v1.BindKind_HadoopUser)

	bsi.Status.Action = "" //remove_action_word(bsi.Status.Action, v1.BackingServiceInstanceActionToBind)
	return

}
func (c *BackingServiceInstanceControllerImpl) bindInstance(dc string, bs *v1.BackingService, bsi *v1.BackingServiceInstance) (result error) {
	glog.Infoln("bsi to bind ", bsi.Name, " and ", dc)

	servicebroker, err := servicebroker_load(c.Client, bs.GenerateName)
	if err != nil {
		return err
	}

	bind_uuid := string(uuid.NewUUID())

	servicebinding := &ServiceBinding{
		ServiceId: bs.Spec.Id,
		PlanId:    bsi.Spec.BackingServicePlanGuid,
		AppGuid:   bsi.Namespace,
		//BindResource: ,
		Parameters:      bsi.Spec.Parameters,
		svc_instance_id: bsi.Spec.InstanceID,
	}

	glog.Infoln("bsi to bind", bsi.Name)

	bindingresponse, err := servicebroker_binding(servicebinding, bind_uuid, servicebroker)
	if err != nil {
		return err
	}

	instanceBinding := v1.InstanceBinding{}
	now := meta_v1.Now()             // unversioned.Now()
	instanceBinding.BoundTime = &now //&unversioned.Now()
	instanceBinding.BindUuid = bind_uuid
	instanceBinding.BindDeploymentConfig = dc
	instanceBinding.Credentials = bindingresponse.Credentials
	// instanceBinding.Credentials = make(map[string]string)
	// instanceBinding.Credentials["Uri"] = bindingresponse.Credentials.Uri
	// instanceBinding.Credentials["Name"] = bindingresponse.Credentials.Name
	// instanceBinding.Credentials["Username"] = bindingresponse.Credentials.Username
	// instanceBinding.Credentials["Password"] = bindingresponse.Credentials.Password
	// instanceBinding.Credentials["Host"] = bindingresponse.Credentials.Host
	// instanceBinding.Credentials["Port"] = bindingresponse.Credentials.Port
	// instanceBinding.Credentials["Vhost"] = bindingresponse.Credentials.Vhost
	// = bindingresponse.SyslogDrainUrl
	// = bindingresponse.RouteServiceUrl

	glog.Infoln("deploymentconfig_inject_envs")

	err = c.deploymentconfig_inject_envs(dc, bsi, &instanceBinding)
	if err != nil {
		return err
	} else {
		bsi.Spec.Binding = append(bsi.Spec.Binding, instanceBinding)
	}

	glog.Infoln("bsi bound. ", bsi.Name)

	bsi.Spec.Bound += 1

	bsi.Status.Phase = v1.BackingServiceInstancePhaseBound
	bsi.Annotations[dc] = v1.BindDeploymentConfigBound

	bsi.Status.Action = "" //remove_action_word(bsi.Status.Action, v1.BackingServiceInstanceActionToBind)
	return

}

func (c *BackingServiceInstanceControllerImpl) unbindInstanceUPS(dc string, bsi *v1.BackingServiceInstance) (err error) {

	glog.Infoln(v1.UPS, " bsi to unbind ", bsi.Name)

	for idx, b := range bsi.Spec.Binding {
		if b.BindDeploymentConfig == dc {

			glog.Infoln("deploymentconfig_clear_envs")
			err = c.deploymentconfig_clear_envs_ups(dc, bsi, &b)
			if err != nil && (!kerrors.IsNotFound(err)) {
				return err
			} else {
				bsi.Spec.Binding = append(bsi.Spec.Binding[:idx], bsi.Spec.Binding[idx+1:]...)
				delete(bsi.Annotations, dc)
			}

			/*
				err = c.deploymentconfig_clear_envs(bsi.Namespace, dc)
				if err != nil {
					return err
				}
			*/
			glog.Infoln("bsi is unbound ", bsi.Name)
			/*
				delete(bsi.Annotations, dc)
				bsi.Spec.BindDeploymentConfig = ""
				bsi.Spec.Credentials = nil
				bsi.Spec.BoundTime = nil
				bsi.Spec.BindUuid = ""
				bsi.Spec.Bound = false
			*/
			bsi.Spec.Bound -= 1
			if bsi.Spec.Bound == 0 {
				bsi.Status.Phase = v1.BackingServiceInstancePhaseUnbound
			}

			break
		}
	}

	bsi.Status.Action = "" //remove_action_word(bsi.Status.Action, v1.BackingServiceInstanceActionToUnbind)

	return

}

func (c *BackingServiceInstanceControllerImpl) updateInstance(bs *v1.BackingService, bsi *v1.BackingServiceInstance) (result error) {
	plan_found := false
	for _, plan := range bs.Spec.Plans {
		if bsi.Spec.BackingServicePlanGuid == plan.Id {
			bsi.Spec.BackingServicePlanName = plan.Name
			plan_found = true
			break
		}
	}

	if !plan_found {
		//c.recorder.Eventf(bsi, kapi.EventTypeNormal, "Updating", "plan (%s) in bs(%s) for bsi (%s) not found",
		//	bsi.Spec.BackingServicePlanGuid, bsi.Spec.BackingServiceName, bsi.Name)
		result = fmt.Errorf("plan (%s) in bs(%s) for bsi (%s) not found",
			bsi.Spec.BackingServicePlanGuid, bsi.Spec.BackingServiceName, bsi.Name)
		return result
	}

	// ...

	glog.Infoln("bsi provisioning servicebroker_load, ", bsi.Name)
	////c.recorder.Eventf(bsi, "Provisioning", "bsi %s provisioning servicebroker_load", bsi.Name)

	servicebroker, err := servicebroker_load(c.Client, bs.GenerateName)
	if err != nil {
		result = err
		return result
	}

	_, err = servicebroker_update_instance(bsi, servicebroker)
	if err != nil {
		return err
	}

	return nil

}

func (c *BackingServiceInstanceControllerImpl) unbindInstanceHadoop(user string, bs *v1.BackingService, bsi *v1.BackingServiceInstance) (result error) {

	glog.Infoln("bsi to unbind hadoop ", bsi.Name, user)

	servicebroker, err := servicebroker_load(c.Client, bs.GenerateName)
	if err != nil {
		return err
	}

	glog.Infoln("servicebroker_unbinding")

	svc := &ServiceBinding{}

	svc.PlanId = bsi.Spec.BackingServicePlanGuid
	svc.ServiceId = bsi.Spec.BackingServiceSpecID
	svc.svc_instance_id = bsi.Spec.InstanceID

	var bindId string

	for idx, b := range bsi.Spec.Binding {
		if b.BindHadoopUser == user {
			bindId = b.BindUuid
			_, err = servicebroker_unbinding(bindId, svc, servicebroker)
			if err != nil {
				return err
			}
			bsi.Spec.Binding = append(bsi.Spec.Binding[:idx], bsi.Spec.Binding[idx+1:]...)
			delete(bsi.Annotations, v1.BindKind_HadoopUser)

			/*
				err = c.deploymentconfig_clear_envs(bsi.Namespace, dc)
				if err != nil {
					return err
				}
			*/
			glog.Infoln("bsi is unbound with hadoop", bsi.Name, user)

			bsi.Spec.Bound -= 1
			if bsi.Spec.Bound == 0 {
				bsi.Status.Phase = v1.BackingServiceInstancePhaseUnbound
			}

			break
		}
	}

	bsi.Status.Action = "" //remove_action_word(bsi.Status.Action, v1.BackingServiceInstanceActionToUnbind)

	return

}
func (c *BackingServiceInstanceControllerImpl) unbindInstance(dc string, bs *v1.BackingService, bsi *v1.BackingServiceInstance) (result error) {

	glog.Infoln("bsi to unbind ", bsi.Name)

	servicebroker, err := servicebroker_load(c.Client, bs.GenerateName)
	if err != nil {
		return err
	}

	glog.Infoln("servicebroker_unbinding")

	svc := &ServiceBinding{}

	svc.PlanId = bsi.Spec.BackingServicePlanGuid
	svc.ServiceId = bsi.Spec.BackingServiceSpecID
	svc.svc_instance_id = bsi.Spec.InstanceID

	var bindId string

	for idx, b := range bsi.Spec.Binding {
		if b.BindDeploymentConfig == dc {
			bindId = b.BindUuid
			_, err = servicebroker_unbinding(bindId, svc, servicebroker)
			if err != nil {
				return err
			}
			glog.Infoln("deploymentconfig_clear_envs")
			err = c.deploymentconfig_clear_envs(dc, bsi, &b)
			if err != nil && (!kerrors.IsNotFound(err)) {
				return err
			} else {
				bsi.Spec.Binding = append(bsi.Spec.Binding[:idx], bsi.Spec.Binding[idx+1:]...)
				delete(bsi.Annotations, dc)
			}

			/*
				err = c.deploymentconfig_clear_envs(bsi.Namespace, dc)
				if err != nil {
					return err
				}
			*/
			glog.Infoln("bsi is unbound ", bsi.Name)
			/*
				delete(bsi.Annotations, dc)
				bsi.Spec.BindDeploymentConfig = ""
				bsi.Spec.Credentials = nil
				bsi.Spec.BoundTime = nil
				bsi.Spec.BindUuid = ""
				bsi.Spec.Bound = false
			*/
			bsi.Spec.Bound -= 1
			if bsi.Spec.Bound == 0 {
				bsi.Status.Phase = v1.BackingServiceInstancePhaseUnbound
			}

			break
		}
	}

	bsi.Status.Action = "" //remove_action_word(bsi.Status.Action, v1.BackingServiceInstanceActionToUnbind)

	return

}

func (c *BackingServiceInstanceControllerImpl) get_deploymentconfig_name(bsi *v1.BackingServiceInstance, stat string) string {
	for dc, bound := range bsi.Annotations {
		if bound == stat {
			return dc
		}
	}
	return ""
}

//=========================================

func (c *BackingServiceInstanceControllerImpl) Get(namespace, name string) (*v1.BackingServiceInstance, error) {
	return c.lister.BackingServiceInstances(namespace).Get(name)
}
