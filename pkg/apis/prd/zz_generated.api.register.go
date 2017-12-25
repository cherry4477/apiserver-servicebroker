/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file was autogenerated by apiregister-gen. Do not edit it manually!

package prd

import (
	"fmt"
	"github.com/kubernetes-incubator/apiserver-builder/pkg/builders"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
)

var (
	InternalBackingService = builders.NewInternalResource(
		"backingservices",
		func() runtime.Object { return &BackingService{} },
		func() runtime.Object { return &BackingServiceList{} },
	)
	InternalBackingServiceStatus = builders.NewInternalResourceStatus(
		"backingservices",
		func() runtime.Object { return &BackingService{} },
		func() runtime.Object { return &BackingServiceList{} },
	)
	InternalBackingServiceInstance = builders.NewInternalResource(
		"backingserviceinstances",
		func() runtime.Object { return &BackingServiceInstance{} },
		func() runtime.Object { return &BackingServiceInstanceList{} },
	)
	InternalBackingServiceInstanceStatus = builders.NewInternalResourceStatus(
		"backingserviceinstances",
		func() runtime.Object { return &BackingServiceInstance{} },
		func() runtime.Object { return &BackingServiceInstanceList{} },
	)
	InternalBindingBackingServiceInstanceREST = builders.NewInternalSubresource(
		"backingserviceinstances", "binding",
		func() runtime.Object { return &Binding{} },
	)
	InternalServiceBroker = builders.NewInternalResource(
		"servicebrokers",
		func() runtime.Object { return &ServiceBroker{} },
		func() runtime.Object { return &ServiceBrokerList{} },
	)
	InternalServiceBrokerStatus = builders.NewInternalResourceStatus(
		"servicebrokers",
		func() runtime.Object { return &ServiceBroker{} },
		func() runtime.Object { return &ServiceBrokerList{} },
	)
	// Registered resources and subresources
	ApiVersion = builders.NewApiGroup("prd.asiainfo.com").WithKinds(
		InternalBackingService,
		InternalBackingServiceStatus,
		InternalBackingServiceInstance,
		InternalBackingServiceInstanceStatus,
		InternalBindingBackingServiceInstanceREST,
		InternalServiceBroker,
		InternalServiceBrokerStatus,
	)

	// Required by code generated by go2idl
	AddToScheme        = ApiVersion.SchemaBuilder.AddToScheme
	SchemeBuilder      = ApiVersion.SchemaBuilder
	localSchemeBuilder = &SchemeBuilder
	SchemeGroupVersion = ApiVersion.GroupVersion
)

// Required by code generated by go2idl
// Kind takes an unqualified kind and returns a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Required by code generated by go2idl
// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// +genclient
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BackingServiceInstance struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   BackingServiceInstanceSpec
	Status BackingServiceInstanceStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Binding struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	BindKind            string
	BindResourceVersion string
	ResourceName        string
	Parameters          map[string]string
}

type BackingServiceInstanceStatus struct {
	Phase         string
	Action        string
	Patch         string
	Reason        string
	LastOperation *LastOperation
}

type BackingServiceInstanceSpec struct {
	InstanceProvisioning
	UserProvidedService
	Binding    []InstanceBinding
	Bound      int
	InstanceID string
	Tags       []string
}

type LastOperation struct {
	State                    string
	Description              string
	AsyncPollIntervalSeconds int
}

type InstanceBinding struct {
	BoundTime            *meta.Time
	BindUuid             string
	BindDeploymentConfig string
	BindHadoopUser       string
	Credentials          map[string]string
}

type UserProvidedService struct {
	Credentials map[string]string
}

type InstanceProvisioning struct {
	DashboardUrl           string
	BackingServiceName     string
	BackingServiceSpecID   string
	BackingServicePlanGuid string
	BackingServicePlanName string
	Parameters             map[string]string
	Creds                  map[string]string
	Accesses               map[string][]string
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BackingService struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   BackingServiceSpec
	Status BackingServiceStatus
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ServiceBroker struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   ServiceBrokerSpec
	Status ServiceBrokerStatus
}

type BackingServiceStatus struct {
	Phase string
}

type ServiceBrokerStatus struct {
	Phase string
}

type ServiceBrokerSpec struct {
	Url        string
	Name       string
	UserName   string
	Password   string
	Finalizers []corev1.FinalizerName
}

type BackingServiceSpec struct {
	Name            string
	Id              string
	Description     string
	Bindable        bool
	PlanUpdateable  bool
	Tags            []string
	Requires        []string
	Metadata        map[string]string
	Plans           []ServicePlan
	DashboardClient map[string]string
}

type ServicePlan struct {
	Name        string
	Id          string
	Description string
	Metadata    ServicePlanMetadata
	Free        bool
}

type ServicePlanMetadata struct {
	Bullets     []string
	Costs       []ServicePlanCost
	DisplayName string
	Customize   map[string]CustomizeSpec
}

type ServicePlanCost struct {
	Amount map[string]float64
	Unit   string
}

type CustomizeSpec struct {
	Default float64
	Max     float64
	Price   float64
	Step    float64
	Unit    string
	Desc    string
}

//
// BackingService Functions and Structs
//
// +k8s:deepcopy-gen=false
type BackingServiceStrategy struct {
	builders.DefaultStorageStrategy
}

// +k8s:deepcopy-gen=false
type BackingServiceStatusStrategy struct {
	builders.DefaultStatusStorageStrategy
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BackingServiceList struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []BackingService
}

func (BackingService) NewStatus() interface{} {
	return BackingServiceStatus{}
}

func (pc *BackingService) GetStatus() interface{} {
	return pc.Status
}

func (pc *BackingService) SetStatus(s interface{}) {
	pc.Status = s.(BackingServiceStatus)
}

func (pc *BackingService) GetSpec() interface{} {
	return pc.Spec
}

func (pc *BackingService) SetSpec(s interface{}) {
	pc.Spec = s.(BackingServiceSpec)
}

func (pc *BackingService) GetObjectMeta() *metav1.ObjectMeta {
	return &pc.ObjectMeta
}

func (pc *BackingService) SetGeneration(generation int64) {
	pc.ObjectMeta.Generation = generation
}

func (pc BackingService) GetGeneration() int64 {
	return pc.ObjectMeta.Generation
}

// Registry is an interface for things that know how to store BackingService.
// +k8s:deepcopy-gen=false
type BackingServiceRegistry interface {
	ListBackingServices(ctx request.Context, options *internalversion.ListOptions) (*BackingServiceList, error)
	GetBackingService(ctx request.Context, id string, options *metav1.GetOptions) (*BackingService, error)
	CreateBackingService(ctx request.Context, id *BackingService) (*BackingService, error)
	UpdateBackingService(ctx request.Context, id *BackingService) (*BackingService, error)
	DeleteBackingService(ctx request.Context, id string) (bool, error)
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched types will panic.
func NewBackingServiceRegistry(sp builders.StandardStorageProvider) BackingServiceRegistry {
	return &storageBackingService{sp}
}

// Implement Registry
// storage puts strong typing around storage calls
// +k8s:deepcopy-gen=false
type storageBackingService struct {
	builders.StandardStorageProvider
}

func (s *storageBackingService) ListBackingServices(ctx request.Context, options *internalversion.ListOptions) (*BackingServiceList, error) {
	if options != nil && options.FieldSelector != nil && !options.FieldSelector.Empty() {
		return nil, fmt.Errorf("field selector not supported yet")
	}
	st := s.GetStandardStorage()
	obj, err := st.List(ctx, options)
	if err != nil {
		return nil, err
	}
	return obj.(*BackingServiceList), err
}

func (s *storageBackingService) GetBackingService(ctx request.Context, id string, options *metav1.GetOptions) (*BackingService, error) {
	st := s.GetStandardStorage()
	obj, err := st.Get(ctx, id, options)
	if err != nil {
		return nil, err
	}
	return obj.(*BackingService), nil
}

func (s *storageBackingService) CreateBackingService(ctx request.Context, object *BackingService) (*BackingService, error) {
	st := s.GetStandardStorage()
	obj, err := st.Create(ctx, object, false)
	if err != nil {
		return nil, err
	}
	return obj.(*BackingService), nil
}

func (s *storageBackingService) UpdateBackingService(ctx request.Context, object *BackingService) (*BackingService, error) {
	st := s.GetStandardStorage()
	obj, _, err := st.Update(ctx, object.Name, rest.DefaultUpdatedObjectInfo(object, builders.Scheme))
	if err != nil {
		return nil, err
	}
	return obj.(*BackingService), nil
}

func (s *storageBackingService) DeleteBackingService(ctx request.Context, id string) (bool, error) {
	st := s.GetStandardStorage()
	_, sync, err := st.Delete(ctx, id, nil)
	return sync, err
}

//
// BackingServiceInstance Functions and Structs
//
// +k8s:deepcopy-gen=false
type BackingServiceInstanceStrategy struct {
	builders.DefaultStorageStrategy
}

// +k8s:deepcopy-gen=false
type BackingServiceInstanceStatusStrategy struct {
	builders.DefaultStatusStorageStrategy
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BackingServiceInstanceList struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []BackingServiceInstance
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BindingList struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []Binding
}

func (BackingServiceInstance) NewStatus() interface{} {
	return BackingServiceInstanceStatus{}
}

func (pc *BackingServiceInstance) GetStatus() interface{} {
	return pc.Status
}

func (pc *BackingServiceInstance) SetStatus(s interface{}) {
	pc.Status = s.(BackingServiceInstanceStatus)
}

func (pc *BackingServiceInstance) GetSpec() interface{} {
	return pc.Spec
}

func (pc *BackingServiceInstance) SetSpec(s interface{}) {
	pc.Spec = s.(BackingServiceInstanceSpec)
}

func (pc *BackingServiceInstance) GetObjectMeta() *metav1.ObjectMeta {
	return &pc.ObjectMeta
}

func (pc *BackingServiceInstance) SetGeneration(generation int64) {
	pc.ObjectMeta.Generation = generation
}

func (pc BackingServiceInstance) GetGeneration() int64 {
	return pc.ObjectMeta.Generation
}

// Registry is an interface for things that know how to store BackingServiceInstance.
// +k8s:deepcopy-gen=false
type BackingServiceInstanceRegistry interface {
	ListBackingServiceInstances(ctx request.Context, options *internalversion.ListOptions) (*BackingServiceInstanceList, error)
	GetBackingServiceInstance(ctx request.Context, id string, options *metav1.GetOptions) (*BackingServiceInstance, error)
	CreateBackingServiceInstance(ctx request.Context, id *BackingServiceInstance) (*BackingServiceInstance, error)
	UpdateBackingServiceInstance(ctx request.Context, id *BackingServiceInstance) (*BackingServiceInstance, error)
	DeleteBackingServiceInstance(ctx request.Context, id string) (bool, error)
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched types will panic.
func NewBackingServiceInstanceRegistry(sp builders.StandardStorageProvider) BackingServiceInstanceRegistry {
	return &storageBackingServiceInstance{sp}
}

// Implement Registry
// storage puts strong typing around storage calls
// +k8s:deepcopy-gen=false
type storageBackingServiceInstance struct {
	builders.StandardStorageProvider
}

func (s *storageBackingServiceInstance) ListBackingServiceInstances(ctx request.Context, options *internalversion.ListOptions) (*BackingServiceInstanceList, error) {
	if options != nil && options.FieldSelector != nil && !options.FieldSelector.Empty() {
		return nil, fmt.Errorf("field selector not supported yet")
	}
	st := s.GetStandardStorage()
	obj, err := st.List(ctx, options)
	if err != nil {
		return nil, err
	}
	return obj.(*BackingServiceInstanceList), err
}

func (s *storageBackingServiceInstance) GetBackingServiceInstance(ctx request.Context, id string, options *metav1.GetOptions) (*BackingServiceInstance, error) {
	st := s.GetStandardStorage()
	obj, err := st.Get(ctx, id, options)
	if err != nil {
		return nil, err
	}
	return obj.(*BackingServiceInstance), nil
}

func (s *storageBackingServiceInstance) CreateBackingServiceInstance(ctx request.Context, object *BackingServiceInstance) (*BackingServiceInstance, error) {
	st := s.GetStandardStorage()
	obj, err := st.Create(ctx, object, false)
	if err != nil {
		return nil, err
	}
	return obj.(*BackingServiceInstance), nil
}

func (s *storageBackingServiceInstance) UpdateBackingServiceInstance(ctx request.Context, object *BackingServiceInstance) (*BackingServiceInstance, error) {
	st := s.GetStandardStorage()
	obj, _, err := st.Update(ctx, object.Name, rest.DefaultUpdatedObjectInfo(object, builders.Scheme))
	if err != nil {
		return nil, err
	}
	return obj.(*BackingServiceInstance), nil
}

func (s *storageBackingServiceInstance) DeleteBackingServiceInstance(ctx request.Context, id string) (bool, error) {
	st := s.GetStandardStorage()
	_, sync, err := st.Delete(ctx, id, nil)
	return sync, err
}

//
// ServiceBroker Functions and Structs
//
// +k8s:deepcopy-gen=false
type ServiceBrokerStrategy struct {
	builders.DefaultStorageStrategy
}

// +k8s:deepcopy-gen=false
type ServiceBrokerStatusStrategy struct {
	builders.DefaultStatusStorageStrategy
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ServiceBrokerList struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []ServiceBroker
}

func (ServiceBroker) NewStatus() interface{} {
	return ServiceBrokerStatus{}
}

func (pc *ServiceBroker) GetStatus() interface{} {
	return pc.Status
}

func (pc *ServiceBroker) SetStatus(s interface{}) {
	pc.Status = s.(ServiceBrokerStatus)
}

func (pc *ServiceBroker) GetSpec() interface{} {
	return pc.Spec
}

func (pc *ServiceBroker) SetSpec(s interface{}) {
	pc.Spec = s.(ServiceBrokerSpec)
}

func (pc *ServiceBroker) GetObjectMeta() *metav1.ObjectMeta {
	return &pc.ObjectMeta
}

func (pc *ServiceBroker) SetGeneration(generation int64) {
	pc.ObjectMeta.Generation = generation
}

func (pc ServiceBroker) GetGeneration() int64 {
	return pc.ObjectMeta.Generation
}

// Registry is an interface for things that know how to store ServiceBroker.
// +k8s:deepcopy-gen=false
type ServiceBrokerRegistry interface {
	ListServiceBrokers(ctx request.Context, options *internalversion.ListOptions) (*ServiceBrokerList, error)
	GetServiceBroker(ctx request.Context, id string, options *metav1.GetOptions) (*ServiceBroker, error)
	CreateServiceBroker(ctx request.Context, id *ServiceBroker) (*ServiceBroker, error)
	UpdateServiceBroker(ctx request.Context, id *ServiceBroker) (*ServiceBroker, error)
	DeleteServiceBroker(ctx request.Context, id string) (bool, error)
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched types will panic.
func NewServiceBrokerRegistry(sp builders.StandardStorageProvider) ServiceBrokerRegistry {
	return &storageServiceBroker{sp}
}

// Implement Registry
// storage puts strong typing around storage calls
// +k8s:deepcopy-gen=false
type storageServiceBroker struct {
	builders.StandardStorageProvider
}

func (s *storageServiceBroker) ListServiceBrokers(ctx request.Context, options *internalversion.ListOptions) (*ServiceBrokerList, error) {
	if options != nil && options.FieldSelector != nil && !options.FieldSelector.Empty() {
		return nil, fmt.Errorf("field selector not supported yet")
	}
	st := s.GetStandardStorage()
	obj, err := st.List(ctx, options)
	if err != nil {
		return nil, err
	}
	return obj.(*ServiceBrokerList), err
}

func (s *storageServiceBroker) GetServiceBroker(ctx request.Context, id string, options *metav1.GetOptions) (*ServiceBroker, error) {
	st := s.GetStandardStorage()
	obj, err := st.Get(ctx, id, options)
	if err != nil {
		return nil, err
	}
	return obj.(*ServiceBroker), nil
}

func (s *storageServiceBroker) CreateServiceBroker(ctx request.Context, object *ServiceBroker) (*ServiceBroker, error) {
	st := s.GetStandardStorage()
	obj, err := st.Create(ctx, object, false)
	if err != nil {
		return nil, err
	}
	return obj.(*ServiceBroker), nil
}

func (s *storageServiceBroker) UpdateServiceBroker(ctx request.Context, object *ServiceBroker) (*ServiceBroker, error) {
	st := s.GetStandardStorage()
	obj, _, err := st.Update(ctx, object.Name, rest.DefaultUpdatedObjectInfo(object, builders.Scheme))
	if err != nil {
		return nil, err
	}
	return obj.(*ServiceBroker), nil
}

func (s *storageServiceBroker) DeleteServiceBroker(ctx request.Context, id string) (bool, error) {
	st := s.GetStandardStorage()
	_, sync, err := st.Delete(ctx, id, nil)
	return sync, err
}