package entities

// Result01 is csv
type Result01 struct {
	JobName     string `csv:"job_name"`
	JobID       string `csv:"job_ID"`
	TaskID      string `csv:"task_ID"`
	Type        string `csv:"type"`
	StartTime   string `csv:"start_time"`
	FinishTime  string `csv:"finish_time"`
	ElapsedTime string `csv:"elapsed_time"`
	Hostname    string `csv:"host"`
}
