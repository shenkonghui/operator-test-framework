package api


type TestJobs struct {
	Name string `yaml:"name"`
	Jobs []Job `yaml:"jobs"`
}

type Job struct {
	yamlFile string `yaml:"yamlFile"`
	MaxWaitTime string `yaml:"maxWaitTime"`
	Verificate []Verificate  `yaml:"verificate"`

}
type Verificate struct {
	JsonPath string `yaml:"jsonPath"`
	Value string `yaml:"value"`
}