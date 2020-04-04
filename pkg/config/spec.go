package config

type Config struct {
	Kind    string
	Version string
	Name    string
	Spec    *Spec `yaml:"spec,omitempty"`
}

type AppType string

const (
	AppTypeHelm     = AppType("helm")
	AppTypeManifest = AppType("manifest")
	configVersion   = "v1alpha"
)

type Spec struct {
	Namespace string `yaml:"namespace"`
	Helm     *HelmSpec          `yaml:"helm,omitempty"`
	Manifest *ManifestSpec      `yaml:"manifest,omitempty"`
	Exclude  *[]ExcludeItemSpec `yaml:"exclude,omitempty"`
}

type HelmSpec struct {
	ChartName  string `yaml:"chartName"`
	Version    string `yaml:"version"`
	RepoName   string `yaml:"repoName"`
	RepoUrl    string `yaml:"repoUrl"`
	ValuesFile string `yaml:"valuesFile"`
}

type ManifestSpec struct {
	Urls    []string
	Template map[string]string
}

type ExcludeItemSpec map[interface{}]interface{}

func (c *Config) AppType() AppType {
	if c.Spec.Helm != nil {
		return AppTypeHelm
	}
	if c.Spec.Manifest != nil {
		return AppTypeManifest
	}

	panic("unknown AppType due to bad config")
}

func (c *Config) HasExcludeSpec() bool {
	return c.Spec != nil && c.Spec.Exclude != nil
}