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

// FakeServiceBrokers implements ServiceBrokerInterface
type FakeServiceBrokers struct {
	Fake *FakePrd
}

var servicebrokersResource = schema.GroupVersionResource{Group: "prd", Version: "", Resource: "servicebrokers"}

var servicebrokersKind = schema.GroupVersionKind{Group: "prd", Version: "", Kind: "ServiceBroker"}

// Get takes name of the serviceBroker, and returns the corresponding serviceBroker object, and an error if there is any.
func (c *FakeServiceBrokers) Get(name string, options v1.GetOptions) (result *prd.ServiceBroker, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(servicebrokersResource, name), &prd.ServiceBroker{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd.ServiceBroker), err
}

// List takes label and field selectors, and returns the list of ServiceBrokers that match those selectors.
func (c *FakeServiceBrokers) List(opts v1.ListOptions) (result *prd.ServiceBrokerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(servicebrokersResource, servicebrokersKind, opts), &prd.ServiceBrokerList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &prd.ServiceBrokerList{}
	for _, item := range obj.(*prd.ServiceBrokerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested serviceBrokers.
func (c *FakeServiceBrokers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(servicebrokersResource, opts))
}

// Create takes the representation of a serviceBroker and creates it.  Returns the server's representation of the serviceBroker, and an error, if there is any.
func (c *FakeServiceBrokers) Create(serviceBroker *prd.ServiceBroker) (result *prd.ServiceBroker, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(servicebrokersResource, serviceBroker), &prd.ServiceBroker{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd.ServiceBroker), err
}

// Update takes the representation of a serviceBroker and updates it. Returns the server's representation of the serviceBroker, and an error, if there is any.
func (c *FakeServiceBrokers) Update(serviceBroker *prd.ServiceBroker) (result *prd.ServiceBroker, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(servicebrokersResource, serviceBroker), &prd.ServiceBroker{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd.ServiceBroker), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeServiceBrokers) UpdateStatus(serviceBroker *prd.ServiceBroker) (*prd.ServiceBroker, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(servicebrokersResource, "status", serviceBroker), &prd.ServiceBroker{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd.ServiceBroker), err
}

// Delete takes name of the serviceBroker and deletes it. Returns an error if one occurs.
func (c *FakeServiceBrokers) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(servicebrokersResource, name), &prd.ServiceBroker{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeServiceBrokers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(servicebrokersResource, listOptions)

	_, err := c.Fake.Invokes(action, &prd.ServiceBrokerList{})
	return err
}

// Patch applies the patch and returns the patched serviceBroker.
func (c *FakeServiceBrokers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *prd.ServiceBroker, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(servicebrokersResource, name, data, subresources...), &prd.ServiceBroker{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd.ServiceBroker), err
}
