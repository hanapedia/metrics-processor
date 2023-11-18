#!/bin/bash

SECRET_NAME=aws-credentials
NAMESPACE=rca

# Initial check for the secret
if kubectl -n $NAMESPACE get secret $SECRET_NAME >/dev/null 2>&1; then
    echo "Secret $SECRET_NAME already exists in namespace $NAMESPACE."
    exit 0
fi

AWS_PROFILE=default

AWS_ACCESS_KEY_ID=$(aws configure get aws_access_key_id --profile $AWS_PROFILE)
AWS_SECRET_ACCESS_KEY=$(aws configure get aws_secret_access_key --profile $AWS_PROFILE)

kubectl create secret generic $SECRET_NAME -n $NAMESPACE --dry-run=client -o yaml \
  --from-literal=aws_access_key_id=$AWS_ACCESS_KEY_ID \
  --from-literal=aws_secret_access_key=$AWS_SECRET_ACCESS_KEY > ./_k8s/$SECRET_NAME.yaml

kubeseal -o yaml < ./_k8s/$SECRET_NAME.yaml > ./_k8s/sealed-$SECRET_NAME.yaml

kubectl apply -n $NAMESPACE -f ./_k8s/sealed-$SECRET_NAME.yaml

while true; do
    if kubectl -n $NAMESPACE get secret $SECRET_NAME; then
        echo "Secret $SECRET_NAME found."
        break
    else
        echo "Secret $SECRET_NAME not found. Waiting..."
        sleep 2
    fi
done
