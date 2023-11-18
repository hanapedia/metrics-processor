package main

import (
	"k8s.io/client-go/dynamic"
	k8sClientGo "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/hanapedia/rca-experiment-runner/infrastructure/chaosmesh"
	"github.com/hanapedia/rca-experiment-runner/infrastructure/env"
	k8sInfra "github.com/hanapedia/rca-experiment-runner/infrastructure/kubernetes"
	"github.com/hanapedia/rca-experiment-runner/pkg/application/service"
)

func main() {
	// Prepare experiment configs
	config, err := env.NewExperimentConfig()
	if err != nil {
		panic(err.Error())
	}

	// Load kubeconfig from KUBECONFIG
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loadingRules,
		configOverrides,
	)
	kubeConfig, err := clientConfig.ClientConfig()
	if err != nil {
		panic(err.Error())
	}

	// prepare kube client for kubernetes API
	clientset, err := k8sClientGo.NewForConfig(kubeConfig)
	if err != nil {
		panic(err.Error())
	}
	kubernetesAdapter := k8sInfra.NewKubernetesAdapter(clientset, config)

	// Prepare kube dynamic config for chaos mesh resource
	dynamicClient, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		panic(err.Error())
	}
	chaosAdapter := chaosmesh.NewChaosMeshAdapter(dynamicClient, config)

	experimentRunner := service.NewExperimentRunner(config, kubernetesAdapter, chaosAdapter)
	err = experimentRunner.RunExperiment()
	if err != nil {
		panic(err.Error())
	}
}
