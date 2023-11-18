package domain

// ChaosResource represents a Chaos Mesh Resource in your domain model
type ChaosResource struct {
	Name       string
	Namespace  string
	Deployment Deployment
	// any other fields you need...
}
