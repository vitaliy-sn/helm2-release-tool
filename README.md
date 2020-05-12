# helm-release-tool

## Build:

```shell
go mod download
go build .
```

## Usage:

```shell
kubectl get cm prometheus.v25 -o json | helm-release-tool info
kubectl get cm prometheus.v25 -o json | helm-release-tool set-status-deployed
kubectl get cm prometheus.v25 -o json | helm-release-tool set-new-release-name <string>
```
