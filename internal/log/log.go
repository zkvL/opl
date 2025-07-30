package log

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	externalip "github.com/glendc/go-external-ip"
	"github.com/xuri/excelize/v2"
)

var logFolder string

// Initialize oplogs folder
func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println("[!] Error getting user home directory: ", err)
		os.Exit(1)
	}
	logFolder = filepath.Join(homeDir, ".oplogs")

	err = os.MkdirAll(logFolder, os.ModePerm)
	if err != nil {
		log.Println("[!] Error creating log folder: ", err)
		os.Exit(1)
	}
}

// Helper function to filter out common commands we don't want to log
// Reduces noise
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

// Helper function to quesry the public IP address from where opl is executed
func GetPublicIP() (net.IP, error) {
	consensus := externalip.DefaultConsensus(nil, nil)
	consensus.UseIPProtocol(4)

	return consensus.ExternalIP()
}

// Logs activity
func LogActivity(entry *LogEntry, debug bool) {
	if filter(entry.Activity) {
		return
	}

	// Open the log file or create it if it doesn't exist
	logFileName := filepath.Join(logFolder, entry.Date[:10]+".json")
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		if debug {
			log.Println("[!] Error opening log file: ", err)
		}
		os.Exit(1)
	}
	defer logFile.Close()

	// Read existing entries or initialize as empty array
	var entries []LogEntry
	stat, err := logFile.Stat()
	if err != nil {
		if debug {
			log.Println("[!] Error getting file stat: ", err)
		}
		os.Exit(1)
	}

	if stat.Size() > 0 {
		decoder := json.NewDecoder(logFile)
		if err := decoder.Decode(&entries); err != nil && err != io.EOF {
			if debug {
				log.Println("[!] Error decoding log file: ", err)
			}
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
		if debug {
			log.Println("[!] Error encoding log entries: ", err)
		}
		os.Exit(1)
	}
}

// Function to format header and table widths
func fmtTable(entries []LogEntry) ([]string, []int) {
	// Prepare headers
	headers := []string{"Operator", "Operator IP(s)", "Timestamp (UTC)", "Command/Activity"}

	// Calculate initial column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}

	// Recalculate column with based on largest data in entries
	for _, entry := range entries {
		if len(entry.Operator) > widths[0] {
			widths[0] = len(entry.Operator)
		}
		ipStr := strings.Join(entry.IPAddr, ", ")
		if len(ipStr) > widths[1] {
			widths[1] = len(ipStr)
		}
		if len(entry.Date) > widths[2] {
			widths[2] = len(entry.Date)
		}
		if len(entry.Activity) > widths[3] {
			widths[3] = len(entry.Activity)
		}
	}

	return headers, widths
}

// Save logs to an XLSX
func fmtXLSX(entries []LogEntry) {
	filename := "opl-timeline.xlsx"
	f := excelize.NewFile()
	sheet := "Timeline"
	f.SetSheetName(f.GetSheetName(0), sheet)

	// Write headers
	headers := []string{"Operator", "Operator IP", "Timestamp (UTC)", "Command/Activity"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Add all entries gathered from json files
	for rowIdx, entry := range entries {
		row := rowIdx + 2
		values := []string{
			entry.Operator,
			strings.Join(entry.IPAddr, ", "),
			entry.Date,
			entry.Activity,
		}
		for colIdx, val := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
			f.SetCellValue(sheet, cell, val)
		}
	}

	if err := f.SaveAs(filename); err == nil {
		fmt.Printf("[+] Logs written successfully to file: %s\n", filename)
	}
}

// Prints all logs in MD format
func fmtMD(entries []LogEntry) {
	headers, widths := fmtTable(entries)

	// Print header and separator
	format := fmt.Sprintf("| %%-%ds | %%-%ds | %%-%ds | %%-%ds |\n", widths[0], widths[1], widths[2], widths[3])
	separator := fmt.Sprintf("|-%s-|-%s-|-%s-|-%s-|\n",
		strings.Repeat("-", widths[0]),
		strings.Repeat("-", widths[1]),
		strings.Repeat("-", widths[2]),
		strings.Repeat("-", widths[3]),
	)

	fmt.Printf(format, headers[0], headers[1], headers[2], headers[3])
	fmt.Print(separator)

	// Add all entries gathered from json files
	for _, entry := range entries {
		ipStr := strings.Join(entry.IPAddr, ", ")
		fmt.Printf(format, entry.Operator, ipStr, entry.Date, entry.Activity)
	}
	fmt.Println()
}

// Print all logs to terminal (default)
func fmtTerminal(entries []LogEntry) {
	headers, widths := fmtTable(entries)

	// Print header and separator
	format := fmt.Sprintf("%%-%ds %%-%ds %%-%ds %%-%ds\n", widths[0], widths[1], widths[2], widths[3])
	fmt.Printf(format, headers[0], headers[1], headers[2], headers[3])
	fmt.Printf("%s\n", strings.Repeat("-", (widths[0]+widths[1]+widths[2]+widths[3])))

	// Add all entries gathered from json files
	for _, entry := range entries {
		ipStr := strings.Join(entry.IPAddr, ", ")
		fmt.Printf(format, entry.Operator, ipStr, entry.Date, entry.Activity)
	}
	fmt.Println()
}

// Parses JSON log files to a LogEntry array
func files2logs(path string, debug bool) []LogEntry {
	var logs []LogEntry

	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			if debug {
				log.Printf("[!] Error accessing path %s: %v\n", path, err)
			}
			return nil
		}

		if info.IsDir() {
			return nil
		}

		logFile, err := os.Open(filePath)
		if err != nil {
			if debug {
				log.Printf("[!] Error opening log file %s: %v\n", filePath, err)
			}
			return nil
		}
		defer logFile.Close()

		var entries []LogEntry
		decoder := json.NewDecoder(logFile)
		if err := decoder.Decode(&entries); err != nil {
			if debug {
				log.Printf("[!] %s is not a valid opl log file: %v\n", filePath, err)
			}
			return nil
		}
		logs = append(logs, entries...)
		return nil
	})

	return logs
}

func ShowLogs(path string, format string, debug bool) {
	logs := files2logs(path, debug)

	switch format {
	case "md":
		fmtMD(logs)
	case "xlsx":
		fmtXLSX(logs)
	default:
		fmtTerminal(logs)
	}
}
