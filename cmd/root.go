package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/comail/colog"
	"github.com/gocarina/gocsv"
	"github.com/plainbanana/olivine/entities"
	"github.com/spf13/cobra"
)

const (
	version = "0.3.1"
	// FEachHost is flag
	FEachHost int = 1 << iota
)

var (
	config entities.Config
	// ErrInputRange : err
	ErrInputRange = errors.New("error: input out of range")
)

func init() {
	cobra.OnInitialize(initHosts)

	RootCmd.AddCommand(versionCmd)

	RootCmd.PersistentFlags().StringVar(&config.LogLevel, "log", "info", "Specify olivine minimun log level. {trace, debug, info, warn, error, alert}")
	RootCmd.Flags().StringVar(&config.Hostfile, "hostfile", "", "Specify target hosts from a hostfile. default target is localhost.")
	RootCmd.Flags().StringVarP(&config.TargetPort, "port", "p", "19888", "Specify the port where target hadoop job history server running on hosts.")

	RootCmd.AddCommand(plotCmd)
	plotCmd.Flags().StringVarP(&config.FileInput, "csv", "c", "olivine.csv", "Specify input csv filepath.")
	plotCmd.Flags().StringVarP(&config.FileOutput, "save", "s", "histories.png", "Specify output image filename.")
}

// RootCmd : test
var RootCmd = &cobra.Command{
	Use:              "olivine",
	Short:            "A command to fetch hadoop job histories.",
	PersistentPreRun: setLogMinLevel,
	Run:              rootcmd,
}

func setLogMinLevel(cmd *cobra.Command, args []string) {
	colog.SetDefaultLevel(colog.LDebug)
	colog.SetMinLevel(getLogMinLevel())
	colog.Register()
	log.Println("debug: colog level", config.LogLevel)
}

func getLogMinLevel() colog.Level {
	switch config.LogLevel {
	case "trace":
		return colog.LTrace
	case "debug":
		return colog.LDebug
	case "info":
		return colog.LInfo
	case "warn":
		return colog.LWarning
	case "error":
		return colog.LError
	case "alert":
		return colog.LAlert
	default:
		return colog.LDebug
	}
}

func initHosts() {
	if config.Hostfile != "" {
		fp, err := os.Open(config.Hostfile)
		if err != nil {
			log.Fatal("error: ", err)
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
		log.Fatal("alert: ", err)
	}

	bodyBytes, _ := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("alert: ", err)
	}

	err = json.Unmarshal(bodyBytes, destInterface)

	return res, err
}

func rootcmd(cmd *cobra.Command, args []string) {
	uris := struct {
		jobs     string
		tasks    string
		attempts string
		counters string
	}{
		jobs:     "jobs",
		tasks:    "tasks",
		attempts: "attempts",
		counters: "counters",
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
			log.Println("alert: GET ", targetURI, err)
		}

		log.Println("trace: jobs ", getJobs, err)
	}
	log.Println("info: found job at ", len(mJobs), "host(s)")

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
				log.Println("info: job", job.ID, job.Name, "has", len(getTask.Tasks.Task), "task(s)", getTask.MapsTotal, "maps", getTask.ReducesTotal, "reduces")
			} else if err != nil {
				log.Println("alert: GET", targetURI, err)
			}

			log.Println("trace: tasks ", getTask, " ERROR is ", err)
		}
	}

	// fmt.Println("map get", mTasks)
	// call GET: /jobs/{JobID}/tasks/{TaskID}/attempts/{TaskAttemptID}
	// call GET: /jobs/{JobID}/tasks/{TaskID}/attempts/{TaskAttemptID}/counters
	var csv []*entities.Result01

	for uri, getTasks := range mTasks {
		for _, task := range getTasks.Tasks.Task {
			var attempt entities.GetTaskAttempt
			var counters entities.GetJobTaskAttemptCounters

			targetURI := strings.Join([]string{uri, task.ID, uris.attempts, task.SuccessfulAttempt}, "/")
			attempt.JobID = getTasks.JobID
			attempt.JobName = getTasks.JobName
			attempt.TaskID = task.ID
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

				targetCountURI := strings.Join([]string{targetURI, uris.counters}, "/")
				_, err = getHistoryAPI(targetCountURI, &counters)
				if err != nil {
					log.Println("alert: GET", targetCountURI, err)
				} else {
					counters.ToResult01(&csvRow)
				}

				csv = append(csv, &csvRow)
				log.Println("trace: csvrow ", csvRow)
			} else if err != nil {
				log.Println("alert: GET", targetURI, err)
			}

			log.Println("trace: response2", attempt, err)
		}
	}

	csvContent, err := gocsv.MarshalString(&csv)
	if err != nil {
		log.Println("alert: failed to marshall string to csv.")
		log.Fatal(err)
	}
	log.Println("info: success to create csv.")

	fmt.Println(csvContent)
}
