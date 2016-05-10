package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

const (
	helpText = `
    reminder help: see this
    reminder list: list all active jobs
    reminder insert "*/10 * * * * *","Every 10 seconds!": activate new job`
	filePath = "../timetable.csv"
)

// ChatMsg is converted to JSON and POSTed to hook
type ChatMsg struct {
	Text string `json:"text"`
}

var routes = map[string]func([]string) string{
	"help":   help,
	"list":   list,
	"insert": insert}

func help(words []string) string {
	log.Println("help")
	return helpText
}

func insert(words []string) string {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0600)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	row := words[1] + "\n"

	file.WriteString(row)

	return "added"
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
		joinedRows[i] = "    " + strings.Join(row, ", ")
	}

	return strings.Join(joinedRows, "\n")
}

func route(words []string) string {
	log.Println(words)
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
