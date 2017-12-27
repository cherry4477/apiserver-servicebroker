/*
@asiainfo.com
*/
package v1

import (
	v1 "github.com/asiainfoldp/apiserver-servicebroker/pkg/apis/prd/v1"
	"github.com/asiainfoldp/apiserver-servicebroker/pkg/client/clientset_generated/clientset/scheme"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	rest "k8s.io/client-go/rest"
)

type PrdV1Interface interface {
	RESTClient() rest.Interface
	BackingServicesGetter
	BackingServiceInstancesGetter
	BindingsGetter
	ServiceBrokersGetter
}

// PrdV1Client is used to interact with features provided by the prd.asiainfo.com group.
type PrdV1Client struct {
	restClient rest.Interface
}

func (c *PrdV1Client) BackingServices() BackingServiceInterface {
	return newBackingServices(c)
}

func (c *PrdV1Client) BackingServiceInstances(namespace string) BackingServiceInstanceInterface {
	return newBackingServiceInstances(c, namespace)
}

func (c *PrdV1Client) Bindings(namespace string) BindingInterface {
	return newBindings(c, namespace)
}

func (c *PrdV1Client) ServiceBrokers() ServiceBrokerInterface {
	return newServiceBrokers(c)
}

// NewForConfig creates a new PrdV1Client for the given config.
func NewForConfig(c *rest.Config) (*PrdV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &PrdV1Client{client}, nil
}

// NewForConfigOrDie creates a new PrdV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *PrdV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new PrdV1Client for the given RESTClient.
func New(c rest.Interface) *PrdV1Client {
	return &PrdV1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *PrdV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
