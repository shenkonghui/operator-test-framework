package api

import (
	"time"
)

type TestJobs struct {
	Name      string      `yaml:"name"`
	Parameter []Parameter `yaml:"parameter"`
	Jobs      []Job       `yaml:"jobs"`
}

type Job struct {
	Name     string `yaml:"name"`
	Cmd      string `yaml:"cmd"`
	yamlFile string `yaml:"yamlFile"`
	// 执行测试后等待的时候，可能operator还没有开始工作
	InitTime   time.Duration `yaml:"initTime"`
	Timeout    time.Duration `yaml:"timeout"`
	Variable   []Variable    `yaml:"variable"`
	Verificate []Verificate  `yaml:"verificate"`
}
type Verificate struct {
	Name     string `yaml:"name"`
	JsonPath string `yaml:"jsonPath"`
	Cmd      string `yaml:"cmd"`
	Value    string `yaml:"value"`
	Operator string `yaml:"operator"`
}

type Parameter struct {
	Name  string `yaml:"key"`
	Value string `yaml:"value"`
}

type Variable struct {
	Key          string `yaml:"key"`
	Value        string `yaml:"value"`
	ValueFromCmd string `yaml:"valueFromCmd"`
}

const (
	OperatorEqual   = "equal"
	OperatorNoEqual = "noEqual"
)

const (
	StatusRunning   = "running"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
)
