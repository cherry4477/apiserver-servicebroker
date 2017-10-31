
/*
@asiainfo.com
*/


// Api versions allow the api contract for a resource to be changed while keeping
// backward compatibility by support multiple concurrent versions
// of the same resource

// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=github.com/asiainfoldp/apiserver-servicebroker/pkg/apis/prd
// +k8s:defaulter-gen=TypeMeta
// +groupName=prd.asiainfo.com
package v1 // import "github.com/asiainfoldp/apiserver-servicebroker/pkg/apis/prd/v1"

