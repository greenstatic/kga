# Core Concepts
kga is simple to use. 

Here is the typical workflow:

1. Initialize an app `kga init <type> <app>`
2. Make the necessary changes to the generated `kga.yaml` file
3. Generate the app `kga generate <app>`
4. Make any new changes to your <app\>/overlay
5. Simply use `kustomize build <app>` to get your final manifest ready for deployment
6. If you wish to make any changes, go ahead and change your <app\>/overlay, then just run `kga generate <app>` again

## General Folder Structure of kga App
When we initialize an app, this is the generic folder structure that will be created.
There are a couple of minor differences between the three different application types: basic, manifest and helm.
For example helm apps have an additional `helm_values.yaml` file in the app root dir.

```text
base/
  excluded/             # Resources that were excluded using kga.yaml exclude 
                          from base/manifests
  manifests/            # Downloaded application manifests
  kustomization.yaml    # kustomization file that brings together all manifests 
                          from base/ 

overlay/                # User specified manifests and patches
  patches/              # Patches that fix base/manifests
  resources/            # User specified additional resources
    namespace.yaml      # Optional automatic namespace manifest (if namespace 
                          is defined in kga.yaml)
  kustomization.yaml    # kustomization file that brings together base/ and 
                          overlay/

kga.yaml                # config for kga
kustomization.yaml      # main kustomization that just links to overlay
                          so we can `kustomize build <app>`
```
