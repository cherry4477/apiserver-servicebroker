package util

import (
	//kapi "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/validation/path"
)

// CheckResourceNameValidity check the validity of the name of a new created resource,
// such as ServiceBroker, BackingService and Backing
func CheckResourceNameValidity(name string, prefix bool) []string {
	errors := path.ValidatePathSegmentName(name, prefix)

	if len(name) < 2 {
		errors = append(errors, "must be at least 2 characters long")
	}

	return errors
}

func DataFoundryFinalizer() string /*kapi.FinalizerName*/ {
	//return "openshift.io/origin" // not a good name, just for history compatibility
	// looks the above one in old df code base has no references.
	return "finalizer.datafoundry"
}
