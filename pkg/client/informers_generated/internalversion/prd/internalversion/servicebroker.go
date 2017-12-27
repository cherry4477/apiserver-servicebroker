/*
@asiainfo.com
*/

// This file was automatically generated by informer-gen

package internalversion

import (
	prd "github.com/asiainfoldp/apiserver-servicebroker/pkg/apis/prd"
	internalclientset "github.com/asiainfoldp/apiserver-servicebroker/pkg/client/clientset_generated/internalclientset"
	internalinterfaces "github.com/asiainfoldp/apiserver-servicebroker/pkg/client/informers_generated/internalversion/internalinterfaces"
	internalversion "github.com/asiainfoldp/apiserver-servicebroker/pkg/client/listers_generated/prd/internalversion"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	time "time"
)

// ServiceBrokerInformer provides access to a shared informer and lister for
// ServiceBrokers.
type ServiceBrokerInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() internalversion.ServiceBrokerLister
}

type serviceBrokerInformer struct {
	factory internalinterfaces.SharedInformerFactory
}

// NewServiceBrokerInformer constructs a new informer for ServiceBroker type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewServiceBrokerInformer(client internalclientset.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				return client.Prd().ServiceBrokers().List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				return client.Prd().ServiceBrokers().Watch(options)
			},
		},
		&prd.ServiceBroker{},
		resyncPeriod,
		indexers,
	)
}

func defaultServiceBrokerInformer(client internalclientset.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewServiceBrokerInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
}

func (f *serviceBrokerInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&prd.ServiceBroker{}, defaultServiceBrokerInformer)
}

func (f *serviceBrokerInformer) Lister() internalversion.ServiceBrokerLister {
	return internalversion.NewServiceBrokerLister(f.Informer().GetIndexer())
}
