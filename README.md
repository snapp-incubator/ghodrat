# GHODRAT

> WebRTC media servers stress testing tool (currently only Janus)

## Deployment

```zsh
# update or create manifests
kubectl apply -f ./deployments/k8s/janus/configmap.yml
kubectl apply -f ./deployments/k8s/janus/job.yml

# delete manifests
kubectl delete -f ./deployments/k8s/janus/configmap.yml
kubectl delete -f ./deployments/k8s/janus/job.yml
```

### troubleshooting image

- docker container run --entrypoint /bin/sh -it --rm ghcr.io/snapp-incubator/ghodrat-janus:v1.1.0
