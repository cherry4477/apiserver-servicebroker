
/*
@asiainfo.com
*/


package v1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/asiainfoldp/apiserver-servicebroker/pkg/apis/prd"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +subresource-request
type Binding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// bind kind is bindking of an instance binding
	BindKind string `json:"bindKind, omitempty"`
	// bindResourceVersion is bindResourceVersion of an instance binding.
	BindResourceVersion string `json:"bindResourceVersion, omitempty"`
	// resourceName of an instance binding
	ResourceName string `json:"resourceName, omitempty"`
	// parameters of a binding request
	Parameters map[string]string `json:"parameters, omitempty"`
}

//=====================================================
//
//=====================================================

const (
	BindKind_DeploymentConfig = "DeploymentConfig"
	BindKind_HadoopUser       = "HadoopUser"
)

//=====================================================
//
//=====================================================

var _ rest.CreaterUpdater = &BindingBackingServiceInstanceREST{}
var _ rest.Patcher = &BindingBackingServiceInstanceREST{}

// +k8s:deepcopy-gen=false
type BindingBackingServiceInstanceREST struct {
	Registry prd.BackingServiceInstanceRegistry
}

func (r *BindingBackingServiceInstanceREST) New() runtime.Object {
	return &Binding{}
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *BindingBackingServiceInstanceREST) Get(ctx request.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return nil, nil
}

func (r *BindingBackingServiceInstanceREST) Create(ctx request.Context, obj runtime.Object, includeUninitialized bool) (runtime.Object, error) {
	bro, ok := obj.(*Binding) // ok may be always true
	if !ok {
		return nil, fmt.Errorf("input is not a bind")
	}
	if bro.BindKind != BindKind_DeploymentConfig && bro.BindKind != BindKind_HadoopUser {
		return nil, fmt.Errorf("unsupported bind type: %s", bro.BindKind)
	}

	bsi, err := r.Registry.GetBackingServiceInstance(ctx, bro.Name, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if bsi.Annotations == nil {
		bsi.Annotations = map[string]string{}
	}
	if bro.BindKind == BindKind_DeploymentConfig {

		if bound := bsi.Annotations[bro.ResourceName]; bound == BindDeploymentConfigBound {
			return nil, fmt.Errorf("'%s' is already bound to this instance.", bro.ResourceName)
		}
		//if bsi.Status.Phase != BackingServiceInstancePhaseUnbound {
		//	return nil, errors.New("back service instance is not in unbound phase")
		//}

		//bs, err := r.Registry.GetBackingService(ctx, bro.Name, &metav1.GetOptions{})
		//if err != nil {
		//	return nil, err
		//}

		/*
			dc, err := r.Registry.GetDeploymentConfig(ctx, bro.Name, &metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			_ = dc
		*/

		// update bsi

		//need debug....bsi.Spec.BindDeploymentConfig = bro.ResourceName // dc.Name

		bsi.Annotations[bro.ResourceName] = BindDeploymentConfigBinding
		bsi.Status.Action = BackingServiceInstanceActionToBind
	} else {
		for _, bd := range bsi.Spec.Binding {
			if bd.BindHadoopUser == bro.ResourceName {
				return nil, fmt.Errorf(" user '%s' is already bound to this instance.", bro.ResourceName)
			}
		}

		if bro.Parameters != nil {
			if bsi.Spec.Parameters == nil {
				bsi.Spec.Parameters = make(map[string]string)
			}
			for k, v := range bro.Parameters {
				bsi.Spec.Parameters[k] = v
			}
		}
		bsi.Annotations[BindKind_HadoopUser] = bro.ResourceName
		bsi.Status.Action = BackingServiceInstanceActionToBind

	}

	bsi, err = r.Registry.UpdateBackingServiceInstance(ctx, bsi)
	if err != nil {
		return nil, err
	}

	return bsi, nil
}

// Update alters the status subset of an object.
func (r *BindingBackingServiceInstanceREST) Update(ctx request.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	bsi, err := r.Registry.GetBackingServiceInstance(ctx, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}

	var oldBind *Binding                            // nil
	obj, err := objInfo.UpdatedObject(ctx, oldBind) // make a diff from old and input
	if obj == nil {
		return nil, false, err
	}

	bro, ok := obj.(*Binding) // ok may be always true
	if !ok {
		return nil, false, fmt.Errorf("input is not a bind")
	}
	if bro.BindKind != BindKind_DeploymentConfig && bro.BindKind != BindKind_HadoopUser {
		return nil, false, fmt.Errorf("unsupported bind type: '%s'", bro.BindKind)
	}

	if bsi.Annotations == nil {
		bsi.Annotations = map[string]string{}
	}

	if bro.BindKind == BindKind_DeploymentConfig {

		if bound, ok := bsi.Annotations[bro.ResourceName]; !ok || bound == "unbound" /*unbound should never happen.*/ {
			return nil, false, fmt.Errorf("DeploymentConfig '%s' not bound to this instance yet.", bro.ResourceName)
		} else {
			bsi.Annotations[bro.ResourceName] = BindDeploymentConfigUnbinding
			bsi.Status.Action = BackingServiceInstanceActionToUnbind
		}
	} else {
		bsi.Annotations[BindKind_HadoopUser] = bro.ResourceName
		bsi.Status.Action = BackingServiceInstanceActionToUnbind
	}

	bsi, err = r.Registry.UpdateBackingServiceInstance(ctx, bsi)
	if err != nil {
		return nil, false, err
	}
	return bsi, true, nil
}
