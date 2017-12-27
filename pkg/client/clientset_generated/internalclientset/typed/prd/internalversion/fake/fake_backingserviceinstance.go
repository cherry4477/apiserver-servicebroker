/*
@asiainfo.com
*/
package fake

import (
	prd "github.com/asiainfoldp/apiserver-servicebroker/pkg/apis/prd"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeBackingServiceInstances implements BackingServiceInstanceInterface
type FakeBackingServiceInstances struct {
	Fake *FakePrd
	ns   string
}

var backingserviceinstancesResource = schema.GroupVersionResource{Group: "prd", Version: "", Resource: "backingserviceinstances"}

var backingserviceinstancesKind = schema.GroupVersionKind{Group: "prd", Version: "", Kind: "BackingServiceInstance"}

// Get takes name of the backingServiceInstance, and returns the corresponding backingServiceInstance object, and an error if there is any.
func (c *FakeBackingServiceInstances) Get(name string, options v1.GetOptions) (result *prd.BackingServiceInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(backingserviceinstancesResource, c.ns, name), &prd.BackingServiceInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*prd.BackingServiceInstance), err
}

// List takes label and field selectors, and returns the list of BackingServiceInstances that match those selectors.
func (c *FakeBackingServiceInstances) List(opts v1.ListOptions) (result *prd.BackingServiceInstanceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(backingserviceinstancesResource, backingserviceinstancesKind, c.ns, opts), &prd.BackingServiceInstanceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &prd.BackingServiceInstanceList{}
	for _, item := range obj.(*prd.BackingServiceInstanceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested backingServiceInstances.
func (c *FakeBackingServiceInstances) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(backingserviceinstancesResource, c.ns, opts))

}

// Create takes the representation of a backingServiceInstance and creates it.  Returns the server's representation of the backingServiceInstance, and an error, if there is any.
func (c *FakeBackingServiceInstances) Create(backingServiceInstance *prd.BackingServiceInstance) (result *prd.BackingServiceInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(backingserviceinstancesResource, c.ns, backingServiceInstance), &prd.BackingServiceInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*prd.BackingServiceInstance), err
}

// Update takes the representation of a backingServiceInstance and updates it. Returns the server's representation of the backingServiceInstance, and an error, if there is any.
func (c *FakeBackingServiceInstances) Update(backingServiceInstance *prd.BackingServiceInstance) (result *prd.BackingServiceInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(backingserviceinstancesResource, c.ns, backingServiceInstance), &prd.BackingServiceInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*prd.BackingServiceInstance), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeBackingServiceInstances) UpdateStatus(backingServiceInstance *prd.BackingServiceInstance) (*prd.BackingServiceInstance, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(backingserviceinstancesResource, "status", c.ns, backingServiceInstance), &prd.BackingServiceInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*prd.BackingServiceInstance), err
}

// Delete takes name of the backingServiceInstance and deletes it. Returns an error if one occurs.
func (c *FakeBackingServiceInstances) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(backingserviceinstancesResource, c.ns, name), &prd.BackingServiceInstance{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBackingServiceInstances) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(backingserviceinstancesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &prd.BackingServiceInstanceList{})
	return err
}

// Patch applies the patch and returns the patched backingServiceInstance.
func (c *FakeBackingServiceInstances) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *prd.BackingServiceInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(backingserviceinstancesResource, c.ns, name, data, subresources...), &prd.BackingServiceInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*prd.BackingServiceInstance), err
}
