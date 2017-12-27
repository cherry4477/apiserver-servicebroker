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

// FakeBackingServices implements BackingServiceInterface
type FakeBackingServices struct {
	Fake *FakePrd
}

var backingservicesResource = schema.GroupVersionResource{Group: "prd", Version: "", Resource: "backingservices"}

var backingservicesKind = schema.GroupVersionKind{Group: "prd", Version: "", Kind: "BackingService"}

// Get takes name of the backingService, and returns the corresponding backingService object, and an error if there is any.
func (c *FakeBackingServices) Get(name string, options v1.GetOptions) (result *prd.BackingService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(backingservicesResource, name), &prd.BackingService{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd.BackingService), err
}

// List takes label and field selectors, and returns the list of BackingServices that match those selectors.
func (c *FakeBackingServices) List(opts v1.ListOptions) (result *prd.BackingServiceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(backingservicesResource, backingservicesKind, opts), &prd.BackingServiceList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &prd.BackingServiceList{}
	for _, item := range obj.(*prd.BackingServiceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested backingServices.
func (c *FakeBackingServices) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(backingservicesResource, opts))
}

// Create takes the representation of a backingService and creates it.  Returns the server's representation of the backingService, and an error, if there is any.
func (c *FakeBackingServices) Create(backingService *prd.BackingService) (result *prd.BackingService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(backingservicesResource, backingService), &prd.BackingService{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd.BackingService), err
}

// Update takes the representation of a backingService and updates it. Returns the server's representation of the backingService, and an error, if there is any.
func (c *FakeBackingServices) Update(backingService *prd.BackingService) (result *prd.BackingService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(backingservicesResource, backingService), &prd.BackingService{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd.BackingService), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeBackingServices) UpdateStatus(backingService *prd.BackingService) (*prd.BackingService, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(backingservicesResource, "status", backingService), &prd.BackingService{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd.BackingService), err
}

// Delete takes name of the backingService and deletes it. Returns an error if one occurs.
func (c *FakeBackingServices) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(backingservicesResource, name), &prd.BackingService{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBackingServices) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(backingservicesResource, listOptions)

	_, err := c.Fake.Invokes(action, &prd.BackingServiceList{})
	return err
}

// Patch applies the patch and returns the patched backingService.
func (c *FakeBackingServices) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *prd.BackingService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(backingservicesResource, name, data, subresources...), &prd.BackingService{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd.BackingService), err
}
