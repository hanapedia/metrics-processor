# Metrics Processor

This project contains a go program that queries prometheus metrics and saves the result matrix as json in S3 object storage. It also includes example deployemnt to Kubernetes as Job using Kustomize.

## Building the program
```sh
make prod TAG=[insert tag here]
```

## Running the program
### Prerequisite
- Kubernetes cluster with following installed
  - Prometheus
  - Linkerd (if you want linkerd metrics)
  - Sealed secrets
- AWS account with S3 Bucket
```sh
# prepare aws secrets
bash prepare_aws_secret.sh

# copy the example env and add values
cp _k8s/overlays/example/metrics-processor.env.example _k8s/overlays/example/metrics-processor.env
vi _k8s/overlays/example/metrics-processor.env

# apply with kustomize
kubectl apply -k _k8s/overlays/example
```
