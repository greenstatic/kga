# Basic App
Basic apps are apps which do not need to pull from the internet manifests and require manual creation of manifests.
This is useful for your own custom applications, where you do not have a manifest file accessible from a public URL or a helm chart.

In principal they are quite similar to manifest typed apps.

The following is a valid _kga.yaml_ configuration:
```yaml
kind: kga-app
version: v1alpha
name: foo
spec:
  type: basic
  exclude:
  - kind: Secret
```

Firt run the `kga init basic <app\>` command, then the user is expected to the add manifests according to the following rules:

* _base/manifests_: Generic manifests for the app (e.g. no ingress rules). We should be able to copy/paste this into another kga basic app without any issues.
* _base/kustomization.yaml_: List all the manifests from _base/manifests_ in the resources list.
* _overlay/resources_: Application instance resources, e.g. Namespace, Ingress, ConfigMap etc.
* _overlay/patches_: Patches to fix the _base/manifests_ manifests.
* _overlay/kustomization.yaml_: List all the manifests from _overlay/resources_, list all _overlay/patches_ and import _base_
* _kustomization.yaml_: Import _overlay_.

Once you finish editing your manifests run:
```bash
kustomize build <app>
```
To verify if the manifests have been correctly overridden, patched and additional resources added.
