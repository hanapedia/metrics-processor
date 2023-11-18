package domain

// Job represents similified Kubernetes Job
type Job struct {
	Name       string
	Namespace  string
	Deployment Deployment
	// any other fields you need...
}
