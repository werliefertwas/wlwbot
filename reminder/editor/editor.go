package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/robfig/cron"
)

const (
	helpText = `
    reminder help: see this
    reminder insert "*/10 * * * * *","Every 10 seconds!": activate new job
    reminder list: list all active jobs with index
    reminder remove 0: remove job at index 0`
	filePath = "../timetable.csv"
)

// ChatMsg is converted to JSON and POSTed to hook
type ChatMsg struct {
	Text string `json:"text"`
}

var routes = map[string]func([]string) string{
	"help":   help,
	"list":   list,
	"insert": insert,
	"remove": remove}

func help(words []string) string {
	return helpText
}

func insert(words []string) string {
	row, readerErr := csv.NewReader(strings.NewReader(words[1])).Read()
	if readerErr != nil {
		return readerErr.Error()
	}
	if len(row) != 2 {
		return "row needs 2 columns"
	}
	_, cronErr := cron.Parse(row[0])
	if cronErr != nil {
		return cronErr.Error()
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	writer := csv.NewWriter(file)
	writerErr := writer.Write(row)
	if writerErr != nil {
		return writerErr.Error()
	}
	writer.Flush()

	return "inserted"
}

func list(words []string) string {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	rows, csvErr := csv.NewReader(file).ReadAll()
	if csvErr != nil {
		panic(csvErr)
	}

	joinedRows := make([]string, len(rows), len(rows))
	for i, row := range rows {
		joinedRows[i] = "    " + strconv.Itoa(i) + ": " + strings.Join(row, ", ")
	}

	return strings.Join(joinedRows, "\n")
}

func remove(words []string) string {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	rows, readerErr := csv.NewReader(file).ReadAll()
	if readerErr != nil {
		panic(readerErr)
	}

	deleteAt, iErr := strconv.Atoi(words[1])
	if iErr != nil {
		panic(iErr)
	}

	msg := "not found"

	keptRows := [][]string{}
	for i, row := range rows {
		if i != deleteAt {
			keptRows = append(keptRows, row)
		} else {
			msg = "removed"
		}
	}

	file.Seek(0, 0)
	file.Truncate(0)
	writerErr := csv.NewWriter(file).WriteAll(keptRows)
	if writerErr != nil {
		panic(writerErr)
	}
	return msg
}

func route(words []string) string {
	return routes[words[0]](words)
}

func serveStatus(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("I'm ok!"))
}

func extractWords(text string) []string {
	return strings.SplitN(text, " ", 3)[1:]
}

func serveBot(w http.ResponseWriter, r *http.Request) {
	l, _ := httputil.DumpRequest(r, true)
	log.Println(string(l))
	t := r.FormValue("text")
	msg := &ChatMsg{route(extractWords(t))}
	msgJSON, _ := json.Marshal(msg)
	w.Write([]byte(msgJSON))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port
	r := mux.NewRouter()
	r.Path("/status").HandlerFunc(serveStatus)
	r.Path("/").HandlerFunc(serveBot)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(port, nil))
}
