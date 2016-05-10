package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/robfig/cron"
)

const hook = "http://localhost:8065/hooks/pdusr3nmwfn4pyhrh83dixr1xo"
const filePath = "timetable.csv"

// ChatMsg is converted to JSON and POSTed to hook
type ChatMsg struct {
	Text string `json:"text"`
}

// CsvLoader provides means to load headerless csv rows from a file
type CsvLoader struct {
	FilePath string
	rows     [][]string
}

// Load returns csv rows and whether they differ from last fetch
func (l *CsvLoader) Load() ([][]string, bool) {
	file, err := os.Open(l.FilePath)
	if err != nil {
		panic(err)
	}

	rows, csvErr := csv.NewReader(file).ReadAll()
	if csvErr != nil {
		panic(csvErr)
	}

	differs := !reflect.DeepEqual(rows, l.rows)
	l.rows = rows

	return rows, differs
}

func remind(text string) {
	jsonMsg, _ := json.Marshal(&ChatMsg{text})
	msgReader := bytes.NewReader(jsonMsg)
	client := &http.Client{}
	r, _ := client.Post(hook, "application/json", msgReader)
	// l, _ := httputil.DumpResponse(r, false)
	log.Println(r.StatusCode, r.Header["X-Request-Id"][0], string(jsonMsg))
}

func main() {
	loader := &CsvLoader{FilePath: filePath}
	cronJobs := &cron.Cron{}

	for true {
		csv, differs := loader.Load()

		if differs {
			cronJobs.Stop()
			cronJobs = cron.New()
			for _, task := range csv {
				cronErr := cronJobs.AddFunc(task[0], func() { remind(task[1]) })
				if cronErr != nil {
					panic(cronErr)
				}
			}
			cronJobs.Start()
		}

		time.Sleep(1 * time.Second)
	}
}
