
/*
@asiainfo.com
*/


package backingservice

import (
	"log"

	"github.com/kubernetes-incubator/apiserver-builder/pkg/builders"

	"github.com/asiainfoldp/apiserver-servicebroker/pkg/apis/prd/v1"
	"github.com/asiainfoldp/apiserver-servicebroker/pkg/controller/sharedinformers"
	listers "github.com/asiainfoldp/apiserver-servicebroker/pkg/client/listers_generated/prd/v1"
)

// +controller:group=prd,version=v1,kind=BackingService,resource=backingservices
type BackingServiceControllerImpl struct {
	builders.DefaultControllerFns

	// lister indexes properties about BackingService
	lister listers.BackingServiceLister
}

// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *BackingServiceControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
	// Use the lister for indexing backingservices labels
	c.lister = arguments.GetSharedInformers().Factory.Prd().V1().BackingServices().Lister()
}

// Reconcile handles enqueued messages
func (c *BackingServiceControllerImpl) Reconcile(u *v1.BackingService) error {
	// Implement controller logic here
	log.Printf("Running reconcile BackingService for %s\n", u.Name)
	return nil
}

func (c *BackingServiceControllerImpl) Get(namespace, name string) (*v1.BackingService, error) {
	return c.lister.Get(name)
}
