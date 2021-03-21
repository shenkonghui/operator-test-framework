package k8s

import (
	"fmt"
	"testing"

	"k8s.io/client-go/tools/clientcmd"
)

func TestK8sServer_CreateObject(t *testing.T) {

	config, _ := clientcmd.BuildConfigFromFlags("", "/Users/shenkonghui/.kube/config")
	K8sServer, err := NewK8sServer(config)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	err = K8sServer.CreateWithFile("/Users/shenkonghui/src/github/operator-test-framework/example/example1/deploy.yaml")
	if err != nil {
		fmt.Errorf(err.Error())
	}
}
