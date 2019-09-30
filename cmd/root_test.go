package cmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/plainbanana/olivine/entities"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "update .golden files")

func TestGetHistoryAPI(t *testing.T) {
	tests := []struct {
		testName string
		dest     string
		json     string
	}{
		{"valid01", "job", `{
			"jobs" : {
			   "job" : [
				  {
					 "submitTime" : 1326381344449,
					 "state" : "SUCCEEDED",
					 "user" : "user1",
					 "reducesTotal" : 1,
					 "mapsCompleted" : 1,
					 "startTime" : 1326381344489,
					 "id" : "job_1326381300833_1_1",
					 "name" : "word count",
					 "reducesCompleted" : 1,
					 "mapsTotal" : 1,
					 "queue" : "default",
					 "finishTime" : 1326381356010
				  },
				  {
					 "submitTime" : 1326381446500,
					 "state" : "SUCCEEDED",
					 "user" : "user1",
					 "reducesTotal" : 1,
					 "mapsCompleted" : 1,
					 "startTime" : 1326381446529,
					 "id" : "job_1326381300833_2_2",
					 "name" : "Sleep job",
					 "reducesCompleted" : 1,
					 "mapsTotal" : 1,
					 "queue" : "default",
					 "finishTime" : 1326381582106
				  }
			   ]
			}
		 }`},
		{"valid02", "task", `{
			"tasks" : {
			   "task" : [
				  {
					 "progress" : 100,
					 "elapsedTime" : 6777,
					 "state" : "SUCCEEDED",
					 "startTime" : 1326381446541,
					 "id" : "task_1326381300833_2_2_m_0",
					 "type" : "MAP",
					 "successfulAttempt" : "attempt_1326381300833_2_2_m_0_0",
					 "finishTime" : 1326381453318
				  },
				  {
					 "progress" : 100,
					 "elapsedTime" : 135559,
					 "state" : "SUCCEEDED",
					 "startTime" : 1326381446544,
					 "id" : "task_1326381300833_2_2_r_0",
					 "type" : "REDUCE",
					 "successfulAttempt" : "attempt_1326381300833_2_2_r_0_0",
					 "finishTime" : 1326381582103
				  }
			   ]
			}
		 }`},
		{"valid03", "taskattempt", `{
			"taskAttempt" : {
			   "assignedContainerId" : "container_1326381300833_0002_01_000002",
			   "progress" : 100,
			   "elapsedTime" : 2638,
			   "state" : "SUCCEEDED",
			   "diagnostics" : "",
			   "rack" : "/98.139.92.0",
			   "nodeHttpAddress" : "host.domain.com:8042",
			   "startTime" : 1326381450680,
			   "id" : "attempt_1326381300833_2_2_m_0_0",
			   "type" : "MAP",
			   "finishTime" : 1326381453318
			}
		 }`},
		{"valid04", "taskattemptcounter", `{
			"jobTaskAttemptCounters" : {
			   "taskAttemptCounterGroup" : [
				  {
					 "counterGroupName" : "org.apache.hadoop.mapreduce.FileSystemCounter",
					 "counter" : [
						{
						   "value" : 2363,
						   "name" : "FILE_BYTES_READ"
						},
						{
						   "value" : 54372,
						   "name" : "FILE_BYTES_WRITTEN"
						},
						{
						   "value" : 0,
						   "name" : "FILE_READ_OPS"
						},
						{
						   "value" : 0,
						   "name" : "FILE_LARGE_READ_OPS"
						},
						{
						   "value" : 0,
						   "name" : "FILE_WRITE_OPS"
						},
						{
						   "value" : 0,
						   "name" : "HDFS_BYTES_READ"
						},
						{
						   "value" : 0,
						   "name" : "HDFS_BYTES_WRITTEN"
						},
					   {
						   "value" : 0,
						   "name" : "HDFS_READ_OPS"
						},
						{
						   "value" : 0,
						   "name" : "HDFS_LARGE_READ_OPS"
						},
						{
						   "value" : 0,
						   "name" : "HDFS_WRITE_OPS"
						}
					 ]
				  },
				  {
					 "counterGroupName" : "org.apache.hadoop.mapreduce.TaskCounter",
					 "counter" : [
						{
						   "value" : 0,
						   "name" : "COMBINE_INPUT_RECORDS"
						},
						{
						   "value" : 0,
						   "name" : "COMBINE_OUTPUT_RECORDS"
						},
						{
						   "value" : 460,
						   "name" : "REDUCE_INPUT_GROUPS"
						},
						{
						   "value" : 2235,
						   "name" : "REDUCE_SHUFFLE_BYTES"
						},
						{
						   "value" : 460,
						   "name" : "REDUCE_INPUT_RECORDS"
						},
						{
						   "value" : 0,
						   "name" : "REDUCE_OUTPUT_RECORDS"
						},
						{
						   "value" : 0,
						   "name" : "SPILLED_RECORDS"
						},
						{
						   "value" : 1,
						   "name" : "SHUFFLED_MAPS"
						},
						{
						   "value" : 0,
						   "name" : "FAILED_SHUFFLE"
						},
						{
						   "value" : 1,
						   "name" : "MERGED_MAP_OUTPUTS"
						},
						{
						   "value" : 26,
						   "name" : "GC_TIME_MILLIS"
						},
						{
						   "value" : 860,
						   "name" : "CPU_MILLISECONDS"
						},
						{
						   "value" : 107839488,
						   "name" : "PHYSICAL_MEMORY_BYTES"
						},
						{
						   "value" : 1123147776,
						   "name" : "VIRTUAL_MEMORY_BYTES"
						},
						{
						   "value" : 57475072,
						   "name" : "COMMITTED_HEAP_BYTES"
						}
					 ]
				  },
				  {
					 "counterGroupName" : "Shuffle Errors",
					 "counter" : [
						{
						   "value" : 0,
						   "name" : "BAD_ID"
						},
						{
						   "value" : 0,
						   "name" : "CONNECTION"
						},
						{
						   "value" : 0,
						   "name" : "IO_ERROR"
						},
						{
						   "value" : 0,
						   "name" : "WRONG_LENGTH"
						},
						{
						   "value" : 0,
						   "name" : "WRONG_MAP"
						},
						{
						   "value" : 0,
						   "name" : "WRONG_REDUCE"
						}
					 ]
				  },
				  {
					 "counterGroupName" : "org.apache.hadoop.mapreduce.lib.output.FileOutputFormatCounter",
					 "counter" : [
						{
						   "value" : 0,
						   "name" : "BYTES_WRITTEN"
						}
					 ]
				  }
			   ],
			   "id" : "attempt_1326381300833_2_2_m_0_0"
			}
		 }`},
	}

	for _, tt := range tests {
		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("content-Type", "application/json")
				fmt.Fprintf(w, tt.json)
				return
			},
		))

		var items interface{}
		switch tt.dest {
		case "job":
			var a entities.GetJobs
			items = &a
		case "task":
			var a entities.GetTasks
			items = &a
		case "taskattempt":
			var a entities.GetTaskAttempt
			items = &a
		case "taskattemptcounter":
			var a entities.GetJobTaskAttemptCounters
			items = &a
		}

		res, err := getHistoryAPI(ts.URL, items)
		if err != nil {
			fmt.Println(ts.URL, res.Body, err)
		}

		fmt.Println(items)
		aa := fmt.Sprintln(items)
		gp := filepath.Join("..", "testdata", t.Name(), tt.testName+".golden")
		b := GetGolden(t, gp, []byte(aa))

		assert.Equal(t, string(aa), string(b), tt.testName)
		ts.Close()
	}
}

// GetGolden : read expected data from .golden file
func GetGolden(t *testing.T, goldenPath string, real []byte) []byte {
	if *update {
		t.Log("update golden file")
		if err := ioutil.WriteFile(goldenPath, real, 0644); err != nil {
			t.Fatalf("failed to update golden file: %s", err)
		}
	}
	g, err := ioutil.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed reading .golden: %s", err)
	}
	return g
}
