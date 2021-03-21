package k8s

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

type K8sServer struct {
	KubeClientset *kubernetes.Clientset
	RestConfig    *rest.Config

	groupVersionKinds []*schema.GroupVersionKind
	exts              []*runtime.RawExtension
}

func NewK8sServer(config *rest.Config) (*K8sServer, error) {
	kubeClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	config.ContentType = runtime.ContentTypeJSON
	config.APIPath = "/apis"
	//config.NegotiatedSerializer =
	return &K8sServer{
		KubeClientset: kubeClientset,
		RestConfig:    config,
	}, nil
}

func (k *K8sServer) CreateWithFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	err = k.Decode(file)
	if err != nil {
		return err
	}
	for _, ext := range k.exts {
		err = k.CreateObject(ext.Object)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k *K8sServer) Decode(r io.Reader) error {
	decoder := yaml.NewYAMLOrJSONDecoder(r, 4096)
	for {
		ext := runtime.RawExtension{}
		if err := decoder.Decode(&ext); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("error parsing: %v", err)
		}
		ext.Raw = bytes.TrimSpace(ext.Raw)
		if len(ext.Raw) == 0 || bytes.Equal(ext.Raw, []byte("null")) {
			continue
		}
		obj, gkv, err := unstructured.UnstructuredJSONScheme.Decode(ext.Raw, nil, nil)
		if err != nil {
			return err
		}
		ext.Object = obj

		k.groupVersionKinds = append(k.groupVersionKinds, gkv)
		k.exts = append(k.exts, &ext)
	}
}

func (k *K8sServer) CreateObject(obj runtime.Object) error {
	groupResources, err := restmapper.GetAPIGroupResources(k.KubeClientset.Discovery())
	if err != nil {
		return err
	}
	rm := restmapper.NewDiscoveryRESTMapper(groupResources)

	gvk := obj.GetObjectKind().GroupVersionKind()
	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}
	mapping, err := rm.RESTMapping(gk, gvk.Version)
	if err != nil {
		return err
	}

	restClient, err := k.newRestClient(*k.RestConfig, schema.GroupVersion{
		Group:   gvk.Group,
		Version: gvk.Version,
	})

	restHelper := resource.NewHelper(restClient, mapping)
	_, err = restHelper.Create("default", true, obj)
	return err
}

func (k *K8sServer) newRestClient(restConfig rest.Config, gv schema.GroupVersion) (rest.Interface, error) {
	restConfig.ContentConfig = resource.UnstructuredPlusDefaultContentConfig()
	restConfig.GroupVersion = &gv
	if len(gv.Group) == 0 {
		restConfig.APIPath = "/api"
	} else {
		restConfig.APIPath = "/apis"
	}

	return rest.RESTClientFor(&restConfig)
}
