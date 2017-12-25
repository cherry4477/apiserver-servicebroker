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

// FakeBackingServiceInstances implements BackingServiceInstanceInterface
type FakeBackingServiceInstances struct {
	Fake *FakePrdV1
	ns   string
}

var backingserviceinstancesResource = schema.GroupVersionResource{Group: "prd.asiainfo.com", Version: "v1", Resource: "backingserviceinstances"}

var backingserviceinstancesKind = schema.GroupVersionKind{Group: "prd.asiainfo.com", Version: "v1", Kind: "BackingServiceInstance"}

// Get takes name of the backingServiceInstance, and returns the corresponding backingServiceInstance object, and an error if there is any.
func (c *FakeBackingServiceInstances) Get(name string, options v1.GetOptions) (result *prd_v1.BackingServiceInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(backingserviceinstancesResource, c.ns, name), &prd_v1.BackingServiceInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*prd_v1.BackingServiceInstance), err
}

// List takes label and field selectors, and returns the list of BackingServiceInstances that match those selectors.
func (c *FakeBackingServiceInstances) List(opts v1.ListOptions) (result *prd_v1.BackingServiceInstanceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(backingserviceinstancesResource, backingserviceinstancesKind, c.ns, opts), &prd_v1.BackingServiceInstanceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &prd_v1.BackingServiceInstanceList{}
	for _, item := range obj.(*prd_v1.BackingServiceInstanceList).Items {
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
func (c *FakeBackingServiceInstances) Create(backingServiceInstance *prd_v1.BackingServiceInstance) (result *prd_v1.BackingServiceInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(backingserviceinstancesResource, c.ns, backingServiceInstance), &prd_v1.BackingServiceInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*prd_v1.BackingServiceInstance), err
}

// Update takes the representation of a backingServiceInstance and updates it. Returns the server's representation of the backingServiceInstance, and an error, if there is any.
func (c *FakeBackingServiceInstances) Update(backingServiceInstance *prd_v1.BackingServiceInstance) (result *prd_v1.BackingServiceInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(backingserviceinstancesResource, c.ns, backingServiceInstance), &prd_v1.BackingServiceInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*prd_v1.BackingServiceInstance), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeBackingServiceInstances) UpdateStatus(backingServiceInstance *prd_v1.BackingServiceInstance) (*prd_v1.BackingServiceInstance, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(backingserviceinstancesResource, "status", c.ns, backingServiceInstance), &prd_v1.BackingServiceInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*prd_v1.BackingServiceInstance), err
}

// Delete takes name of the backingServiceInstance and deletes it. Returns an error if one occurs.
func (c *FakeBackingServiceInstances) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(backingserviceinstancesResource, c.ns, name), &prd_v1.BackingServiceInstance{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBackingServiceInstances) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(backingserviceinstancesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &prd_v1.BackingServiceInstanceList{})
	return err
}

// Patch applies the patch and returns the patched backingServiceInstance.
func (c *FakeBackingServiceInstances) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *prd_v1.BackingServiceInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(backingserviceinstancesResource, c.ns, name, data, subresources...), &prd_v1.BackingServiceInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*prd_v1.BackingServiceInstance), err
}
