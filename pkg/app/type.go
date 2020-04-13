package app

import (
	"fmt"
)

type Type interface {
	AppType() string
	init(c *Config, path string) error
	initConfig(*Config)
	generate(c *Config, path string) error
}

func CreateType(typeStr string) Type {
	switch typeStr {
	case BasicType:
		return &Basic{}
	case ManifestType:
		return &Manifest{}
	case HelmType:
		return &Helm{}
	}

	panic("invalid app type string")
}

type InvalidTypeStringError string

func (s InvalidTypeStringError) Error() string {
	return fmt.Sprintf("app type '%s' is an invalid, can only be: basic, manifest or helm", string(s))
}

func ValidateTypeString(type_ string) (err error) {
	defer func() {
		if recover() != nil {
			err = InvalidTypeStringError(type_)
		}
	}()

	CreateType(type_)
	return nil
}
