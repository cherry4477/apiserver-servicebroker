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

// BackingServiceInstancesGetter has a method to return a BackingServiceInstanceInterface.
// A group's client should implement this interface.
type BackingServiceInstancesGetter interface {
	BackingServiceInstances(namespace string) BackingServiceInstanceInterface
}

// BackingServiceInstanceInterface has methods to work with BackingServiceInstance resources.
type BackingServiceInstanceInterface interface {
	Create(*prd.BackingServiceInstance) (*prd.BackingServiceInstance, error)
	Update(*prd.BackingServiceInstance) (*prd.BackingServiceInstance, error)
	UpdateStatus(*prd.BackingServiceInstance) (*prd.BackingServiceInstance, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*prd.BackingServiceInstance, error)
	List(opts v1.ListOptions) (*prd.BackingServiceInstanceList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *prd.BackingServiceInstance, err error)
	BackingServiceInstanceExpansion
}

// backingServiceInstances implements BackingServiceInstanceInterface
type backingServiceInstances struct {
	client rest.Interface
	ns     string
}

// newBackingServiceInstances returns a BackingServiceInstances
func newBackingServiceInstances(c *PrdClient, namespace string) *backingServiceInstances {
	return &backingServiceInstances{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the backingServiceInstance, and returns the corresponding backingServiceInstance object, and an error if there is any.
func (c *backingServiceInstances) Get(name string, options v1.GetOptions) (result *prd.BackingServiceInstance, err error) {
	result = &prd.BackingServiceInstance{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("backingserviceinstances").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of BackingServiceInstances that match those selectors.
func (c *backingServiceInstances) List(opts v1.ListOptions) (result *prd.BackingServiceInstanceList, err error) {
	result = &prd.BackingServiceInstanceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("backingserviceinstances").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested backingServiceInstances.
func (c *backingServiceInstances) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("backingserviceinstances").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a backingServiceInstance and creates it.  Returns the server's representation of the backingServiceInstance, and an error, if there is any.
func (c *backingServiceInstances) Create(backingServiceInstance *prd.BackingServiceInstance) (result *prd.BackingServiceInstance, err error) {
	result = &prd.BackingServiceInstance{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("backingserviceinstances").
		Body(backingServiceInstance).
		Do().
		Into(result)
	return
}

// Update takes the representation of a backingServiceInstance and updates it. Returns the server's representation of the backingServiceInstance, and an error, if there is any.
func (c *backingServiceInstances) Update(backingServiceInstance *prd.BackingServiceInstance) (result *prd.BackingServiceInstance, err error) {
	result = &prd.BackingServiceInstance{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("backingserviceinstances").
		Name(backingServiceInstance.Name).
		Body(backingServiceInstance).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *backingServiceInstances) UpdateStatus(backingServiceInstance *prd.BackingServiceInstance) (result *prd.BackingServiceInstance, err error) {
	result = &prd.BackingServiceInstance{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("backingserviceinstances").
		Name(backingServiceInstance.Name).
		SubResource("status").
		Body(backingServiceInstance).
		Do().
		Into(result)
	return
}

// Delete takes name of the backingServiceInstance and deletes it. Returns an error if one occurs.
func (c *backingServiceInstances) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("backingserviceinstances").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *backingServiceInstances) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("backingserviceinstances").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched backingServiceInstance.
func (c *backingServiceInstances) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *prd.BackingServiceInstance, err error) {
	result = &prd.BackingServiceInstance{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("backingserviceinstances").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
