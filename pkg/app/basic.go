package app

import "errors"

const BasicType = "basic"

type Basic struct {
}

func (_ Basic) AppType() string {
	return BasicType
}

func (b *Basic) init(c *Config, path string) error {
	if err := createInitAppStructureBase(path); err != nil {
		return err
	}

	return nil
}

func (b *Basic) initConfig(c *Config) {
	c.Spec = &Spec{Type: BasicType}
}

func (b *Basic) generate(c *Config, path string) error {
	return errors.New("you cannot run generate on a basic app")
}
