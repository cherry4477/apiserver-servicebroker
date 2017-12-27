/*
@asiainfo.com
*/
package fake

import (
	internalversion "github.com/asiainfoldp/apiserver-servicebroker/pkg/client/clientset_generated/internalclientset/typed/prd/internalversion"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakePrd struct {
	*testing.Fake
}

func (c *FakePrd) BackingServices() internalversion.BackingServiceInterface {
	return &FakeBackingServices{c}
}

func (c *FakePrd) BackingServiceInstances(namespace string) internalversion.BackingServiceInstanceInterface {
	return &FakeBackingServiceInstances{c, namespace}
}

func (c *FakePrd) ServiceBrokers() internalversion.ServiceBrokerInterface {
	return &FakeServiceBrokers{c}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakePrd) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
