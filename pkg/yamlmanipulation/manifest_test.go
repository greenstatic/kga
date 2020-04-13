package yamlmanipulation

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestExcludeItemMatchResource(t *testing.T) {
	assert := assert.New(t)

	_excludeItemYAML1 := `apiVersion: v1
kind: Secret
metadata:
  labels:
    foo: bar  # This is a comment
`

	_excludeItemResource1 := `apiVersion: v1
kind: Secret
metadata:
  name: foo-bar  # Another comment
  labels:
    abba: babba
    foo: bar
type: Opaque
data:
  username: admin
  password: admin
`

	_excludeItemResource2 := `apiVersion: v1
kind: Secret
metadata:
  name: not-foo
type: Opaque
data:
  username: admin
  password: admin
`

	_excludeItemResource3 := `apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  labels:
    app: foo
    chart: foo-1.14.3
    heritage: Helm
    release: foo
  name: foo
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: foo
subjects:
  - kind: ServiceAccount
    name: foo
    namespace: foo
`

	_excludeItemYAML2 := `apiVersion: v1
kind: Random
spec:
  somekey:
  - one
  - two
  - three
`

	_excludeItemResource4 := `apiVersion: v1
kind: Random
metadata:
  labels:
    foo: bar
  name: some-random-resource
spec:
  somekey:
  - four
  - three
  - one
  - two
  - four
`

	_excludeItemResource5 := `apiVersion: v1
kind: Random
metadata:
  labels:
    foo: bar
  name: some-random-resource
spec:
  somekey: somevalue
`

	_excludeItemResource6 := `apiVersion: v1
kind: Random
metadata:
  labels:
    foo: bar
  name: some-random-resource
spec:
  somekey: []
`

	_excludeItemYAML3 := `apiVersion: v1
kind: Random
spec:
  somekey:
  - 1
  - 2
  - 3
`
	_excludeItemResource7 := `apiVersion: v1
kind: Random
metadata:
  labels:
    foo: bar
  name: some-random-resource
spec:
  somekey:
  - 4
  - 3
  - 1
  - 2
  - 4
`

	type test struct {
		exclude      string
		resource     string
		expectResult bool
	}

	tests := []test{
		test{
			_excludeItemYAML1,
			_excludeItemResource1,
			true,
		},
		{
			_excludeItemYAML1,
			_excludeItemResource2,
			false,
		},
		{
			_excludeItemYAML1,
			_excludeItemResource3,
			false,
		},
		{
			_excludeItemYAML2,
			_excludeItemResource4,
			true,
		},
		{
			_excludeItemYAML2,
			_excludeItemResource5,
			false,
		},
		{
			_excludeItemYAML2,
			_excludeItemResource6,
			false,
		},
		{
			_excludeItemYAML3,
			_excludeItemResource7,
			true,
		},
		{
			_excludeItemYAML3,
			_excludeItemResource6,
			false,
		},
	}

	for i, tst := range tests {
		exclude := new(interface{})
		if err := yaml.Unmarshal([]byte(tst.exclude), &exclude); err != nil {
			panic(err)
		}

		resource := new(interface{})
		if err := yaml.Unmarshal([]byte(tst.resource), &resource); err != nil {
			panic(err)
		}

		ans, err := ExcludeItemMatchesResource(exclude, resource)
		assert.NoError(err, fmt.Sprintf("Failed test: %d", i+1))
		assert.Equal(tst.expectResult, ans, fmt.Sprintf("Failed test: %d", i+1))
	}

}
