package lasso

import (
	controllerruntime "github.com/rancher/lasso/controller-runtime"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/lasso/pkg/dynamic"
	"github.com/rancher/wrangler/pkg/apply"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
)

type Runtime struct {
	Apply   apply.Apply
	Backend *Backend
}

func NewRuntime(cfg *rest.Config, scheme *runtime.Scheme) (*Runtime, error) {
	factory, err := controller.NewSharedControllerFactoryFromConfig(cfg, scheme)
	if err != nil {
		return nil, err
	}

	restClient, err := rest.UnversionedRESTClientFor(cfg)
	if err != nil {
		return nil, err
	}

	dc := discovery.NewDiscoveryClient(restClient)
	cache, err := controllerruntime.NewNewCacheFunc(factory.SharedCacheFactory(), dynamic.New(dc))(cfg, cache.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, err
	}

	client, err := cluster.DefaultNewClient(cache, cfg, client.Options{
		Scheme: scheme,
	})

	return &Runtime{
		Apply:   apply.New(dc, apply.NewClientFactory(cfg)),
		Backend: NewBackend(factory, client, cache),
	}, nil
}