package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	logger "github.com/zkvL/opl/internal/log"
)

var wg sync.WaitGroup

func main() {
	cmd := flag.Bool("cmd", false, "If set, will log the executed command along with the host public IP")
	print := flag.String("print", "", "Print log file or folder")
	flag.Parse()

	if *print != "" {
		logger.PrintLogs(*print)
		return
	}

	if len(os.Args) < 2 {
		fmt.Printf("Usage: \t%s [-cmd] 'cmd or activity to log'\n\t%s [-print <log-file-or-folder>]\n", os.Args[0], os.Args[0])
		os.Exit(1)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		logEntry := logger.NewEntry(*cmd)
		logger.LogCmd(logEntry)
	}()
	wg.Wait()
}
