#!/bin/bash

echo "Building Server Imgage...\n"
docker build -t mmarci96/server-test:latest ./apps/server/
docker push mmarci96/server-test:latest
kubectl rollout restart deployment -n game-test server-test

echo "Building Go proxy and static files...\n"
docker build -t mmarci96/proxy:latest .
docker push mmarci96/proxy:latest
kubectl rollout restart deployment -n game-test proxy 

kubectl apply -f ./k8s/

