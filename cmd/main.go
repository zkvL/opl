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
	print := flag.String("print", "", "Print log file or folder")
	noip := flag.Bool("noip", false, "If true, won't log the IP address of the current machine.")
	enable := flag.String("enable", "", "Configure shell environment to log commands")
	disable := flag.String("disable", "", "Disable shell environment from logging commands")
	flag.Parse()

	if *enable != "" {
		logger.Enable(*enable)
		return
	}

	if *disable != "" {
		logger.Disable(*disable)
		return
	}

	if *print != "" {
		logger.PrintLogs(*print)
		return
	}

	if len(os.Args) < 2 {
		fmt.Printf("Usage: \t%s [-print <log-file-or-folder>]\n\t%s [-noip] 'actiity-to-log'\n", os.Args[0], os.Args[0])
		os.Exit(1)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		logEntry := logger.NewEntry(*noip)
		logger.LogCommand(logEntry)
	}()
	wg.Wait()
}
