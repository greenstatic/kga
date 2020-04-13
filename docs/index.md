# Overview of kga
<img src="assets/kga_logo.svg" style="width: 565px; max-width: 100%"/>

kga is a **CLI tool** to help you **manage your GitOps Kubernetes applications**.
It leverages industry standard tools like [Kustomize](https://kustomize.io) and [Helm](https://helm.sh) to enable easier GitOps practices when creating resources for your Kubernetes applications.   

!!! Warning
    kga is currently being developed, although we make sure our GitHub master branch always builds, you might experience some bugs or unfinished features.

## Okay, what REALLY is kga?
It is a simple CLI tool that replaces many bash scripts that one develops when journeying into Kubernetes GitOps.
The main aim of kga is to be the glue that holds established Kubernetes and GitOps tools and techniques together.
It does this by helping you manage the entire life cycle of an application:

1. Initial creation of GitOps application structure
2. Downloading of manifests
3. Overriding and [kustomizing](https://kustomize.io) the application's manifests
4. Static checks for current industry practices
5. Update the applications manifests when a new version is available 

!!! Note
    We have not implemented updated and check commands yet.

kga is only needed when creating and updating applications.
Once kga does its magic, you are just left with plain old YAML files, that you already know how to work with.

We try to reduce technical debt by sticking to Kubernetes established tooling such as Kustomization and Helm.
We also try to minimize the 

## Features
* 3 different application types:
    * [Basic](user-guide/basic-app.md) - you define all your manifests and update them by yourself
    * [Manifest](user-guide/manifest-app.md) - provide a URL and version where we can fetch manifests
    * [Helm](user-guide/helm-app.md) - provide a chart and helm override values and we will build your manifests
* Namespace overrides for manifest apps
* Exclusion of user defined resources
* Manifest URL templating 

## Why kga?
We are firm believers in maintaining applications in Kubernetes by practicing GitOps.
One of the challenges we faced when doing Kubernetes GitOps was the lack of tooling for maintain apps.
This and the number of shell scripts that were difficult to maintain as projects grew bigger made us develop a tool that would ease the process.
