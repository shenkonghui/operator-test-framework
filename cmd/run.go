/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"operator-test-framework/pkg/api"
	job "operator-test-framework/pkg/api"
	"operator-test-framework/pkg/util"
	"os"
	"strings"
	"time"

	"github.com/prometheus/tsdb/fileutil"

	"k8s.io/apimachinery/pkg/util/json"

	"k8s.io/klog/v2"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configPath string
var parameter string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run test job",
	RunE: func(cmd *cobra.Command, args []string) error {
		klog.V(2).Info("start to test operator")
		hostname, _ := os.Hostname()
		jobs := &job.TestJobs{Name: fmt.Sprintf("%s_%s", hostname, time.Now().Format("20060102150405"))}
		ouput := api.OutPut{Name: jobs.Name}

		// 判断文件存在
		exist, err := util.PathExists(configPath)
		if !exist {
			return fmt.Errorf("path is not exist: %s", configPath)
		}
		if err != nil {
			return err
		}

		// 支持目录和文件
		isDir := util.IsDir(configPath)
		files := []string{}
		if isDir {
			files, err = fileutil.ReadDir(configPath)
		} else {
			files = append(files, configPath)
		}

		// 遍历所有文件
		for _, cfgFile := range files {
			file := cfgFile
			if isDir {
				file = configPath + "/" + cfgFile
			}
			// 判断下级是否是目录
			if util.IsDir(file) {
				continue
			} else {
				yamlFile, _ := ioutil.ReadFile(file)
				err := yaml.Unmarshal(yamlFile, jobs)
				if err != nil {
					return err
				}
				jobs.Parameter = util.ConvertStrToPara(parameter, jobs.Parameter)
				if len(jobs.Parameter) > 0 {
					yamlFileStr := string(yamlFile)
					for _, para := range jobs.Parameter {
						if para.Name != "" {
							name := fmt.Sprintf("#{%s}", para.Name)
							yamlFileStr = strings.ReplaceAll(yamlFileStr, name, para.Value)
						}
					}
					err := yaml.Unmarshal([]byte(yamlFileStr), jobs)
					if err != nil {
						return err
					}

				}
				klog.V(2).Info("parsed the job file successfully")
				startTime := time.Now()
				err, testResult, runResult := runJob(jobs)
				if err != nil {
					klog.V(2).Info(err.Error())
					runResult.Status = api.StatusFailed
				}
				if testResult {
					runResult.Status = api.StatusCompleted
				}
				runResult.Time = time.Now().Sub(startTime).String()
				ouput.TestRunResults = append(ouput.TestRunResults, runResult)

			}
		}

		defer func() {
			result, _ := json.Marshal(ouput)
			fmt.Println(string(result))
			if ouput.TestRunResults[0].Status == api.StatusFailed {
				os.Exit(2)
			}
		}()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.PersistentFlags().StringVarP(&configPath, "configPath", "c", "", "配置文件路径，可以指定目录")
	rootCmd.PersistentFlags().StringVarP(&parameter, "parameter", "p", "", "指定变量，将替换配置文件中的配置, 比如name1=abc,name2=bcd")

}

func runJob(jobs *job.TestJobs) (error, bool, api.TestRunResults) {
	testRun := api.TestRunResults{TestRunName: jobs.Name}
	for i, job := range jobs.Jobs {
		klog.V(2).Infof("run job %s", job.Name)
		testRun.CaseResults = append(testRun.CaseResults, api.CaseResult{CaseName: job.Name})
		testRun.CaseResults[i].CaseName = job.Name

		shell := job.Cmd

		Parameter := make(map[string]string)

		// 设置变量
		if len(job.Variable) != 0 {
			for _, p := range job.Variable {
				if p.ValueFromCmd != "" {
					cmdResult, err := util.ExecCmd(p.ValueFromCmd)
					if err != nil {
						testRun.Status = api.StatusFailed
						return err, false, testRun
					}
					Parameter[p.Key] = string(cmdResult)
				} else {
					Parameter[p.Key] = string(p.Value)
				}
			}
		}

		// 执行测试命令
		if shell != "" {
			_, err := util.ExecCmd(shell)
			if err != nil {
				testRun.CaseResults[i].ErrorLog = err.Error()
				testRun.Status = api.StatusFailed
				return err, false, testRun
			}

		}

		// 等待初始化时间
		if job.InitTime == 0 {
			job.InitTime = time.Second * 2
		}
		time.Sleep(job.InitTime)
		startTime := time.Now()
		if job.Timeout == 0 {
			job.Timeout = time.Second * 60
		}
		endTime := startTime.Add(job.Timeout)

		klog.V(2).Infof("start to verificate job: %v", job.Name)
		// 开始执行验证
		for {
			// 超时
			if time.Now().After(endTime) {
				testRun.CaseResults[i].Time = time.Now().Sub(startTime).String()
				testRun.CaseResults[i].Status = "failed"
				testRun.CaseResults[i].ErrorLog = fmt.Sprintf("timeout: %s", testRun.CaseResults[i].ErrorLog)
				klog.V(2).Infof("verificate job %v timeout", job.Name)
				return fmt.Errorf("timeout"), false, testRun
			}

			// 验证
			verificateTrue, err, caseResults := verificate(job.Verificate, Parameter, testRun.CaseResults[i])
			testRun.CaseResults[i] = caseResults
			if err != nil {
				errLog := fmt.Sprintf("verificate err: %s", err.Error())
				testRun.CaseResults[i].ErrorLog = errLog
				return err, false, testRun
			}

			// 验证成功
			if verificateTrue {
				klog.V(2).Infof("verificate job %v success,", job.Name)
				testRun.CaseResults[i].Time = time.Now().Sub(startTime).String()
				testRun.CaseResults[i].Status = "completed"
				break
			}
			if job.Interval == 0 {
				job.Interval = time.Second
			}
			time.Sleep(job.Interval)
		}
	}
	return nil, true, testRun
}

// 验证
func verificate(verifes []job.Verificate, parameter map[string]string, caseResult api.CaseResult) (bool, error, api.CaseResult) {
	for _, verife := range verifes {
		cmdResult, err := util.ExecCmd(verife.Cmd)
		if err != nil {
			caseResult.ErrorLog = err.Error()
			return false, err, caseResult
		}
		if len(verife.Value) == 0 {
			return false, fmt.Errorf("value is null"), caseResult
		}

		// 替换变量
		var value string
		if ([]byte(verife.Value)[0]) == '@' {
			value = parameter[string([]byte(verife.Value)[1:])]
		} else {
			value = verife.Value
		}
		switch verife.Operator {
		// 验证操作，目前只包含等于/不等于
		case api.OperatorNoEqual:
			if string(cmdResult) == value {
				caseResult.Status = api.StatusFailed
				caseResult.ErrorLog = fmt.Sprintf("%v is equal %v", string(cmdResult), value)
				klog.V(2).Info(caseResult.ErrorLog)
				return false, nil, caseResult
			} else {
				caseResult.ErrorLog = ""
				caseResult.Status = api.StatusCompleted
			}
		case api.OperatorEqual:
			if string(cmdResult) != value {
				caseResult.Status = api.StatusFailed
				caseResult.ErrorLog = fmt.Sprintf("%v is not equal %v", string(cmdResult), value)
				klog.V(2).Info(caseResult.ErrorLog)
				return false, nil, caseResult
			} else {
				caseResult.Status = api.StatusCompleted
				caseResult.ErrorLog = ""
			}
		default:
			if string(cmdResult) != value {
				caseResult.Status = api.StatusFailed
				caseResult.ErrorLog = fmt.Sprintf("%v is not equal %v", string(cmdResult), value)
				klog.V(2).Info(caseResult.ErrorLog)
				return false, nil, caseResult
			} else {
				caseResult.Status = api.StatusCompleted
				caseResult.ErrorLog = ""
			}
		}

	}
	return true, nil, caseResult

}
