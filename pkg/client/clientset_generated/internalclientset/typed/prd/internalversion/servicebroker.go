/*
@asiainfo.com
*/
package internalversion

import (
	prd "github.com/asiainfoldp/apiserver-servicebroker/pkg/apis/prd"
	scheme "github.com/asiainfoldp/apiserver-servicebroker/pkg/client/clientset_generated/internalclientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ServiceBrokersGetter has a method to return a ServiceBrokerInterface.
// A group's client should implement this interface.
type ServiceBrokersGetter interface {
	ServiceBrokers() ServiceBrokerInterface
}

// ServiceBrokerInterface has methods to work with ServiceBroker resources.
type ServiceBrokerInterface interface {
	Create(*prd.ServiceBroker) (*prd.ServiceBroker, error)
	Update(*prd.ServiceBroker) (*prd.ServiceBroker, error)
	UpdateStatus(*prd.ServiceBroker) (*prd.ServiceBroker, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*prd.ServiceBroker, error)
	List(opts v1.ListOptions) (*prd.ServiceBrokerList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *prd.ServiceBroker, err error)
	ServiceBrokerExpansion
}

// serviceBrokers implements ServiceBrokerInterface
type serviceBrokers struct {
	client rest.Interface
}

// newServiceBrokers returns a ServiceBrokers
func newServiceBrokers(c *PrdClient) *serviceBrokers {
	return &serviceBrokers{
		client: c.RESTClient(),
	}
}

// Get takes name of the serviceBroker, and returns the corresponding serviceBroker object, and an error if there is any.
func (c *serviceBrokers) Get(name string, options v1.GetOptions) (result *prd.ServiceBroker, err error) {
	result = &prd.ServiceBroker{}
	err = c.client.Get().
		Resource("servicebrokers").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ServiceBrokers that match those selectors.
func (c *serviceBrokers) List(opts v1.ListOptions) (result *prd.ServiceBrokerList, err error) {
	result = &prd.ServiceBrokerList{}
	err = c.client.Get().
		Resource("servicebrokers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested serviceBrokers.
func (c *serviceBrokers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Resource("servicebrokers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a serviceBroker and creates it.  Returns the server's representation of the serviceBroker, and an error, if there is any.
func (c *serviceBrokers) Create(serviceBroker *prd.ServiceBroker) (result *prd.ServiceBroker, err error) {
	result = &prd.ServiceBroker{}
	err = c.client.Post().
		Resource("servicebrokers").
		Body(serviceBroker).
		Do().
		Into(result)
	return
}

// Update takes the representation of a serviceBroker and updates it. Returns the server's representation of the serviceBroker, and an error, if there is any.
func (c *serviceBrokers) Update(serviceBroker *prd.ServiceBroker) (result *prd.ServiceBroker, err error) {
	result = &prd.ServiceBroker{}
	err = c.client.Put().
		Resource("servicebrokers").
		Name(serviceBroker.Name).
		Body(serviceBroker).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *serviceBrokers) UpdateStatus(serviceBroker *prd.ServiceBroker) (result *prd.ServiceBroker, err error) {
	result = &prd.ServiceBroker{}
	err = c.client.Put().
		Resource("servicebrokers").
		Name(serviceBroker.Name).
		SubResource("status").
		Body(serviceBroker).
		Do().
		Into(result)
	return
}

// Delete takes name of the serviceBroker and deletes it. Returns an error if one occurs.
func (c *serviceBrokers) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("servicebrokers").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *serviceBrokers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Resource("servicebrokers").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched serviceBroker.
func (c *serviceBrokers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *prd.ServiceBroker, err error) {
	result = &prd.ServiceBroker{}
	err = c.client.Patch(pt).
		Resource("servicebrokers").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
