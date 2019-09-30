package entities

import (
	"encoding/json"
	"log"
	"strconv"
)

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
			StartTime         int64       `json:"startTime"`
			FinishTime        int64       `json:"finishTime"`
			ElapsedTime       int         `json:"elapsedTime"`
			Progress          json.Number `json:"progress"`
			ID                string      `json:"id"`
			State             string      `json:"state"`
			Type              string      `json:"type"`
			SuccessfulAttempt string      `json:"successfulAttempt"`
		} `json:"task"`
	} `json:"tasks"`
	JobID        string
	JobName      string
	ReducesTotal int
	MapsTotal    int
}

// GetTaskAttempts is struct
type GetTaskAttempts struct {
	TaskAttempts struct {
		TaskAttempt []struct {
			StartTime           int64       `json:"startTime"`
			FinishTime          int64       `json:"finishTime"`
			ElapsedTime         int         `json:"elapsedTime"`
			Progress            json.Number `json:"progress"`
			ID                  string      `json:"id"`
			Rack                string      `json:"rack"`
			State               string      `json:"state"`
			Status              string      `json:"status"`
			NodeHTTPAddress     string      `json:"nodeHttpAddress"`
			Diagnostics         string      `json:"diagnostics"`
			Type                string      `json:"type"`
			AssignedContainerID string      `json:"assignedContainerId"`
		} `json:"taskAttempt"`
	} `json:"taskAttempts"`
	JobID   string
	JobName string
	TaskID  string
}

// GetTaskAttempt is struct
type GetTaskAttempt struct {
	TaskAttempt struct {
		AssignedContainerID string      `json:"assignedContainerId"`
		Progress            json.Number `json:"progress"`
		ElapsedTime         int         `json:"elapsedTime"`
		State               string      `json:"state"`
		Diagnostics         string      `json:"diagnostics"`
		Rack                string      `json:"rack"`
		NodeHTTPAddress     string      `json:"nodeHttpAddress"`
		StartTime           int64       `json:"startTime"`
		ID                  string      `json:"id"`
		Type                string      `json:"type"`
		FinishTime          int64       `json:"finishTime"`
	} `json:"taskAttempt"`
	JobID   string
	JobName string
	TaskID  string
}

// GetJobTaskAttemptCountersInterface is interface
type GetJobTaskAttemptCountersInterface interface {
	ToResult01(*Result01)
}

// GetJobTaskAttemptCounters is struct
type GetJobTaskAttemptCounters struct {
	JobTaskAttemptCounters struct {
		TaskAttemptCounterGroup []struct {
			CounterGroupName string `json:"counterGroupName"`
			Counter          []struct {
				Value int    `json:"value"`
				Name  string `json:"name"`
			} `json:"counter"`
		} `json:"taskAttemptCounterGroup"`
		ID string `json:"id"`
	} `json:"jobTaskAttemptCounters"`
}

// ToResult01 : map to result01 csv struct
func (c GetJobTaskAttemptCounters) ToResult01(in *Result01) {
	for _, v := range c.JobTaskAttemptCounters.TaskAttemptCounterGroup {
		for _, vv := range v.Counter {
			val := vv.Value
			name := vv.Name

			switch name {
			case "FILE_BYTES_READ":
				(*in).FileBytesRead = strconv.Itoa(val)
			case "FILE_BYTES_WRITTEN":
				(*in).FileBytesWritten = strconv.Itoa(val)
			case "FILE_READ_OPS":
				(*in).FileReadOps = strconv.Itoa(val)
			case "FILE_LARGE_READ_OPS":
				(*in).FileLargeReadOps = strconv.Itoa(val)
			case "FILE_WRITE_OPS":
				(*in).FileWriteOps = strconv.Itoa(val)
			case "HDFS_BYTES_READ":
				(*in).HDFSBytesRead = strconv.Itoa(val)
			case "HDFS_BYTES_WRITTEN":
				(*in).HDFSBytesWritten = strconv.Itoa(val)
			case "HDFS_READ_OPS":
				(*in).HDFSReadOps = strconv.Itoa(val)
			case "HDFS_LARGE_READ_OPS":
				(*in).HDFSLargeReadOps = strconv.Itoa(val)
			case "HDFS_WRITE_OPS":
				(*in).HDFSWriteOps = strconv.Itoa(val)
			case "COMBINE_INPUT_RECORDS":
				(*in).CombineInputRecords = strconv.Itoa(val)
			case "COMBINE_OUTPUT_RECORDS":
				(*in).CombineOutputRecords = strconv.Itoa(val)
			case "REDUCE_INPUT_GROUPS":
				(*in).ReduceInputGroups = strconv.Itoa(val)
			case "REDUCE_SHUFFLE_BYTES":
				(*in).ReduceShuffleBytes = strconv.Itoa(val)
			case "REDUCE_INPUT_RECORDS":
				(*in).ReduceInputRecords = strconv.Itoa(val)
			case "REDUCE_OUTPUT_RECORDS":
				(*in).ReduceOutputRecords = strconv.Itoa(val)
			case "SPILLED_RECORDS":
				(*in).SpilledRecords = strconv.Itoa(val)
			case "SHUFFLED_MAPS":
				(*in).ShuffledMaps = strconv.Itoa(val)
			case "FAILED_SHUFFLE":
				(*in).FailedShuffle = strconv.Itoa(val)
			case "MERGED_MAP_OUTPUTS":
				(*in).MergedMapOutputs = strconv.Itoa(val)
			case "GC_TIME_MILLIS":
				(*in).GCTimeMillis = strconv.Itoa(val)
			case "CPU_MILLISECONDS":
				(*in).CPUMilliseconds = strconv.Itoa(val)
			case "PHYSICAL_MEMORY_BYTES":
				(*in).PhysicalMemoryBytes = strconv.Itoa(val)
			case "VIRTUAL_MEMORY_BYTES":
				(*in).VirtualMemoryBytes = strconv.Itoa(val)
			case "COMMITTED_HEAP_BYTES":
				(*in).CommitedHeapBytes = strconv.Itoa(val)
			case "BAD_ID":
				(*in).BadID = strconv.Itoa(val)
			case "CONNECTION":
				(*in).Connection = strconv.Itoa(val)
			case "IO_ERROR":
				(*in).IOError = strconv.Itoa(val)
			case "WRONG_LENGTH":
				(*in).WrongLength = strconv.Itoa(val)
			case "WRONG_MAP":
				(*in).WrongMap = strconv.Itoa(val)
			case "WRONG_REDUCE":
				(*in).WrongReduce = strconv.Itoa(val)
			case "BYTES_WRITTEN":
				(*in).BytesWritten = strconv.Itoa(val)
			default:
				log.Println("warn:", v.CounterGroupName, name, "is undefined")
			}
		}
	}
}
