/*
@asiainfo.com
*/
package v1

import (
	v1 "github.com/asiainfoldp/apiserver-servicebroker/pkg/apis/prd/v1"
	scheme "github.com/asiainfoldp/apiserver-servicebroker/pkg/client/clientset_generated/clientset/scheme"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	Create(*v1.ServiceBroker) (*v1.ServiceBroker, error)
	Update(*v1.ServiceBroker) (*v1.ServiceBroker, error)
	UpdateStatus(*v1.ServiceBroker) (*v1.ServiceBroker, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.ServiceBroker, error)
	List(opts meta_v1.ListOptions) (*v1.ServiceBrokerList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ServiceBroker, err error)
	ServiceBrokerExpansion
}

// serviceBrokers implements ServiceBrokerInterface
type serviceBrokers struct {
	client rest.Interface
}

// newServiceBrokers returns a ServiceBrokers
func newServiceBrokers(c *PrdV1Client) *serviceBrokers {
	return &serviceBrokers{
		client: c.RESTClient(),
	}
}

// Get takes name of the serviceBroker, and returns the corresponding serviceBroker object, and an error if there is any.
func (c *serviceBrokers) Get(name string, options meta_v1.GetOptions) (result *v1.ServiceBroker, err error) {
	result = &v1.ServiceBroker{}
	err = c.client.Get().
		Resource("servicebrokers").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ServiceBrokers that match those selectors.
func (c *serviceBrokers) List(opts meta_v1.ListOptions) (result *v1.ServiceBrokerList, err error) {
	result = &v1.ServiceBrokerList{}
	err = c.client.Get().
		Resource("servicebrokers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested serviceBrokers.
func (c *serviceBrokers) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Resource("servicebrokers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a serviceBroker and creates it.  Returns the server's representation of the serviceBroker, and an error, if there is any.
func (c *serviceBrokers) Create(serviceBroker *v1.ServiceBroker) (result *v1.ServiceBroker, err error) {
	result = &v1.ServiceBroker{}
	err = c.client.Post().
		Resource("servicebrokers").
		Body(serviceBroker).
		Do().
		Into(result)
	return
}

// Update takes the representation of a serviceBroker and updates it. Returns the server's representation of the serviceBroker, and an error, if there is any.
func (c *serviceBrokers) Update(serviceBroker *v1.ServiceBroker) (result *v1.ServiceBroker, err error) {
	result = &v1.ServiceBroker{}
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

func (c *serviceBrokers) UpdateStatus(serviceBroker *v1.ServiceBroker) (result *v1.ServiceBroker, err error) {
	result = &v1.ServiceBroker{}
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
func (c *serviceBrokers) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("servicebrokers").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *serviceBrokers) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Resource("servicebrokers").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched serviceBroker.
func (c *serviceBrokers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ServiceBroker, err error) {
	result = &v1.ServiceBroker{}
	err = c.client.Patch(pt).
		Resource("servicebrokers").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
