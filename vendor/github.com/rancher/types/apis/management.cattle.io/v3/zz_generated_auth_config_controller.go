package v3

import (
	"context"

	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/controller"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

var (
	AuthConfigGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "AuthConfig",
	}
	AuthConfigResource = metav1.APIResource{
		Name:         "authconfigs",
		SingularName: "authconfig",
		Namespaced:   false,
		Kind:         AuthConfigGroupVersionKind.Kind,
	}
)

type AuthConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AuthConfig
}

type AuthConfigHandlerFunc func(key string, obj *AuthConfig) error

type AuthConfigLister interface {
	List(namespace string, selector labels.Selector) (ret []*AuthConfig, err error)
	Get(namespace, name string) (*AuthConfig, error)
}

type AuthConfigController interface {
	Informer() cache.SharedIndexInformer
	Lister() AuthConfigLister
	AddHandler(handler AuthConfigHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type AuthConfigInterface interface {
	ObjectClient() *clientbase.ObjectClient
	Create(*AuthConfig) (*AuthConfig, error)
	GetNamespace(name, namespace string, opts metav1.GetOptions) (*AuthConfig, error)
	Get(name string, opts metav1.GetOptions) (*AuthConfig, error)
	Update(*AuthConfig) (*AuthConfig, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespace(name, namespace string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*AuthConfigList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() AuthConfigController
	AddSyncHandler(sync AuthConfigHandlerFunc)
	AddLifecycle(name string, lifecycle AuthConfigLifecycle)
}

type authConfigLister struct {
	controller *authConfigController
}

func (l *authConfigLister) List(namespace string, selector labels.Selector) (ret []*AuthConfig, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*AuthConfig))
	})
	return
}

func (l *authConfigLister) Get(namespace, name string) (*AuthConfig, error) {
	var key string
	if namespace != "" {
		key = namespace + "/" + name
	} else {
		key = name
	}
	obj, exists, err := l.controller.Informer().GetIndexer().GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(schema.GroupResource{
			Group:    AuthConfigGroupVersionKind.Group,
			Resource: "authConfig",
		}, name)
	}
	return obj.(*AuthConfig), nil
}

type authConfigController struct {
	controller.GenericController
}

func (c *authConfigController) Lister() AuthConfigLister {
	return &authConfigLister{
		controller: c,
	}
}

func (c *authConfigController) AddHandler(handler AuthConfigHandlerFunc) {
	c.GenericController.AddHandler(func(key string) error {
		obj, exists, err := c.Informer().GetStore().GetByKey(key)
		if err != nil {
			return err
		}
		if !exists {
			return handler(key, nil)
		}
		return handler(key, obj.(*AuthConfig))
	})
}

type authConfigFactory struct {
}

func (c authConfigFactory) Object() runtime.Object {
	return &AuthConfig{}
}

func (c authConfigFactory) List() runtime.Object {
	return &AuthConfigList{}
}

func (s *authConfigClient) Controller() AuthConfigController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.authConfigControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(AuthConfigGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &authConfigController{
		GenericController: genericController,
	}

	s.client.authConfigControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type authConfigClient struct {
	client       *Client
	ns           string
	objectClient *clientbase.ObjectClient
	controller   AuthConfigController
}

func (s *authConfigClient) ObjectClient() *clientbase.ObjectClient {
	return s.objectClient
}

func (s *authConfigClient) Create(o *AuthConfig) (*AuthConfig, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*AuthConfig), err
}

func (s *authConfigClient) Get(name string, opts metav1.GetOptions) (*AuthConfig, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*AuthConfig), err
}

func (s *authConfigClient) GetNamespace(name, namespace string, opts metav1.GetOptions) (*AuthConfig, error) {
	obj, err := s.objectClient.GetNamespace(name, namespace, opts)
	return obj.(*AuthConfig), err
}

func (s *authConfigClient) Update(o *AuthConfig) (*AuthConfig, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*AuthConfig), err
}

func (s *authConfigClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *authConfigClient) DeleteNamespace(name, namespace string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespace(name, namespace, options)
}

func (s *authConfigClient) List(opts metav1.ListOptions) (*AuthConfigList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*AuthConfigList), err
}

func (s *authConfigClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *authConfigClient) Patch(o *AuthConfig, data []byte, subresources ...string) (*AuthConfig, error) {
	obj, err := s.objectClient.Patch(o.Name, o, data, subresources...)
	return obj.(*AuthConfig), err
}

func (s *authConfigClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *authConfigClient) AddSyncHandler(sync AuthConfigHandlerFunc) {
	s.Controller().AddHandler(sync)
}

func (s *authConfigClient) AddLifecycle(name string, lifecycle AuthConfigLifecycle) {
	sync := NewAuthConfigLifecycleAdapter(name, s, lifecycle)
	s.AddSyncHandler(sync)
}
