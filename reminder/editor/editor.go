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
- help: see this
- list: list all active jobs`
	botPrefix = "reminder "
	filePath  = "../timetable.csv"
)

// ChatMsg is converted to JSON and POSTed to hook
type ChatMsg struct {
	Text string `json:"text"`
}

var routes = map[string]func([]string) string{
	"help": help,
	"list": list}

func help(words []string) string {
	log.Println("help")
	log.Println(words[1:])
	return helpText
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
	return routes[words[0]](words)
}

func serveStatus(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("I'm ok!"))
}

func extractWords(text string) []string {
	return strings.SplitN(strings.TrimLeft(text, botPrefix), " ", 2)
}

func serveBot(w http.ResponseWriter, r *http.Request) {
	l, _ := httputil.DumpRequest(r, true)
	log.Println(string(l))
	msg := &ChatMsg{
		route(extractWords(r.FormValue("text")))}
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
