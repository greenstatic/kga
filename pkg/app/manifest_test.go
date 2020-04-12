package app

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManifestUrlApplyTemplate(t *testing.T) {
	assert_ := assert.New(t)

	type test struct {
		url string

		resultUrl string
		isError   bool
	}

	c := Config{
		Kind:    ConfigKind,
		Version: ConfigVersion,
		Name:    "foo",
		Spec:    &Spec{
			Namespace: "foo",
			Type:      ManifestType,
			Manifest:  &ManifestSpec{
				Version:  "v1.2.3",
				Urls: []string{"http://example.com/{{ .Version }}/manifest.yaml"},
				Template: map[string]string{"bar": "foo", "one": "1"},
			},
		},
	}

	tests := []test{
		{
			url:       "http://example.com/{{ .Version }}/manifest.yaml",
			resultUrl: "http://example.com/v1.2.3/manifest.yaml",
			isError:   false,
		},
		{
			url:       "http://example.com/{{ .Version }}/manifest.yaml",
			resultUrl: "http://example.com/v1.2.3/manifest.yaml",
			isError:   false,
		},
		{
			url:       "http://example.com/{{ .Version }}/{{ .Config.Name }}.yaml",
			resultUrl: "http://example.com/v1.2.3/foo.yaml",
			isError:   false,
		},
		{
			url:       "http://example.com/{{ .version }}/manifest.yaml",
			resultUrl: "",
			isError:   true,
		},
		{
			url:       "http://example.com/{{ .Versionz }}/manifest.yaml",
			resultUrl: "",
			isError:   true,
		},
	}

	for i, tst := range tests {
		url, err := manifestUrlApplyTemplate(&c, tst.url)

		assert_.Equal(tst.resultUrl, url, fmt.Sprintf("Failed test: %d", i + 1))
		if tst.isError {
			assert_.Error(err, fmt.Sprintf("Failed test: %d", i + 1))
		} else {
			assert_.NoError(err, fmt.Sprintf("Failed test: %d", i + 1))
		}
	}
}
