package log

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var logFolder string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println("[!] Error getting user home directory: ", err)
		os.Exit(1)
	}
	logFolder = filepath.Join(homeDir, "operator-logs")

	err = os.MkdirAll(logFolder, os.ModePerm)
	if err != nil {
		log.Println("[!] Error creating log folder: ", err)
		os.Exit(1)
	}
}

func filter(cmd string) bool {
	filters := []string{"alias", "cd", "chmod", "chown", "cp", "exit", "find", "id", "kill", "ls", "locate", "make", "man", "mkdir", "mv", "nano", "opl", "ps", "pwd", "uname", "vim", "which", "whoami"}

	inputCmd := strings.Fields(cmd)
	if len(inputCmd) > 0 {
		first := inputCmd[0]
		for _, fCmd := range filters {
			if strings.HasPrefix(first, fCmd) {
				return true
			}
		}
	}
	return false
}

func LogCmd(entry *LogEntry) {
	if filter(entry.Command) {
		return
	}

	// Open the log file or create it if it doesn't exist
	logFileName := filepath.Join(logFolder, entry.Date[:10]+".json")
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("[!] Error opening log file: ", err)
		os.Exit(1)
	}
	defer logFile.Close()

	// Read existing entries or initialize as empty array
	var entries []LogEntry
	stat, err := logFile.Stat()
	if err != nil {
		log.Println("[!] Error getting file stat: ", err)
		os.Exit(1)
	}

	if stat.Size() > 0 {
		decoder := json.NewDecoder(logFile)
		if err := decoder.Decode(&entries); err != nil && err != io.EOF {
			log.Println("[!] Error decoding log file: ", err)
			os.Exit(1)
		}
	}

	// Append the new entry
	entries = append(entries, *entry)

	// Write all entries back to the file
	logFile.Seek(0, 0)
	logFile.Truncate(0)
	encoder := json.NewEncoder(logFile)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(entries); err != nil {
		log.Println("[!] Error encoding log entries: ", err)
		os.Exit(1)
	}
}

func PrintLogs(path string) {
	header := fmt.Sprintf("%-20s %-25s %-20s %-20s\n", "Operator", "Timestamp (UTC)", "Operator IP", "Command/Activity")
	fmt.Printf("%s%s\n", header, strings.Repeat("-", len(header)))

	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("[!] Error accessing path: ", err)
			return nil
		}

		if info.IsDir() {
			return nil
		}

		logFile, err := os.Open(filePath)
		if err != nil {
			log.Println("[!] Error opening log file: ", err)
			return nil
		}
		defer logFile.Close()

		var entries []LogEntry
		decoder := json.NewDecoder(logFile)
		if err := decoder.Decode(&entries); err != nil {
			log.Println("[!] Error decoding log file: ", err)
			return nil
		}

		for _, entry := range entries {
			fmt.Printf("%-20s %-25s %-20s %-20s\n", entry.Operator, entry.Date, entry.IPAddr, entry.Command)
		}
		fmt.Println()
		return nil
	})
}
