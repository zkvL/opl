package log

import (
	"log"
	"net"
	"os"
	"time"

	externalip "github.com/glendc/go-external-ip"
)

type LogEntry struct {
	Date     string `json:"date"`
	Command  string `json:"command"`
	IPAddr   string `json:"ipaddr"`
	Operator string `json:"operator,omitempty"`
}

func NewEntry(cmd bool) *LogEntry {
	operator := os.Getenv("OPERATOR")
	date := time.Now().UTC().Format("2006-01-02 15:04:05 GMT")

	if cmd {
		command := os.Args[2]
		ip, err := getPublicIP()
		if err != nil {
			log.Println("[!] Error getting the public IP")
		}
		return &LogEntry{Date: date, Command: command, IPAddr: ip.String(), Operator: operator}
	}

	command := os.Args[1]
	return &LogEntry{Date: date, Command: command, IPAddr: "", Operator: operator}
}

func getPublicIP() (net.IP, error) {
	consensus := externalip.DefaultConsensus(nil, nil)
	consensus.UseIPProtocol(4)

	return consensus.ExternalIP()
}
