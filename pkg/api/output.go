package api

type OutPut struct {
	name           string           `json:"name"`
	TestRunResults []TestRunResults `json:"test_run_results"`
}

type TestRunResults struct {
	TestRunName string       `json:"test_run_name"`
	Status      string       `json:"status"`
	Time        string       `json:"time"`
	CaseResults []CaseResult `json:"case_results"`
}
type CaseResult struct {
	CaseName string `json:"case_name"`
	Status   string `json:"status"`
	ErrorLog string `json:"error_log"`
	Time     string `json:"time"`
}
