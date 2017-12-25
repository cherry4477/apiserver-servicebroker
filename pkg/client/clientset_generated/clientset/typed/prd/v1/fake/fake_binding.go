/*
@asiainfo.com
*/
package fake

import (
	prd_v1 "github.com/asiainfoldp/apiserver-servicebroker/pkg/apis/prd/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeBindings implements BindingInterface
type FakeBindings struct {
	Fake *FakePrdV1
	ns   string
}

var bindingsResource = schema.GroupVersionResource{Group: "prd.asiainfo.com", Version: "v1", Resource: "bindings"}

var bindingsKind = schema.GroupVersionKind{Group: "prd.asiainfo.com", Version: "v1", Kind: "Binding"}

// Get takes name of the binding, and returns the corresponding binding object, and an error if there is any.
func (c *FakeBindings) Get(name string, options v1.GetOptions) (result *prd_v1.Binding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(bindingsResource, c.ns, name), &prd_v1.Binding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*prd_v1.Binding), err
}

// List takes label and field selectors, and returns the list of Bindings that match those selectors.
func (c *FakeBindings) List(opts v1.ListOptions) (result *prd_v1.BindingList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(bindingsResource, bindingsKind, c.ns, opts), &prd_v1.BindingList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &prd_v1.BindingList{}
	for _, item := range obj.(*prd_v1.BindingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested bindings.
func (c *FakeBindings) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(bindingsResource, c.ns, opts))

}

// Create takes the representation of a binding and creates it.  Returns the server's representation of the binding, and an error, if there is any.
func (c *FakeBindings) Create(binding *prd_v1.Binding) (result *prd_v1.Binding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(bindingsResource, c.ns, binding), &prd_v1.Binding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*prd_v1.Binding), err
}

// Update takes the representation of a binding and updates it. Returns the server's representation of the binding, and an error, if there is any.
func (c *FakeBindings) Update(binding *prd_v1.Binding) (result *prd_v1.Binding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(bindingsResource, c.ns, binding), &prd_v1.Binding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*prd_v1.Binding), err
}

// Delete takes name of the binding and deletes it. Returns an error if one occurs.
func (c *FakeBindings) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(bindingsResource, c.ns, name), &prd_v1.Binding{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBindings) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(bindingsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &prd_v1.BindingList{})
	return err
}

// Patch applies the patch and returns the patched binding.
func (c *FakeBindings) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *prd_v1.Binding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(bindingsResource, c.ns, name, data, subresources...), &prd_v1.Binding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*prd_v1.Binding), err
}
