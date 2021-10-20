# GHODRAT

> WebRTC media servers stress testing tool (currently only Janus)

## Deployment

```zsh
kubectl create -f ./deployments/k8s/janus/configmap.yml
kubectl create -f ./deployments/k8s/janus/deployment.yml
```

### troubleshooting image

- docker container run --entrypoint /bin/sh -it --rm ghcr.io/snapp-incubator/ghodrat-janus:develop
