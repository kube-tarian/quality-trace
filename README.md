<p align="center"><b>Quality Trace is a framework for end-to-end testing based on OpenTelemetry tracing.</b></p>

<h4 align="center">
    <a href="https://github.com/kube-tarian/quality-trace/discussions">Discussions</a> 
</h4>

<h4 align="center">

[![Docker Image CI](https://github.com/kube-tarian/quality-trace/actions/workflows/docker-image.yaml/badge.svg)](https://github.com/kube-tarian/quality-trace/actions/workflows/docker-image.yaml)
[![CodeQL](https://github.com/kube-tarian/quality-trace/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/kube-tarian/quality-trace/actions/workflows/codeql-analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kube-tarian/quality-trace)](https://goreportcard.com/report/github.com/kube-tarian/quality-trace)

[![Price](https://img.shields.io/badge/price-FREE-0098f7.svg)](https://github.com/kube-tarian/quality-trace/blob/main/LICENSE)
[![Discussions](https://badgen.net/badge/icon/discussions?label=open)](https://github.com/kube-tarian/quality-trace/discussions)
[![Code of Conduct](https://badgen.net/badge/icon/code-of-conduct?label=open)](./code-of-conduct.md)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

</h4>

<hr>

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

