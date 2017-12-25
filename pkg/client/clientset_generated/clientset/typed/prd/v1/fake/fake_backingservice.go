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

// FakeBackingServices implements BackingServiceInterface
type FakeBackingServices struct {
	Fake *FakePrdV1
}

var backingservicesResource = schema.GroupVersionResource{Group: "prd.asiainfo.com", Version: "v1", Resource: "backingservices"}

var backingservicesKind = schema.GroupVersionKind{Group: "prd.asiainfo.com", Version: "v1", Kind: "BackingService"}

// Get takes name of the backingService, and returns the corresponding backingService object, and an error if there is any.
func (c *FakeBackingServices) Get(name string, options v1.GetOptions) (result *prd_v1.BackingService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(backingservicesResource, name), &prd_v1.BackingService{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd_v1.BackingService), err
}

// List takes label and field selectors, and returns the list of BackingServices that match those selectors.
func (c *FakeBackingServices) List(opts v1.ListOptions) (result *prd_v1.BackingServiceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(backingservicesResource, backingservicesKind, opts), &prd_v1.BackingServiceList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &prd_v1.BackingServiceList{}
	for _, item := range obj.(*prd_v1.BackingServiceList).Items {
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
func (c *FakeBackingServices) Create(backingService *prd_v1.BackingService) (result *prd_v1.BackingService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(backingservicesResource, backingService), &prd_v1.BackingService{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd_v1.BackingService), err
}

// Update takes the representation of a backingService and updates it. Returns the server's representation of the backingService, and an error, if there is any.
func (c *FakeBackingServices) Update(backingService *prd_v1.BackingService) (result *prd_v1.BackingService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(backingservicesResource, backingService), &prd_v1.BackingService{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd_v1.BackingService), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeBackingServices) UpdateStatus(backingService *prd_v1.BackingService) (*prd_v1.BackingService, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(backingservicesResource, "status", backingService), &prd_v1.BackingService{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd_v1.BackingService), err
}

// Delete takes name of the backingService and deletes it. Returns an error if one occurs.
func (c *FakeBackingServices) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(backingservicesResource, name), &prd_v1.BackingService{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBackingServices) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(backingservicesResource, listOptions)

	_, err := c.Fake.Invokes(action, &prd_v1.BackingServiceList{})
	return err
}

// Patch applies the patch and returns the patched backingService.
func (c *FakeBackingServices) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *prd_v1.BackingService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(backingservicesResource, name, data, subresources...), &prd_v1.BackingService{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd_v1.BackingService), err
}
