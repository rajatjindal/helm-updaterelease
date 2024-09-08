package main

import (
	"fmt"
	"os"

	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

var configFlags = genericclioptions.NewConfigFlags(true)

// following are just example placeholders, replace it with your values
// you can get it by running "kubectl get secret -n <namespace>"
const (
	helmReleaseSecretName = "sh.helm.release.v1.spin-operator.v1"
	helmReleaseNamespace  = "spin-operator"
)

func main() {
	kubeclient, err := getKubernetesClientset()
	if err != nil {
		fmt.Println("ERROR: ", err)
		os.Exit(1)
	}

	storage := driver.NewSecrets(kubeclient.CoreV1().Secrets(helmReleaseNamespace))
	releaseMd, err := storage.Get(helmReleaseSecretName)
	if err != nil {
		fmt.Println("ERROR: ", err)
		os.Exit(1)
	}

	fmt.Printf("Release -> %#v\n", releaseMd.Namespace)

	//make changes that you want to make to the release metadata and then
	releaseMd.SetStatus(release.StatusDeployed, "deployed successfully")

	err = storage.Update(helmReleaseSecretName, releaseMd)
	if err != nil {
		fmt.Println("ERROR: ", err)
		os.Exit(1)
	}
}

func getKubernetesClientset() (kubernetes.Interface, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
