
/*
@asiainfo.com
*/


package backingservice_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"
	"github.com/kubernetes-incubator/apiserver-builder/pkg/test"

	"github.com/asiainfoldp/apiserver-servicebroker/pkg/apis"
	"github.com/asiainfoldp/apiserver-servicebroker/pkg/client/clientset_generated/clientset"
	"github.com/asiainfoldp/apiserver-servicebroker/pkg/openapi"
	"github.com/asiainfoldp/apiserver-servicebroker/pkg/controller/sharedinformers"
	"github.com/asiainfoldp/apiserver-servicebroker/pkg/controller/backingservice"
)

var testenv *test.TestEnvironment
var config *rest.Config
var cs *clientset.Clientset
var shutdown chan struct{}
var controller *backingservice.BackingServiceController
var si *sharedinformers.SharedInformers

func TestBackingService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "BackingService Suite", []Reporter{test.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	testenv = test.NewTestEnvironment()
	config = testenv.Start(apis.GetAllApiBuilders(), openapi.GetOpenAPIDefinitions)
	cs = clientset.NewForConfigOrDie(config)

	shutdown = make(chan struct{})
	si = sharedinformers.NewSharedInformers(config, shutdown)
	controller = backingservice.NewBackingServiceController(config, si)
	controller.Run(shutdown)
})

var _ = AfterSuite(func() {
	close(shutdown)
	testenv.Stop()
})
