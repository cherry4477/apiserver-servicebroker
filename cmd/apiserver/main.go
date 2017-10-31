
/*
@asiainfo.com
*/


package main

import (
	// Make sure glide gets these dependencies
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "github.com/go-openapi/loads"

	"github.com/kubernetes-incubator/apiserver-builder/pkg/cmd/server"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // Enable cloud provider auth

	"github.com/asiainfoldp/apiserver-servicebroker/pkg/apis"
	"github.com/asiainfoldp/apiserver-servicebroker/pkg/openapi"
)

func main() {
	version := "v0"
	server.StartApiServer("/registry/asiainfo.com", apis.GetAllApiBuilders(), openapi.GetOpenAPIDefinitions, "Api", version)
}
