package main

import (
	"encoding/csv"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/robfig/cron"
)

// CsvLoader provides means to load csv rows from a file
type CsvLoader struct {
	Address string
	rows    [][]string
}

// Load returns csv rows and whether they differ from last fetch
func (l *CsvLoader) Load() ([][]string, bool) {
	file, err := os.Open(l.Address)
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

func main() {
	loader := &CsvLoader{Address: "timetable.csv"}
	c := &cron.Cron{}

	for true {
		csv, differs := loader.Load()

		log.Println("Checking cronjobs")
		if differs {
			log.Println("Loading cronjobs")
			c.Stop()
			c = cron.New()
			for _, task := range csv {
				cronErr := c.AddFunc(task[0], func() { log.Println(task[1]) })
				if cronErr != nil {
					panic(cronErr)
				}
			}
			c.Start()
		}

		time.Sleep(5 * time.Second)
	}
}
