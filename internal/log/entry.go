package log

import (
	"os"
	"time"
)

var operator = os.Getenv("OPERATOR")
var date = time.Now().UTC().Format("2006-01-02 15:04:05 UTC")

type LogEntry struct {
	Date     string   `json:"date"`
	Activity string   `json:"activity"`
	IPAddr   []string `json:"ipaddr"`
	Operator string   `json:"operator,omitempty"`
}

func NewLogEntry(activity string, ips []string) *LogEntry {
	return &LogEntry{Date: date, Activity: activity, IPAddr: ips, Operator: operator}
}
