package port

import "github.com/hanapedia/rca-experiment-runner/pkg/domain"

// KubernetesClientPort defines the interface for interactions with the Kubernetes API
type KubernetesClientPort interface {
	// GetDeploymentsWithAnnotation retrieves deployments in the given namespace with the specified annotation
	GetDeploymentsWithOutAnnotation(namespace string, annotationKey string, annotationValue string) ([]domain.Deployment, error)
	// CreateAndApplyJobResource creates a new job resource and applies it for the specified deployment
	CreateAndApplyJobResource(deployment domain.Deployment) error
}

// ChaosExperimentsPort defines the interface for interactions with chaos experiment tools
type ChaosExperimentsPort interface {
	// CreateAndApplyNetworkDelay creates a new chaos resource for the specified deployment and apply it.
	CreateAndApplyNetworkDelay(deployment domain.Deployment) error
}
