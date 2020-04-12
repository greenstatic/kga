package app

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestConfigVerify(t *testing.T) {
	assert_ := assert.New(t)

	type test struct {
		config string
		err    error
	}

	tests := []test{
		{`kind: kga-app
version: v1alpha
name: foo
spec:
  type: basic
`, nil},
		{`kind: kga-app
version: v1alpha
name: foo
spec:
  type: manifest
  manifest:
    version: v1.0.0
    urls:
    - https://example.com/{{.Version}}/manifest.yaml
`, nil},
		{`kind: kga-app
version: v1alpha
name: foo
spec:
  type: helm
  namespace: foo
  helm:
    chartName: nginx-ingress
    version: 1.34.3
    repoName: stable
    repoUrl: https://kubernetes-charts.storage.googleapis.com
    valuesFile: helm_values.yaml
`, nil},
		{`version: v1alpha
name: foo
spec:
  type: basic
`, ConfigBadFieldValueError("")},
		{`kind: kgaapp
version: v1alpha
name: foo
spec:
  type: basic
`, ConfigBadFieldValueError("")},
		{`kind: kga-app
version: v1
name: foo
spec:
  type: basic
`, ConfigBadFieldValueError("")},
		{`kind: kga-app
version: v1alpha
name: foo
`, ConfigMissingFieldError("")},
		{`kind: kga-app
version: v1alpha
spec:
  type: basic
`, ConfigBadFieldValueError("")},
		{`kind: kga-app
version: v1alpha
name: foo
spec:
  type: manifest
  manifest:
    urls:
    - https://example.com/latest/manifest.yaml
`, ConfigMissingFieldError("")},
		{`kind: kga-app
version: v1alpha
name: foo
spec:
  type: manifest
  manifest:
    version: v1.0.0
`, ConfigMissingFieldError("")},
		{`kind: kga-app
version: v1alpha
name: foo
spec:
  type: helm
  helm:
    chartName: nginx-ingress
    version: 1.34.3
    repoName: stable
    repoUrl: https://kubernetes-charts.storage.googleapis.com
    valuesFile: helm_values.yaml
`, ConfigMissingFieldError("")},
		{`kind: kga-app
version: v1alpha
name: foo
spec:
  type: helm
  namespace: foo
  helm:
    chartName: nginx-ingress
    version: 1.34.3
    repoName: stable
    repoUrl: https://kubernetes-charts.storage.googleapis.com
`, nil},
		{`kind: kga-app
version: v1alpha
name: foo
spec:
  type: helm
  namespace: foo
  helm:
    version: 1.34.3
    repoName: stable
    repoUrl: https://kubernetes-charts.storage.googleapis.com
    valuesFile: helm_values.yaml
`, ConfigMissingFieldError("")},
		{`kind: kga-app
version: v1alpha
name: foo
spec:
  type: helm
  namespace: foo
  helm:
    chartName: nginx-ingress
    repoName: stable
    repoUrl: https://kubernetes-charts.storage.googleapis.com
    valuesFile: helm_values.yaml
`, ConfigMissingFieldError("")},
		{`kind: kga-app
version: v1alpha
name: foo
spec:
  type: helm
  namespace: foo
  helm:
    chartName: nginx-ingress
    version: 1.34.3
    repoUrl: https://kubernetes-charts.storage.googleapis.com
    valuesFile: helm_values.yaml
`, ConfigMissingFieldError("")},
		{`kind: kga-app
version: v1alpha
name: foo
spec:
  type: helm
  namespace: foo
  helm:
    chartName: nginx-ingress
    version: 1.34.3
    repoName: stable
    valuesFile: helm_values.yaml
`, ConfigMissingFieldError("")},
		{`kind: kga-app
version: v1alpha
name: foo
spec:
  type: helm
  namespace: foo
  helm:
    chartName: nginx-ingress
    version: 1.34.3
    repoName: stable
    repoUrl: https://kubernetes-charts.storage.googleapis.com
    valuesFile: helm_values.yaml
`, nil},
		{`kind: kga-app
version: v1alpha
name: m
spec:
  namespace: m
  type: manifest
  manifest:
    version: stable
    urls:
      - https://raw.githubusercontent.com/argoproj/argo-cd/{{.Version}}/manifests/install.yaml
  exclude:
  - apiVersion: "v1"
`, nil},
	}

	for i, tst := range tests {
		c := Config{}
		b := []byte(tst.config)
		err := yaml.Unmarshal(b, &c)
		assert_.NoError(err, fmt.Sprintf("Failed test: %d, returned: %s", i+1, err))

		err = c.Verify()
		assert_.IsType(tst.err, err, fmt.Sprintf("Failed test: %d, returned: %s", i+1, err))
	}
}
