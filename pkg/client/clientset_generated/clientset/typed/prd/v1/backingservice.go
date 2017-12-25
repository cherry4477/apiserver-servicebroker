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

// BackingServicesGetter has a method to return a BackingServiceInterface.
// A group's client should implement this interface.
type BackingServicesGetter interface {
	BackingServices() BackingServiceInterface
}

// BackingServiceInterface has methods to work with BackingService resources.
type BackingServiceInterface interface {
	Create(*v1.BackingService) (*v1.BackingService, error)
	Update(*v1.BackingService) (*v1.BackingService, error)
	UpdateStatus(*v1.BackingService) (*v1.BackingService, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.BackingService, error)
	List(opts meta_v1.ListOptions) (*v1.BackingServiceList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.BackingService, err error)
	BackingServiceExpansion
}

// backingServices implements BackingServiceInterface
type backingServices struct {
	client rest.Interface
}

// newBackingServices returns a BackingServices
func newBackingServices(c *PrdV1Client) *backingServices {
	return &backingServices{
		client: c.RESTClient(),
	}
}

// Get takes name of the backingService, and returns the corresponding backingService object, and an error if there is any.
func (c *backingServices) Get(name string, options meta_v1.GetOptions) (result *v1.BackingService, err error) {
	result = &v1.BackingService{}
	err = c.client.Get().
		Resource("backingservices").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of BackingServices that match those selectors.
func (c *backingServices) List(opts meta_v1.ListOptions) (result *v1.BackingServiceList, err error) {
	result = &v1.BackingServiceList{}
	err = c.client.Get().
		Resource("backingservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested backingServices.
func (c *backingServices) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Resource("backingservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a backingService and creates it.  Returns the server's representation of the backingService, and an error, if there is any.
func (c *backingServices) Create(backingService *v1.BackingService) (result *v1.BackingService, err error) {
	result = &v1.BackingService{}
	err = c.client.Post().
		Resource("backingservices").
		Body(backingService).
		Do().
		Into(result)
	return
}

// Update takes the representation of a backingService and updates it. Returns the server's representation of the backingService, and an error, if there is any.
func (c *backingServices) Update(backingService *v1.BackingService) (result *v1.BackingService, err error) {
	result = &v1.BackingService{}
	err = c.client.Put().
		Resource("backingservices").
		Name(backingService.Name).
		Body(backingService).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *backingServices) UpdateStatus(backingService *v1.BackingService) (result *v1.BackingService, err error) {
	result = &v1.BackingService{}
	err = c.client.Put().
		Resource("backingservices").
		Name(backingService.Name).
		SubResource("status").
		Body(backingService).
		Do().
		Into(result)
	return
}

// Delete takes name of the backingService and deletes it. Returns an error if one occurs.
func (c *backingServices) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("backingservices").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *backingServices) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Resource("backingservices").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched backingService.
func (c *backingServices) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.BackingService, err error) {
	result = &v1.BackingService{}
	err = c.client.Patch(pt).
		Resource("backingservices").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
