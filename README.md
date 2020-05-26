# helm-release-tool

## Build:

```shell
./build.sh
```

## Usage examples:

```shell
kubectl get cm example.v1 -o json | helm-release-tool info
kubectl get cm example.v1 -o json | helm-release-tool get-manifest
kubectl get cm example.v1 -o json | helm-release-tool set-manifest <path-to-new-manifest> | kubectl apply -f -
kubectl get cm example.v1 -o json | helm-release-tool set-status-deployed | kubectl apply -f -
kubectl get cm example.v1 -o json | helm-release-tool set-release-name <string> | kubectl apply -f -
kubectl get cm example.v1 -o json | helm-release-tool set-release-namespace <string> | kubectl apply -f -
```
