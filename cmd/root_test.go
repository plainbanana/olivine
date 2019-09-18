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
