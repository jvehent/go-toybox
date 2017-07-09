package main

import (
	"log"
	"log/syslog"
)

func main() {
	slog, err := syslog.Dial(
		"udp",
		"localhost:514",
		syslog.LOG_LOCAL5|syslog.LOG_INFO,
		"SecuringDevOpsSyslog")
	defer slog.Close()
	if err != nil {
		log.Fatal("error:", err)
	}
	slog.Alert("This is an alert log")
	slog.Info("This is just information log")
}
