# helm-release-tool

## Build:

```shell
./build.sh
```

## Usage examples:

```shell
kubectl get cm example.v1 -o json | helm-release-tool info | kubectl apply -f -
kubectl get cm example.v1 -o json | helm-release-tool get-manifests | kubectl apply -f -
kubectl get cm example.v1 -o json | helm-release-tool set-status-deployed | kubectl apply -f -
kubectl get cm example.v1 -o json | helm-release-tool set-release-name <string> | kubectl apply -f -
kubectl get cm example.v1 -o json | helm-release-tool set-release-namespace <string> | kubectl apply -f -
```
