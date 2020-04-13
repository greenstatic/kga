# Getting Started
## Installation
See the [installation instructions](installation.md) on how to install kga.

We assume you have a working kga, Kustomize and Helm installation.

## Your First kga Manifest App
This short guide will demonstrate how to create a [Kubernetes Dashboard](https://github.com/kubernetes/dashboard/) kga manifest app.
What this means is that we will get the applications (Kubernetes Dashboard) manifests from a web URL.

### 1. Initialize our Kubernetes Dashboard app
```bash
kga init manifest dashboard
```
This will created a directory called `dashboard` which contains the application structure for a Kustomize GitOps application.
If you are curious about the structure, checkout out the [short explanation on the Core Concepts](core-concepts.md#general-folder-structure-of-kga-app) page.

### 2. Update the `dashboard/kga.yaml` file
```yaml
kind: kga-app
version: v1alpha
name: dashboard
spec:
  namespace: dashboard
  type: manifest
  manifest:
    version: v2.0.0-rc7
    urls:
    - https://raw.githubusercontent.com/kubernetes/dashboard/{{ .Version }}/aio/deploy/recommended.yaml
  exclude:
  - kind: Namespace
  - kind: Secret
```

We got the URL from [Kubernetes Dashboard's Getting Starting instructions](https://github.com/kubernetes/dashboard/#getting-started).

### 3. Generate the manifests
```bash
kga generate dashboard
```

The final output line should be: `Successfully generated kga app`.

Now checkout the result.
We have all the required manfiests, our own namespace that overrides all the namespaced resources from the base/manifests (base manifest has the namespace _kubernetes-dashboard_ while our namespace is named _dashboard_).

### 4. Build using Kustomize
```bash
kustomize build dashboard
```

Kustomize can now be used to build all the manifests into one long manifest that can be piped into kubectl (or use kubectl -f <app\>) to deploy manually.

Thats it!

For more details about manifest apps visit: [User Guide / Manifest](user-guide/manifest-app.md).


## Your First kga Helm App
Now lets show how to make a kga Helm app.
This time we will demonstrate by creating a [NGINX Ingress](https://hub.helm.sh/charts/stable/nginx-ingress) kga helm app.

### 1. Initialize our NGINX Ingress ap
```bash
kga init helm nginx-ingress
```

### 2. Update the `nginx-ingress.yaml` file
```yaml
kind: kga-app
version: v1alpha
name: nginx-ingress
spec:
  namespace: nginx-ingress
  type: helm
  helm:
    chartName: nginx-ingress
    version: 1.34.3
    repoName: stable
    repoUrl: https://kubernetes-charts.storage.googleapis.com
    valuesFile: helm_values.yaml
  exclude:
  - kind: Secret
```

### 3. Generate the manifest
```bash
kga generate nginx-ingress
```

The final output line should be: `Successfully generated kga app`.

### 4. Build using Kustomize
```bash
kustomize build nginx-ingress
```

Thats it, you are done.

For more details about helm apps visit: [User Guide / Helm](user-guide/helm-app.md).

## Other
It is also possible to create an app that does not rely on 3rd party providers/services.
Such an app is called _basic_. 
Using `kga init basic <app\>` we create the now familiar project structure and you yourself then create the necessary _base/manifests_ files and link them together in all the kustomization.yaml files.
`kga generate <app\>` will fail on a basic app, since we do not have anything to pull from the internet and process.

For more details about basic apps visit: [User Guide / Basic](user-guide/basic-app.md).
