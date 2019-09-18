package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/plainbanana/hadoop-jobhistoryfetcher/entities"
	"github.com/spf13/cobra"
)

var (
	config entities.Config
)

func init() {
	cobra.OnInitialize(initHosts)
	RootCmd.PersistentFlags().StringVar(&config.Hostfile, "hostfile", "", "Specify target hosts from a hostfile. default target is localhost.")
	RootCmd.PersistentFlags().StringVarP(&config.TargetPort, "port", "p", "19888", "Specify the port where target hadoop job history server running on hosts.")
}

// RootCmd : test
var RootCmd = &cobra.Command{
	Use:   "olivine",
	Short: "A command to fetch hadoop job histories.",
	Run:   rootcmd,
}

func initHosts() {
	if config.Hostfile != "" {
		fp, err := os.Open(config.Hostfile)
		if err != nil {
			log.Fatal(err)
		}
		defer fp.Close()

		scanner := bufio.NewScanner(fp)

		for scanner.Scan() {
			if host := strings.TrimSpace(scanner.Text()); host != "" {
				config.Hosts = append(config.Hosts, host)
			}
		}
	} else {
		config.Hosts = append(config.Hosts, "localhost")
	}
}

func getHistoryAPI(url string, destInterface interface{}) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	bodyBytes, _ := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bodyBytes, destInterface)

	return res, err
}

func rootcmd(cmd *cobra.Command, args []string) {
	uris := struct {
		jobs     string
		tasks    string
		attempts string
	}{
		jobs:     "jobs",
		tasks:    "tasks",
		attempts: "attempts",
	}

	// call GET: /jobs
	mJobs := make(map[string]entities.GetJobs)

	for _, host := range config.Hosts {
		var getJobs entities.GetJobs

		targetURI := "http://" + host + ":" + config.TargetPort + "/ws/v1/history/mapreduce/" + uris.jobs
		_, err := getHistoryAPI(targetURI, &getJobs)
		if err == nil && len(getJobs.Jobs.Job) != 0 {
			mJobs[targetURI] = getJobs
		} else if err != nil {
			log.Println("Error: GET ", targetURI, err)
		}

		// fmt.Println("getjobs", getJobs, err)
	}
	log.Println("found job at ", len(mJobs), "host(s)")

	// call GET: /jobs/{JobID}/tasks
	mTasks := make(map[string]entities.GetTasks)

	for uri, getJobs := range mJobs {
		for _, job := range getJobs.Jobs.Job {
			var getTask entities.GetTasks

			targetURI := strings.Join([]string{uri, job.ID, uris.tasks}, "/")
			getTask.JobID = job.ID
			getTask.JobName = job.Name
			getTask.MapsTotal = job.MapsTotal
			getTask.ReducesTotal = job.ReducesTotal
			_, err := getHistoryAPI(targetURI, &getTask)
			if err == nil && len(getTask.Tasks.Task) != 0 {
				mTasks[targetURI] = getTask
				log.Println("job", job.ID, job.Name, "has", len(getTask.Tasks.Task), "task(s)", getTask.MapsTotal, "maps", getTask.ReducesTotal, "reduces")
			} else if err != nil {
				log.Println("Error: GET", targetURI, err)
			}

			// fmt.Println("response1", getTask, "ERR", err)
		}
	}

	// fmt.Println("map get", mTasks)
	// call GET: /jobs/{JobID}/tasks/{TaskID}/attempts/{TaskAttemptID}
	var csv []*entities.Result01

	for uri, getTasks := range mTasks {
		for _, task := range getTasks.Tasks.Task {
			var attempt entities.GetTaskAttempt

			targetURI := strings.Join([]string{uri, task.ID, uris.attempts, task.SuccessfulAttempt}, "/")
			attempt.JobID = getTasks.JobID
			attempt.JobName = getTasks.JobName
			attempt.TaskID = task.ID
			// fmt.Println("attemptLog", attempt)
			_, err := getHistoryAPI(targetURI, &attempt)
			if err == nil {
				csvRow := entities.Result01{
					JobName:     attempt.JobName,
					JobID:       attempt.JobID,
					TaskID:      attempt.TaskID,
					StartTime:   fmt.Sprint(attempt.TaskAttempt.StartTime),
					FinishTime:  fmt.Sprint(attempt.TaskAttempt.FinishTime),
					ElapsedTime: fmt.Sprint(attempt.TaskAttempt.ElapsedTime),
					Hostname:    attempt.TaskAttempt.NodeHTTPAddress,
					Type:        attempt.TaskAttempt.Type,
				}

				csv = append(csv, &csvRow)
				// fmt.Println("csvrow", csvRow)
			} else if err != nil {
				log.Println("ERROR: GET", targetURI, err)
			}

			// fmt.Println("response2", attempt, err)
		}
	}

	csvContent, err := gocsv.MarshalString(&csv)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(csvContent)
}
