
/*
@asiainfo.com
*/


package v1

import (
	"log"

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

// BackingServiceInstance
// +k8s:openapi-gen=true
// +resource:path=backingserviceinstances,strategy=BackingServiceInstanceStrategy
// +subresource:request=Binding,path=binding,rest=BindingBackingServiceInstanceREST
type BackingServiceInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackingServiceInstanceSpec   `json:"spec,omitempty"`
	Status BackingServiceInstanceStatus `json:"status,omitempty"`
}

// BackingServiceInstanceSpec defines the desired state of BackingServiceInstance
type BackingServiceInstanceSpec struct {
	// description of an instance.
	InstanceProvisioning `json:"provisioning, omitempty"`
	// description of an user-provided-service
	UserProvidedService `json:"userprovidedservice, omitempty"`
	// bindings of an instance
	Binding []InstanceBinding `json:"binding, omitempty"`
	// binding number of an instance
	Bound int `json:"bound, omitempty"`
	// id of an instance
	InstanceID string `json:"instance_id, omitempty"`
	// tags of an instance
	Tags []string `json:"tags, omitempty"`
}

// BackingServiceInstanceStatus defines the observed state of BackingServiceInstance
type BackingServiceInstanceStatus struct {
	// phase is the current lifecycle phase of the instance
	Phase string/*BackingServiceInstancePhase*/ `json:"phase,omitempty"`
	// action is the action of the instance
	Action string/*BackingServiceInstanceAction*/ `json:"action,omitempty"`
	// updating an instance need an flag.
	Patch string/*BackingServiceInstancePatch*/ `json:"patch,omitempty"`
	// detail of failure.
	Reason string `json:"reason,omitempty"`
	//last operation  of a instance provisioning
	LastOperation *LastOperation `json:"last_operation,omitempty"`
}

// UserProvidedService describe an user-provided-service
type UserProvidedService struct {
	Credentials map[string]string `json:"credentials, omitempty"`
}

// InstanceProvisioning describe an InstanceProvisioning detail
type InstanceProvisioning struct {
	// dashboard url of an instance
	DashboardUrl string `json:"dashboard_url, omitempty"`
	// bs name of an instance
	BackingServiceName string `json:"backingservice_name, omitempty"`
	// bs id of an instance
	BackingServiceSpecID string `json:"backingservice_spec_id, omitempty"`
	// bs plan id of an instance
	BackingServicePlanGuid string `json:"backingservice_plan_guid, omitempty"`
	// bs plan name of an instance
	BackingServicePlanName string `json:"backingservice_plan_name, omitempty"`
	// parameters of an instance
	Parameters map[string]string `json:"parameters, omitempty"`
	// credentials of an instance
	Creds map[string]string `json:"credentials, omitempty"`
	// access of hadoop insance... hardcode hack.
	Accesses map[string][]string `json:"accesses,omitempty"`
}

// InstanceBinding describe an instance binding.
type InstanceBinding struct {
	// bound time of an instance binding
	BoundTime *metav1.Time `json:"bound_time,omitempty"`
	// bind uid of an instance binding
	BindUuid string `json:"bind_uuid, omitempty"`
	// deploymentconfig of an binding.
	BindDeploymentConfig string `json:"bind_deploymentconfig,omitempty"`
	// bind to hadoopuser
	BindHadoopUser string `json:"bind_hadoop_user,omitempty"`
	// credentials of an instance binding
	Credentials map[string]string `json:"credentials, omitempty"`
}

// LastOperation describe last operation of an instance provisioning
type LastOperation struct {
	// state of last operation
	State string `json:"state"`
	// description of last operation
	Description string `json:"description"`
	// async_poll_interval_seconds of a last operation
	AsyncPollIntervalSeconds int `json:"async_poll_interval_seconds, omitempty"`
}

type BackingServiceInstancePhase string
type BackingServiceInstanceAction string
type BackingServiceInstancePatch string

const (
	BackingServiceInstancePhaseProvisioning string = /*BackingServiceInstancePhase*/ "Provisioning"
	BackingServiceInstancePhaseUnbound      string = /*BackingServiceInstancePhase*/ "Unbound"
	BackingServiceInstancePhaseBound        string = /*BackingServiceInstancePhase*/ "Bound"
	BackingServiceInstancePhaseDeleted      string = /*BackingServiceInstancePhase*/ "Deleted"
	BackingServiceInstancePhaseFailure      string = /*BackingServiceInstancePhase*/ "Failure"

	BackingServiceInstancePatchUpdating string = /*BackingServiceInstancePatch*/ "Updating"
	BackingServiceInstancePatchUpdate   string = /*BackingServiceInstancePatch*/ "Update"
	BackingServiceInstancePatchUpdated  string = /*BackingServiceInstancePatch*/ "Updated"
	BackingServiceInstancePatchFailure  string = /*BackingServiceInstancePatch*/ "Failure"

	BackingServiceInstanceActionToBind   string = /*BackingServiceInstanceAction*/ "_ToBind"
	BackingServiceInstanceActionToUnbind string = /*BackingServiceInstanceAction*/ "_ToUnbind"
	BackingServiceInstanceActionToDelete string = /*BackingServiceInstanceAction*/ "_ToDelete"

	BindDeploymentConfigBinding   string = "binding"
	BindDeploymentConfigUnbinding string = "unbinding"
	BindDeploymentConfigBound     string = "bound"

	UPS string = "USER-PROVIDED-SERVICE"
)

//=====================================================
//
//=====================================================

func validateBackingServiceInstanceSpec(spec *prd.BackingServiceInstanceSpec) field.ErrorList {
	allErrs := field.ErrorList{}

	specPath := field.NewPath("spec")
	if spec.BackingServiceName == "" {
		allErrs = append(allErrs, field.Invalid(specPath.Child("backingservice_name"), spec.BackingServiceName, "BackingServiceName must be specified"))
	}
	if spec.BackingServicePlanGuid == "" {
		allErrs = append(allErrs, field.Invalid(specPath.Child("backingservice_plan_guid"), spec.BackingServicePlanGuid, "BackingServicePlanGuid must be specified"))
	}

	return allErrs
}

// Validate checks that an instance of BackingServiceInstance is well formed
func (s BackingServiceInstanceStrategy) Validate(ctx request.Context, obj runtime.Object) field.ErrorList {
	o := obj.(*prd.BackingServiceInstance)
	log.Printf("Validating fields for BackingServiceInstance %s\n", o.Name)

	//errors := field.ErrorList{}
	// perform validation here and add to errors using field.Invalid
	errors := genericvalidation.ValidateObjectMeta(&o.ObjectMeta, s.NamespaceScoped(), prdutil.CheckResourceNameValidity, field.NewPath("metadata"))
	errors = append(errors, validateBackingServiceInstanceSpec(&o.Spec)...)
	return errors
}

// DefaultingFunction sets default BackingServiceInstance field values
func (BackingServiceInstanceSchemeFns) DefaultingFunction(o interface{}) {
	obj := o.(*BackingServiceInstance)
	// set default field values here
	log.Printf("Defaulting fields for BackingServiceInstance %s\n", obj.Name)

	if len(obj.Status.Phase) == 0 {
		obj.Status.Phase = BackingServiceInstancePhaseProvisioning
	}
}
