package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	logger "github.com/zkvl/opl/internal/log"
)

var wg sync.WaitGroup

func main() {
	print := flag.String("print", "", "Print log file or folder")
	runCmd := flag.Bool("runCmd", true, "If false, will log the command and pass the execution to shell environment")
	flag.Parse()

	if *print != "" {
		logger.PrintLogs(*print)
		return
	}

	if (len(os.Args) < 2 && *runCmd) || (len(os.Args) < 3 && !*runCmd) {
		fmt.Printf("Usage: %s [-print <log-file-or-folder>] command [args...]\n", os.Args[0])
		os.Exit(1)
	}

	wg.Add(2)

	go func() {
		defer wg.Done()
		logEntry := logger.NewEntry(*runCmd)
		logger.LogCommand(logEntry)
	}()

	if *runCmd {
		go executeCommand(os.Args[1:], *runCmd)
	} else {
		go executeCommand(os.Args[3:], *runCmd)
	}
	wg.Wait()
}

func executeCommand(command []string, runCmd bool) {
	defer wg.Done()

	if runCmd {
		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

		if err := cmd.Run(); err != nil {
			log.Println("[!] Error executing command:", err)
			os.Exit(1)
		}
	}
}
