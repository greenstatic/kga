package app

import (
	"github.com/greenstatic/kga/pkg/log"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
)

func kustomizeExecutablePath() string {
	env := os.Getenv("KUSTOMIZE")
	if env != "" {
		return env
	}
	return "kustomize" // assume that is is in the users path
}

func kustomizeBuildSuccseeds(path string) error {
	cmd := kustomizeExecutablePath()
	log.Infof("Checking if kustomize can build: %s build %s", cmd, path)
	k := exec.Command(cmd, "build", path)
	return k.Run()
}

func kustomizeAddKeyValue(kustomization, key, value string) string {
	if kustomization == "" {
		node := make(map[string]interface{})
		node[key] = value
		kustomizationUpdated, err := yaml.Marshal(&node)
		if err != nil {
			panic(err)
		}
		return string(kustomizationUpdated)
	}

	doc := yaml.Node{}

	if err := yaml.Unmarshal([]byte(kustomization), &doc); err != nil {
		panic(err)
	}

	if !(len(doc.Content) > 0 && doc.Content[0].Kind == yaml.MappingNode) {
		panic(errors.New("kustomization yaml is not a mapping node (root node)"))
	}

	keyExists := false

	for i := 0; i < len(doc.Content[0].Content); i += 2 {
		nodeKey := doc.Content[0].Content[i]
		nodeValue := doc.Content[0].Content[i+1]

		if nodeKey.Kind == yaml.ScalarNode && nodeKey.Value == key {
			nodeValue.Value = value
			keyExists = true
		}
	}

	if !keyExists {
		nodeKey := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: key,
			Tag:   "!!str",
		}
		nodeValue := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: value,
			Tag:   "!!str",
		}

		doc.Content[0].Content = append(doc.Content[0].Content, nodeKey)
		doc.Content[0].Content = append(doc.Content[0].Content, nodeValue)
	}

	kustomizationUpdated, err := yaml.Marshal(&doc)
	if err != nil {
		panic(err)
	}

	return string(kustomizationUpdated)
}

func ComperatorEqualStrings(existing, value string) bool {
	return existing == value
}

func ComperatorEqualStringPathsWrapper(basePath string) func(string, string) bool {
	return func(existing, value string) bool {
		a := existing
		if !filepath.IsAbs(existing) {
			a = filepath.Join(basePath, existing)
		}

		b := value
		if !filepath.IsAbs(value) {
			b = filepath.Join(basePath, value)
		}

		return a == b
	}
}


func kustomizeAddListElement(kustomization, key, value string, comperator func(existing, value string) bool) string {
	if kustomization == "" {
		node := make(map[string][]string)
		node[key] = []string{value}
		kustomizationUpdated, err := yaml.Marshal(&node)
		if err != nil {
			panic(err)
		}
		return string(kustomizationUpdated)
	}

	doc := yaml.Node{}

	if err := yaml.Unmarshal([]byte(kustomization), &doc); err != nil {
		panic(err)
	}

	if !(len(doc.Content) > 0 && doc.Content[0].Kind == yaml.MappingNode) {
		panic(errors.New("kustomization yaml is not a mapping node (root node)"))
	}

	keyExists := false

	for i := 0; i < len(doc.Content[0].Content); i += 2 {
		nodeKey := doc.Content[0].Content[i]
		nodeValue := doc.Content[0].Content[i+1]

		if nodeKey.Kind == yaml.ScalarNode && nodeKey.Value == key && nodeValue.Kind == yaml.SequenceNode {
			duplicate := false
			for _, item := range nodeValue.Content {
				if comperator(item.Value, value) {
					duplicate = true
				}
			}

			if !duplicate {
				nodeValue.Content = append(nodeValue.Content, &yaml.Node{
					Kind:        yaml.ScalarNode,
					Tag:         "!!str",
					Value:       value,
				})
			}
			keyExists = true
			break
		}
	}

	if !keyExists {
		nodeKey := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: key,
			Tag:   "!!str",
		}
		nodeContentValue := &yaml.Node{
			Kind: yaml.ScalarNode,
			Value: value,
			Tag: "!!str",
		}
		nodeValue := &yaml.Node{
			Kind:  yaml.SequenceNode,
			Content: make([]*yaml.Node, 0),
			Tag:   "!!seq",
		}

		nodeValue.Content = append(nodeValue.Content, nodeContentValue)

		doc.Content[0].Content = append(doc.Content[0].Content, nodeKey)
		doc.Content[0].Content = append(doc.Content[0].Content, nodeValue)
	}

	kustomizationUpdated, err := yaml.Marshal(&doc)
	if err != nil {
		panic(err)
	}

	return string(kustomizationUpdated)
}
