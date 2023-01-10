# quality-trace
Quality Trace is a framework for end-to-end testing based on OpenTelemetry tracing.


###  How to install and run quality-trace:

#### Prerequisites
* A Kubernetes cluster 
* Helm binary

#### Prepare Namespace
```bash
kubectl create namespace quality-trace
```

#### Installation
```bash
helm repo add quality-trace https://kube-tarian.github.io/quality-trace
helm repo update

helm upgrade -i quality-trace quality-trace/quality-trace -n quality-trace
```

