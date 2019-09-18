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
	if err, ok := err.(*json.SyntaxError); ok {
		fmt.Println(string(bodyBytes[err.Offset-15 : err.Offset+15]))
	}

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
	fmt.Println(uris)

	// call GET: /jobs
	mJobs := make(map[string]entities.GetJobs)

	for _, host := range config.Hosts {
		var getJobs entities.GetJobs

		targetURI := "http://" + host + ":" + config.TargetPort + "/ws/v1/history/mapreduce/" + uris.jobs
		res, err := getHistoryAPI(targetURI, &getJobs)
		if err == nil {
			mJobs[host] = getJobs
		}

		fmt.Println("getjobs", res.Body, getJobs, err)
	}
	log.Println("found job at ", len(mJobs), "host(s)")

	// call GET: /jobs/{JobID}/tasks
	mTasks := make(map[string]entities.GetTasks)

	for host, getJobs := range mJobs {
		for _, job := range getJobs.Jobs.Job {
			var getTask entities.GetTasks

			targetURI := "http://" + host + ":" + config.TargetPort + "/ws/v1/history/mapreduce/" + uris.jobs + "/" + job.ID + "/tasks"
			res, err := getHistoryAPI(targetURI, &getTask)
			if err == nil {
				mTasks[host] = getTask
				log.Println("job", job.ID, job.Name, "has", len(getTask.Tasks.Task), "task(s) at", host)
			}

			fmt.Println("response", res.Body, getTask)
		}
	}

	fmt.Println("root command", config.Hostfile, config.Hosts)
}
