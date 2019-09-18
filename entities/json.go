package entities

// GetJobs is struct
type GetJobs struct {
	Jobs struct {
		Job []struct {
			SubmitTime       int64  `json:"submitTime"`
			State            string `json:"state"`
			User             string `json:"user"`
			ReducesTotal     int    `json:"reducesTotal"`
			MapsCompleted    int    `json:"mapsCompleted"`
			StartTime        int64  `json:"startTime"`
			ID               string `json:"id"`
			Name             string `json:"name"`
			ReducesCompleted int    `json:"reducesCompleted"`
			MapsTotal        int    `json:"mapsTotal"`
			Queue            string `json:"queue"`
			FinishTime       int64  `json:"finishTime"`
		} `json:"job"`
	} `json:"jobs"`
}

// GetTasks is struct
type GetTasks struct {
	Tasks struct {
		Task []struct {
			StartTime         int64  `json:"startTime"`
			FinishTime        int64  `json:"finishTime"`
			ElapsedTime       int    `json:"elapsedTime"`
			Progress          int    `json:"progress"`
			ID                string `json:"id"`
			State             string `json:"state"`
			Type              string `json:"type"`
			SuccessfulAttempt string `json:"successfulAttempt"`
		} `json:"task"`
	} `json:"tasks"`
}

// GetTaskAttempts is struct
type GetTaskAttempts struct {
	TaskAttempts struct {
		TaskAttempt []struct {
			StartTime           int64  `json:"startTime"`
			FinishTime          int64  `json:"finishTime"`
			ElapsedTime         int    `json:"elapsedTime"`
			Progress            int    `json:"progress"`
			ID                  string `json:"id"`
			Rack                string `json:"rack"`
			State               string `json:"state"`
			Status              string `json:"status"`
			NodeHTTPAddress     string `json:"nodeHttpAddress"`
			Diagnostics         string `json:"diagnostics"`
			Type                string `json:"type"`
			AssignedContainerID string `json:"assignedContainerId"`
		} `json:"taskAttempt"`
	} `json:"taskAttempts"`
}
