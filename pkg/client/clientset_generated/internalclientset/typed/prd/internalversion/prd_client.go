/*
@asiainfo.com
*/
package internalversion

import (
	"github.com/asiainfoldp/apiserver-servicebroker/pkg/client/clientset_generated/internalclientset/scheme"
	rest "k8s.io/client-go/rest"
)

type PrdInterface interface {
	RESTClient() rest.Interface
	BackingServicesGetter
	BackingServiceInstancesGetter
	ServiceBrokersGetter
}

// PrdClient is used to interact with features provided by the prd group.
type PrdClient struct {
	restClient rest.Interface
}

func (c *PrdClient) BackingServices() BackingServiceInterface {
	return newBackingServices(c)
}

func (c *PrdClient) BackingServiceInstances(namespace string) BackingServiceInstanceInterface {
	return newBackingServiceInstances(c, namespace)
}

func (c *PrdClient) ServiceBrokers() ServiceBrokerInterface {
	return newServiceBrokers(c)
}

// NewForConfig creates a new PrdClient for the given config.
func NewForConfig(c *rest.Config) (*PrdClient, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &PrdClient{client}, nil
}

// NewForConfigOrDie creates a new PrdClient for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *PrdClient {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new PrdClient for the given RESTClient.
func New(c rest.Interface) *PrdClient {
	return &PrdClient{c}
}

func setConfigDefaults(config *rest.Config) error {
	g, err := scheme.Registry.Group("prd")
	if err != nil {
		return err
	}

	config.APIPath = "/apis"
	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}
	if config.GroupVersion == nil || config.GroupVersion.Group != g.GroupVersion.Group {
		gv := g.GroupVersion
		config.GroupVersion = &gv
	}
	config.NegotiatedSerializer = scheme.Codecs

	if config.QPS == 0 {
		config.QPS = 5
	}
	if config.Burst == 0 {
		config.Burst = 10
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *PrdClient) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
