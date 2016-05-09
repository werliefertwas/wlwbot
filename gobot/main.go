package main

import (
  "fmt"
  "net/http"
  "os"
)

func main() {
    http.HandleFunc("/gobot", handler)
    http.ListenAndServe(":" + os.Getenv("PORT"), nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "{\"token\":\"8cehrcw8mtfd9cee1uau575twc\",\"username\": \"gobot\", \"text\":\"Folgende branches sind auf den Sandboxen %s!\"}", r.PostFormValue("text"))
}
