/*
@asiainfo.com
*/
package fake

import (
	v1 "github.com/asiainfoldp/apiserver-servicebroker/pkg/client/clientset_generated/clientset/typed/prd/v1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakePrdV1 struct {
	*testing.Fake
}

func (c *FakePrdV1) BackingServices() v1.BackingServiceInterface {
	return &FakeBackingServices{c}
}

func (c *FakePrdV1) BackingServiceInstances(namespace string) v1.BackingServiceInstanceInterface {
	return &FakeBackingServiceInstances{c, namespace}
}

func (c *FakePrdV1) Bindings(namespace string) v1.BindingInterface {
	return &FakeBindings{c, namespace}
}

func (c *FakePrdV1) ServiceBrokers() v1.ServiceBrokerInterface {
	return &FakeServiceBrokers{c}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakePrdV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
