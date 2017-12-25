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

// FakeServiceBrokers implements ServiceBrokerInterface
type FakeServiceBrokers struct {
	Fake *FakePrdV1
}

var servicebrokersResource = schema.GroupVersionResource{Group: "prd.asiainfo.com", Version: "v1", Resource: "servicebrokers"}

var servicebrokersKind = schema.GroupVersionKind{Group: "prd.asiainfo.com", Version: "v1", Kind: "ServiceBroker"}

// Get takes name of the serviceBroker, and returns the corresponding serviceBroker object, and an error if there is any.
func (c *FakeServiceBrokers) Get(name string, options v1.GetOptions) (result *prd_v1.ServiceBroker, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(servicebrokersResource, name), &prd_v1.ServiceBroker{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd_v1.ServiceBroker), err
}

// List takes label and field selectors, and returns the list of ServiceBrokers that match those selectors.
func (c *FakeServiceBrokers) List(opts v1.ListOptions) (result *prd_v1.ServiceBrokerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(servicebrokersResource, servicebrokersKind, opts), &prd_v1.ServiceBrokerList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &prd_v1.ServiceBrokerList{}
	for _, item := range obj.(*prd_v1.ServiceBrokerList).Items {
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
func (c *FakeServiceBrokers) Create(serviceBroker *prd_v1.ServiceBroker) (result *prd_v1.ServiceBroker, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(servicebrokersResource, serviceBroker), &prd_v1.ServiceBroker{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd_v1.ServiceBroker), err
}

// Update takes the representation of a serviceBroker and updates it. Returns the server's representation of the serviceBroker, and an error, if there is any.
func (c *FakeServiceBrokers) Update(serviceBroker *prd_v1.ServiceBroker) (result *prd_v1.ServiceBroker, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(servicebrokersResource, serviceBroker), &prd_v1.ServiceBroker{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd_v1.ServiceBroker), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeServiceBrokers) UpdateStatus(serviceBroker *prd_v1.ServiceBroker) (*prd_v1.ServiceBroker, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(servicebrokersResource, "status", serviceBroker), &prd_v1.ServiceBroker{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd_v1.ServiceBroker), err
}

// Delete takes name of the serviceBroker and deletes it. Returns an error if one occurs.
func (c *FakeServiceBrokers) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(servicebrokersResource, name), &prd_v1.ServiceBroker{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeServiceBrokers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(servicebrokersResource, listOptions)

	_, err := c.Fake.Invokes(action, &prd_v1.ServiceBrokerList{})
	return err
}

// Patch applies the patch and returns the patched serviceBroker.
func (c *FakeServiceBrokers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *prd_v1.ServiceBroker, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(servicebrokersResource, name, data, subresources...), &prd_v1.ServiceBroker{})
	if obj == nil {
		return nil, err
	}
	return obj.(*prd_v1.ServiceBroker), err
}
