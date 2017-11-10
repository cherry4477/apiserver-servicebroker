/*
@asiainfo.com
*/

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"

	kapi "k8s.io/api/core/v1"
	genericvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/asiainfoldp/apiserver-servicebroker/pkg/apis/prd"
	prdutil "github.com/asiainfoldp/apiserver-servicebroker/pkg/apis/prd/util"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// ServiceBroker
// +k8s:openapi-gen=true
// +resource:path=servicebrokers,strategy=ServiceBrokerStrategy
type ServiceBroker struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceBrokerSpec   `json:"spec,omitempty"`
	Status ServiceBrokerStatus `json:"status,omitempty"`
}

// ServiceBrokerSpec defines the desired state of ServiceBroker
type ServiceBrokerSpec struct {
	// url defines the address of a ServiceBroker service
	Url string `json:"url" description:"url defines the address of a ServiceBroker service"`
	// name defines the name of a ServiceBroker service
	Name string `json:"name" description:"name defines the name of a ServiceBroker service"`
	// username defines the username to access ServiceBroker service
	UserName string `json:"username" description:"username defines the username to access ServiceBroker service"`
	// password defines the password to access ServiceBroker service
	Password string `json:"password" description:"password defines the password to access ServiceBroker service"`
	// Finalizers is an opaque list of values that must be empty to permanently remove object from storage
	Finalizers []kapi.FinalizerName `json:"finalizers,omitempty" description:"an opaque list of values that must be empty to permanently remove object from storage"`
}

// ServiceBrokerStatus defines the observed state of ServiceBroker
type ServiceBrokerStatus struct {
	// Phase is the current lifecycle phase of the project
	Phase string /*ServiceBrokerPhase*/ `json:"phase,omitempty" description:"phase is the current lifecycle phase of the servicebroker"`
}

const (
	ServiceBroker_Label = "asiainfo.io/servicebroker"
)

type ServiceBrokerPhase string

const (
	// These are internal finalizer values to Origin
	// Replaced with prdutil.DataFoundryFinalizer()() now.
	//FinalizerDatafoundry kapi.FinalizerName = "openshift.io/origin"

	// ServiceBroker_PhaseNew is create by administrator.
	ServiceBroker_PhaseNew string = /*ServiceBrokerPhase*/ "New"

	// ServiceBroker_PhaseActive indicates that servicebroker service working well.
	ServiceBroker_PhaseActive string = /*ServiceBrokerPhase*/ "Active"

	// ServiceBroker_PhaseFailed indicates that servicebroker stopped.
	ServiceBroker_PhaseFailed string = /*ServiceBrokerPhase*/ "Failed"

	// ServiceBroker_PhaseDeleting indicates that servicebroker is going to be deleted.
	ServiceBroker_PhaseDeleting string = /*ServiceBrokerPhase*/ "Deleting"

	// ServiceBroker_PingTimer indicates that servicebroker last ping time.
	ServiceBroker_PingTimer string = "ServiceBroker/LastPing"

	// ServiceBroker_NewRetryTimes indicates that new servicebroker retry times.
	ServiceBroker_NewRetryTimes string = "ServiceBroker/NewRetryTimes"

	// ServiceBroker_RefreshTimer indicates that servicebroker last refresh backingservice time.
	ServiceBroker_RefreshTimer string = "ServiceBroker/LastRefresh"
)

//=====================================================
//
//=====================================================

// Validate checks that an instance of ServiceBroker is well formed
func (s ServiceBrokerStrategy) Validate(ctx request.Context, obj runtime.Object) field.ErrorList {
	o := obj.(*prd.ServiceBroker)
	//log.Printf("Validating fields for ServiceBroker %s\n", o.Name)

	//errors := field.ErrorList{}

	// perform validation here and add to errors using field.Invalid
	errors := genericvalidation.ValidateObjectMeta(&o.ObjectMeta, s.NamespaceScoped(), prdutil.CheckResourceNameValidity, field.NewPath("metadata"))

	return errors
}

func (ServiceBrokerStrategy) NamespaceScoped() bool { return false }

func (ServiceBrokerStatusStrategy) NamespaceScoped() bool { return false }

// DefaultingFunction sets default ServiceBroker field values
func (ServiceBrokerSchemeFns) DefaultingFunction(o interface{}) {
	obj := o.(*ServiceBroker)
	// set default field values here
	//log.Printf("Defaulting fields for ServiceBroker %s\n", obj.Name)

	if len(obj.Status.Phase) == 0 {
		obj.Status.Phase = ServiceBroker_PhaseNew
	}
}

func (s ServiceBrokerStrategy) PrepareForCreate(ctx request.Context, obj runtime.Object) {
	s.DefaultStorageStrategy.PrepareForCreate(ctx, obj)

	o := obj.(*prd.ServiceBroker)
	//log.Printf("PrepareForCreating for ServiceBroker %s\n", o.Name)

	// If the sb doesn't have a finalizer, k8s default api server will delete the sb when the DELETE API is called.
	// To avoid this happening, each sb will be set a finalizer.
	// k8s default api server will not delete the sb when the DELETE API is called.
	// The metadata.deletionTimestamp field of the bs will be set instead.
	// So controller has a chance to check if there are active bsi(s) existing.
	o.SetFinalizers(append(o.GetFinalizers(), prdutil.DataFoundryFinalizer()))
}
