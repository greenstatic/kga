package yamlmanipulation

import (
	"errors"
	"gopkg.in/yaml.v2"
	"reflect"
)

// Pass pointers to the variables that hold the YAML definitions. Returns true if the excludeItem value matches
// the resource value.
func ExcludeItemMatchesResource(excludeItem *interface{}, resource *interface{}) (bool, error) {
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

					m, err := ExcludeItemMatchesResource(&left, &right)
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

				m, err := ExcludeItemMatchesResource(&left, &right)
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


func ExcludeResourceFromManifest(excludeStr, manifest string) (excludedManifest, newManifest string, err error) {
	exclude := new(interface{})

	if err := yaml.Unmarshal([]byte(excludeStr), exclude); err != nil {
		return "", "", err
	}

	yamls := DocumentsSeperate(manifest)

	fileManifestDocuments := make([]Document, 0)
	excludedDocuments := make([]Document, 0)

	for _, yamlDoc := range yamls {
		ymlDoc := new(interface{})
		if err := yaml.Unmarshal([]byte(string(yamlDoc)), ymlDoc); err != nil {
			return "", "", err
		}

		excludeDoc := reflect.ValueOf(*exclude).Interface()
		match, err := ExcludeItemMatchesResource(&excludeDoc, ymlDoc)
		if err != nil {
			return "", "", err
		}

		if match {
			excludedDocuments = append(excludedDocuments, yamlDoc)

		} else {
			fileManifestDocuments = append(fileManifestDocuments, yamlDoc)
		}
	}

	excludedManifest = DocumentsJoin(excludedDocuments)
	newManifest = DocumentsJoin(fileManifestDocuments)
	err = nil
	return
}
