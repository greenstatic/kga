package generate

import (
	"errors"
	"github.com/ghodss/yaml"
	"github.com/greenstatic/kga/internal/layout"
	"github.com/greenstatic/kga/pkg/config"
	"github.com/greenstatic/kga/pkg/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func RemoveExcludedBaseManifests(appPath string, excludedManifests *[]config.ExcludeItemSpec) error {
	pathBase := filepath.Join(appPath, "base")
	pathManifests := filepath.Join(pathBase, "manifests")

	files, err := ioutil.ReadDir(pathManifests)
	if err != nil {
		return err
	}

	keptResourcesCount := 0
	excludedResourcesCount := 0

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		content, err := ioutil.ReadFile(filepath.Join(pathManifests, f.Name()))
		if err != nil {
			return err
		}

		contentStr := string(content)
		yamls := seperateYAMLDocuments(contentStr)

		fileManifestDocuments := make([]YAMLDocument, 0)
		excludedDocuments := make([]YAMLDocument, 0)

		for _, exclude := range *excludedManifests {
			for _, yamlDoc := range yamls {
				ymlDoc := new(interface{})
				if err := yaml.Unmarshal([]byte(string(yamlDoc)), ymlDoc); err != nil {
					return err
				}

				excludeDoc := reflect.ValueOf(exclude).Interface()
				match, err := excludeItemMatchResource(&excludeDoc, ymlDoc)
				if err != nil {
					return err
				}

				if match {
					//excludedDocuments, _ := excludedDocuments[f.Name()]
					excludedDocuments = append(excludedDocuments, yamlDoc)

				} else {
					fileManifestDocuments = append(fileManifestDocuments, yamlDoc)
				}
			}
		}

		// Override the manifest file with the non excluded YAML documents
		manifestContent := joinYAMLDocuments(fileManifestDocuments)
		if err := ioutil.WriteFile(filepath.Join(pathManifests, f.Name()), []byte(manifestContent), os.FileMode(0640)); err != nil {
			return err
		}
		keptResourcesCount += len(fileManifestDocuments)

		// Add excluded documents to <app>/base/excluded/
		if len(excludedDocuments) > 0 {
			exists, err := layout.FileOrDirExists(filepath.Join(pathBase, "excluded"))
			if err != nil {
				return err
			}

			if !exists {
				// Create excluded dir
				if err := os.Mkdir(filepath.Join(pathBase, "excluded"), os.FileMode(0755)); err != nil {
					return err
				}
			}

			// Add file with excluded content
			excludedContent := joinYAMLDocuments(excludedDocuments)
			if err := ioutil.WriteFile(filepath.Join(pathBase, "excluded", f.Name()), []byte(excludedContent), os.FileMode(0640)); err != nil {
				return err
			}
		}
		excludedResourcesCount += len(excludedDocuments)
	}

	log.Infof("Kept: %d resources, excluded: %d resources", keptResourcesCount, excludedResourcesCount)

	return nil
}

// Pass pointers to the variables that hold the YAML definitions. Returns true if the excludeItem value matches
// the resource value.
func excludeItemMatchResource(excludeItem *interface{}, resource *interface{}) (bool, error) {
	if excludeItem == nil || resource == nil {
		return false, errors.New("comparing nil value")
	}

	switch reflect.TypeOf(*excludeItem).Kind() {
	case reflect.Map:
		if reflect.TypeOf(*resource).Kind() != reflect.Map {
			return false, nil
		}

		for _, v := range reflect.ValueOf(*excludeItem).MapKeys() {
			match := false
			for _, v2 := range reflect.ValueOf(*resource).MapKeys() {
				if reflect.ValueOf(v.Interface()).String() == reflect.ValueOf(v2.Interface()).String() {
					// Key matches, check if value matches as well
					left := reflect.ValueOf(*excludeItem).MapIndex(v).Interface()
					right_ := reflect.ValueOf(*resource).MapIndex(v2)

					if !right_.IsValid() || right_.IsZero() || right_.IsNil() {
						return false, nil
					}
					right := right_.Interface()

					m, err := excludeItemMatchResource(&left, &right)
					if err != nil {
						return false, err
					}

					if m {
						match = true
						break
					}
				}
			}

			if !match {
				return false, nil
			}
		}

		return true, nil

	case reflect.Slice:
		if reflect.TypeOf(*resource).Kind() != reflect.Slice {
			return false, nil
		}

		for i := 0; i < reflect.ValueOf(*excludeItem).Len(); i++ {
			left := reflect.ValueOf(*excludeItem).Index(i).Interface()
			match := false
			for j := 0; j < reflect.ValueOf(*resource).Len(); j++ {
				right := reflect.ValueOf(*excludeItem).Index(j).Interface()

				m, err := excludeItemMatchResource(&left, &right)
				if err != nil {
					return false, err
				}

				if m {
					match = true
					break
				}
			}

			if !match {
				return false, nil
			}
		}

		return true, nil

	case reflect.String:
		leftStr := reflect.ValueOf(*excludeItem).String()
		if reflect.TypeOf(*resource).Kind() != reflect.String {
			return false, nil
		}

		rightStr := reflect.ValueOf(*resource).String()
		return leftStr == rightStr, nil

	case reflect.Int:
		leftInt := reflect.ValueOf(*excludeItem).Int()
		if reflect.TypeOf(*resource).Kind() != reflect.Int {
			return false, nil
		}

		rightInt := reflect.ValueOf(*resource).Int()
		return leftInt == rightInt, nil

	default:
		return false, errors.New("unsupported kind: " + reflect.TypeOf(*excludeItem).Kind().String())
	}
}

type YAMLDocument string

func seperateYAMLDocuments(fileContent string) []YAMLDocument {
	docRaw := strings.Split(fileContent, "---")

	documents := make([]YAMLDocument, 0, len(docRaw))
	for _, docContent := range docRaw {

		docContentWithoutWhitespace := strings.ReplaceAll(docContent, " ", "")
		docContentWithoutWhitespace = strings.ReplaceAll(docContentWithoutWhitespace, "\n", "")
		docContentWithoutWhitespace = strings.ReplaceAll(docContentWithoutWhitespace, "\r", "")
		docContentWithoutWhitespace = strings.ReplaceAll(docContentWithoutWhitespace, "\t", "")

		if len(docContentWithoutWhitespace) > 0 {
			documents = append(documents, YAMLDocument(strings.Trim(docContent, "\n ")))
		}
	}

	if len(documents) == 0 {
		documents = append(documents, YAMLDocument(docRaw[0]))
	}

	return documents
}

func joinYAMLDocuments(documents []YAMLDocument) string {
	xs := make([]string, 0, len(documents))

	for _, v := range documents {
		xs = append(xs, string(v))
	}

	joined := strings.Join(xs, "\n---\n")

	addTrailingNewline := true
	if len(documents) == 1 {
		docContent := string(documents[0])
		docContentWithoutWhitespace := strings.ReplaceAll(docContent, " ", "")
		docContentWithoutWhitespace = strings.ReplaceAll(docContentWithoutWhitespace, "\n", "")
		docContentWithoutWhitespace = strings.ReplaceAll(docContentWithoutWhitespace, "\r", "")
		docContentWithoutWhitespace = strings.ReplaceAll(docContentWithoutWhitespace, "\t", "")

		if len(docContentWithoutWhitespace) == 0 {
			addTrailingNewline = false
		}
	}

	if addTrailingNewline {
		joined += "\n"
	}

	return joined
}
