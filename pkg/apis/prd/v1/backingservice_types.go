/*
@asiainfo.com
*/

package v1

import (
	genericvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/asiainfoldp/apiserver-servicebroker/pkg/apis/prd"
	prdutil "github.com/asiainfoldp/apiserver-servicebroker/pkg/apis/prd/util"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// BackingService
// +k8s:openapi-gen=true
// +resource:path=backingservices,strategy=BackingServiceStrategy
type BackingService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackingServiceSpec   `json:"spec,omitempty"`
	Status BackingServiceStatus `json:"status,omitempty"`
}

// BackingServiceSpec defines the desired state of BackingService
type BackingServiceSpec struct {
	// name of backingservice
	Name string `json:"name" description:"name of backingservice"`
	// id of backingservice
	Id string `json:"id" description:"id of backingservice"`
	// description of a backingservice
	Description string `json:"description" description:"description of a backingservice"`
	// is backingservice bindable
	Bindable bool `json:"bindable" description:"is backingservice bindable?"`
	// is  backingservice plan updateable
	PlanUpdateable bool `json:"plan_updateable, omitempty" description:"is  backingservice plan updateable"`
	// list of backingservice tags of BackingService
	Tags []string `json:"tags, omitempty" description:"list of backingservice tags of BackingService"`
	// require condition of backingservice
	Requires []string `json:"requires, omitempty" description:"require condition of backingservice"`

	// metadata of backingservice
	Metadata map[string]string `json:"metadata, omitempty" description:"metadata of backingservice"`
	// plans of a backingservice
	Plans []ServicePlan `json:"plans" description:"plans of a backingservice"`
	// DashboardClient of backingservic
	DashboardClient map[string]string `json:"dashboard_client" description:"DashboardClient of backingservice"`
}

// BackingServiceStatus defines the observed state of BackingService
type BackingServiceStatus struct {
	// phase is the current lifecycle phase of the servicebroker
	Phase string /*BackingServicePhase*/ `json:"phase,omitempty" description:"phase is the current lifecycle phase of the servicebroker"`
}

// ServicePlan describe a ServicePlan
type ServicePlan struct {
	// name of a ServicePlan
	Name string `json:"name"`
	//id of a ServicePlan
	Id string `json:"id"`
	// description of a ServicePlan
	Description string `json:"description"`
	// metadata of a ServicePlan
	Metadata ServicePlanMetadata `json:"metadata, omitempty"`
	// is this plan free or not
	Free bool `json:"free, omitempty"`
}

// ServicePlanMetadata describe a ServicePlanMetadata
type ServicePlanMetadata struct {
	// bullets of a ServicePlanMetadata
	Bullets []string `json:"bullets, omitempty"`
	// costs of a ServicePlanMetadata
	Costs []ServicePlanCost `json:"costs, omitempty"`
	// displayname of a ServicePlanMetadata
	DisplayName string `json:"displayName, omitempty"`
	// customize options of a plan.
	Customize map[string]CustomizeSpec `json:"customize, omitempty"`
}

type CustomizeSpec struct {
	Default float64 `json:"default,omitempty"`
	Max     float64 `json:"max,omitempty"`
	Price   float64 `json:"price,omitempty"`
	Step    float64 `json:"step,omitempty"`
	Unit    string  `json:"unit,omitempty"`
	Desc    string  `json:"desc,omitempty"`
}

//TODO amount should be a array object...

// ServicePlanCost describe a ServicePlanCost
type ServicePlanCost struct {
	// amount of a ServicePlanCost
	Amount map[string]float64 `json:"amount, omitempty"`
	// unit of a ServicePlanCost
	Unit string `json:"unit, omitempty"`
}

const BackingService_NS = "openshift"

type BackingServicePhase string

const (
	BackingService_PhaseActive   string = /*BackingServicePhase*/ "Active"
	BackingService_PhaseInactive string = /*BackingServicePhase*/ "Inactive"
)

//=====================================================
//
//=====================================================

// Validate checks that an instance of BackingService is well formed
func (s BackingServiceStrategy) Validate(ctx request.Context, obj runtime.Object) field.ErrorList {
	o := obj.(*prd.BackingService)
	//log.Printf("Validating fields for BackingService %s\n", o.Name)

	//errors := field.ErrorList{}
	// perform validation here and add to errors using field.Invalid
	errors := genericvalidation.ValidateObjectMeta(&o.ObjectMeta, s.NamespaceScoped(), prdutil.CheckResourceNameValidity, field.NewPath("metadata"))
	return errors
}

func (BackingServiceStrategy) NamespaceScoped() bool { return false }

func (BackingServiceStatusStrategy) NamespaceScoped() bool { return false }

// DefaultingFunction sets default BackingService field values
func (BackingServiceSchemeFns) DefaultingFunction(o interface{}) {
	//obj := o.(*BackingService)
	// set default field values here
	//log.Printf("Defaulting fields for BackingService %s\n", obj.Name)
}
