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

// BindingsGetter has a method to return a BindingInterface.
// A group's client should implement this interface.
type BindingsGetter interface {
	Bindings(namespace string) BindingInterface
}

// BindingInterface has methods to work with Binding resources.
type BindingInterface interface {
	Create(*v1.Binding) (*v1.Binding, error)
	Update(*v1.Binding) (*v1.Binding, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.Binding, error)
	List(opts meta_v1.ListOptions) (*v1.BindingList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Binding, err error)
	BindingExpansion
}

// bindings implements BindingInterface
type bindings struct {
	client rest.Interface
	ns     string
}

// newBindings returns a Bindings
func newBindings(c *PrdV1Client, namespace string) *bindings {
	return &bindings{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the binding, and returns the corresponding binding object, and an error if there is any.
func (c *bindings) Get(name string, options meta_v1.GetOptions) (result *v1.Binding, err error) {
	result = &v1.Binding{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("bindings").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Bindings that match those selectors.
func (c *bindings) List(opts meta_v1.ListOptions) (result *v1.BindingList, err error) {
	result = &v1.BindingList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("bindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested bindings.
func (c *bindings) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("bindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a binding and creates it.  Returns the server's representation of the binding, and an error, if there is any.
func (c *bindings) Create(binding *v1.Binding) (result *v1.Binding, err error) {
	result = &v1.Binding{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("bindings").
		Body(binding).
		Do().
		Into(result)
	return
}

// Update takes the representation of a binding and updates it. Returns the server's representation of the binding, and an error, if there is any.
func (c *bindings) Update(binding *v1.Binding) (result *v1.Binding, err error) {
	result = &v1.Binding{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("bindings").
		Name(binding.Name).
		Body(binding).
		Do().
		Into(result)
	return
}

// Delete takes name of the binding and deletes it. Returns an error if one occurs.
func (c *bindings) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("bindings").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *bindings) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("bindings").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched binding.
func (c *bindings) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Binding, err error) {
	result = &v1.Binding{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("bindings").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
