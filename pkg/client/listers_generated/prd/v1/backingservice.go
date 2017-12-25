/*
@asiainfo.com
*/

// This file was automatically generated by lister-gen

package v1

import (
	v1 "github.com/asiainfoldp/apiserver-servicebroker/pkg/apis/prd/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// BackingServiceLister helps list BackingServices.
type BackingServiceLister interface {
	// List lists all BackingServices in the indexer.
	List(selector labels.Selector) (ret []*v1.BackingService, err error)
	// Get retrieves the BackingService from the index for a given name.
	Get(name string) (*v1.BackingService, error)
	BackingServiceListerExpansion
}

// backingServiceLister implements the BackingServiceLister interface.
type backingServiceLister struct {
	indexer cache.Indexer
}

// NewBackingServiceLister returns a new BackingServiceLister.
func NewBackingServiceLister(indexer cache.Indexer) BackingServiceLister {
	return &backingServiceLister{indexer: indexer}
}

// List lists all BackingServices in the indexer.
func (s *backingServiceLister) List(selector labels.Selector) (ret []*v1.BackingService, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.BackingService))
	})
	return ret, err
}

// Get retrieves the BackingService from the index for a given name.
func (s *backingServiceLister) Get(name string) (*v1.BackingService, error) {
	key := &v1.BackingService{ObjectMeta: meta_v1.ObjectMeta{Name: name}}
	obj, exists, err := s.indexer.Get(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("backingservice"), name)
	}
	return obj.(*v1.BackingService), nil
}
