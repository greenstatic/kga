# kga - Kubernetes GitOps App CLI tool
A CLI tool to help you manage your GitOps Kubernetes applications.

This project is under heavy development, all configurations, flags and commands are probably going to break in future versions.


The main aim of this tool is to be the glue that holds established Kubernetes and GitOps tools and techniques to manage and maintain applications in a GitOps fashion.
kga tries to reduce the repetitive tasks a DevOps engineer must perform to create/update an application.

kga is only needed when creating and updating applications.
Once an application is created you can edit the resource YAMLs as you previously have, by using kustomization.

We try to reduce technical debt by sticking to Kubernetes established tooling such as: Helm and Kustomization.


## kga Application Structure
The following is the kga folder structure for a Helm based application that kga generates.
```
app/
    base/
        manifests/              <- All chart resources
        kustomization.yaml
    
    overlay/                    <- Our resources
        pathces/                <- Patches the base manifests
        resources/              <- Additional resources e.g. namespace resource
        kustomization.yaml
    
    helm_values.yaml            <- Used when generating the chart template
    kga.yaml                    <- kga configuration file
    kustomization.yaml          <- Entrypoint into the app (kubectl apply -k .)
```

## Installation
```bash
go install github.com/greenstatic/kga/cmd/kga
```

### Installation from Source
1. Clone repo
2. Run `make`
3. The built executable should be in `./bin/kga`

## Usage
1. Move into the directory where you wish to save your Kubernetes applications
2. Run: `kga create <app-name>`
3. If this is a Helm based app, edit `app-name/helm_values.yaml`
4. Edit `app-name/kga.yaml` with all the #TODO fields filled out
5. Run: `kga generate <app-name>`
6. Now you can simply do `kubectl apply -k <app-name>` or `kustomize build <app-name>` to view your app's resources

## Example kga.yaml file
### App Type Helm
```yaml
kind: kga-app
version: v1alpha
name: nginx-ingress
spec:
  namespace: nginx-ingress
  helm:
    chartName: nginx-ingress
    version: 1.34.3
    repoName: stable
    repoUrl: https://kubernetes-charts.storage.googleapis.com/
    valuesFile: ./helm_values.yaml

  # Used just to demonstrate the usage of exclude spec
  exclude:
  - apiVersion: v1
    kind: ServiceAccount
    metadata:
      labels:
        app: nginx-ingress
```

### App Type Manifest
```yaml
kind: kga-app
version: v1alpha
name: kubernetes-dashboard
spec:
  namespace: kubernetes-dashboard
  manifest:
    urls:
      - "https://raw.githubusercontent.com/kubernetes/dashboard/{{ .version }}/{{ .foo }}/deploy/recommended.yaml"
    template:
      version: v2.0.0-rc7
      foo: aio

  # Used just to demonstrate the usage of exclude spec
  exclude:
  - kind: Secret
```

## Why Did We Develop kga?
We are firm believers in maintaining applications in Kubernetes by practicing GitOps.
One of the challenges we faced when doing Kubernetes GitOps was the lack of tooling for maintain apps.
This and the number of shell scripts that were difficult to maintain as projects grew bigger made us develop a tool that would ease the process.