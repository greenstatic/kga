package yamlmanipulation

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDocumentsSeperate(t *testing.T) {
	assert := assert.New(t)

	content0 := `

`
	content1 := `foo: bar
spec:
  one:
  - 1
  - 2
  - 3
  - 4
`
	content2 := `---
foo: bar
`
	content3 := `---
foo: bar
---
wow: foo
`
	content4 := `---
foo: bar
---
wow: foo
---
`
	content5 := `---
foo: bar
---
wow: foo
---
two: one
`

	type test struct {
		content                       string
		expectedNumberOfYAMLDocuments int
	}

	tests := []test{
		{
			content0,
			1,
		},
		{
			content1,
			1,
		},
		{
			content2,
			1,
		},
		{
			content3,
			2,
		},
		{
			content4,
			2,
		},
		{
			content5,
			3,
		},
	}

	for i, tst := range tests {
		documents := DocumentsSeperate(tst.content)
		assert.Equal(tst.expectedNumberOfYAMLDocuments, len(documents), fmt.Sprintf("Failed test: %d", i + 1))
	}
}

func TestDocumentsJoin(t *testing.T) {
	assert := assert.New(t)

	content0 := `

`
	content1 := `foo: bar
spec:
 one:
 - 1  # WOW
 - 2
 - 3
 - 4
`
	content2 := `---
foo: bar
`
	content2_ := `foo: bar
`
	content3 := `---
# TODO
foo: bar
---
wow: foo
`
	content3_ := `# TODO
foo: bar
---
wow: foo
`
	content4 := `---
foo: bar
---
wow: foo
---
`
	content4_ := `foo: bar
---
wow: foo
`
	content5 := `---
foo: bar
---
wow: foo
---
two: one
`
	content5_ := `foo: bar
---
wow: foo
---
two: one
`

	type test struct {
		content  string
		expected string
	}

	tests := []test{
		{
			content0,
			content0,
		},
		{
			content1,
			content1,
		},
		{
			content2,
			content2_,
		},
		{
			content3,
			content3_,
		},
		{
			content4,
			content4_,
		},
		{
			content5,
			content5_,
		},
	}

	for i, tst := range tests {
		documents := DocumentsSeperate(tst.content)
		joined := DocumentsJoin(documents)
		assert.Equal(tst.expected, joined, fmt.Sprintf("Failed test: %d", i))
	}
}
