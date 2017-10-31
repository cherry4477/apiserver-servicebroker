
/*
@asiainfo.com
*/


package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/asiainfoldp/apiserver-servicebroker/pkg/apis/prd/v1"
	. "github.com/asiainfoldp/apiserver-servicebroker/pkg/client/clientset_generated/clientset/typed/prd/v1"
)

var _ = Describe("BackingServiceInstance", func() {
	var instance BackingServiceInstance
	var expected BackingServiceInstance
	var client BackingServiceInstanceInterface

	BeforeEach(func() {
		instance = BackingServiceInstance{}
		instance.Name = "instance-1"

		expected = instance
	})

	AfterEach(func() {
		client.Delete(instance.Name, &metav1.DeleteOptions{})
	})

	Describe("when sending a binding request", func() {
		It("should return success", func() {
			client = cs.PrdV1Client.Backingserviceinstances("backingserviceinstance-test-binding")
			_, err := client.Create(&instance)
			Expect(err).ShouldNot(HaveOccurred())

			binding := &Binding{}
			binding.Name = instance.Name
			restClient := cs.PrdV1Client.RESTClient()
			err = restClient.Post().Namespace("backingserviceinstance-test-binding").
				Name(instance.Name).
				Resource("backingserviceinstances").
				SubResource("binding").
				Body(binding).Do().Error()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
