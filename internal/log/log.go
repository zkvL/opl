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
		log.Println("[!] Error getting user home directory:", err)
		os.Exit(1)
	}
	logFolder = filepath.Join(homeDir, "operator-logs")

	err = os.MkdirAll(logFolder, os.ModePerm)
	if err != nil {
		log.Println("[!] Error creating log folder:", err)
		os.Exit(1)
	}
}

func Enable(shell string) {
	fishConfFileName := filepath.Join(os.Getenv("HOME") + "/.config/fish/config.fish")
	fishConfig := "\nfunction logCmd --on-event fish_prompt\n  set cmd $history[1]\n  opl \"$cmd\"\nend\n"
	zshConfFileName := filepath.Join(os.Getenv("HOME") + "/.zshrc")
	zshConfig := "\npreexec() { opl \"${1}\" }\n"

	switch shell {
	case "fish":
		fmt.Println("[-] Configuring fish environment")
		f, err := os.OpenFile(fishConfFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("[!] Error opening fish config file:", err)
		}
		defer f.Close()
		if _, err := f.WriteString(fishConfig); err != nil {
			log.Println("[!] Error configuring opl within fish config file:", err)
		}
	case "zsh":
		fmt.Println("[-] Configuring zsh environment")
		f, err := os.OpenFile(zshConfFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("[!] Error opening zsh config file:", err)
		}
		defer f.Close()
		if _, err := f.WriteString(zshConfig); err != nil {
			log.Println("[!] Error configuring opl within zsh config file:", err)
		}
	default:
		log.Printf("Shell %s not supported\n", shell)
		os.Exit(1)
	}
}

func Disable(shell string) {
	fishConfFileName := filepath.Join(os.Getenv("HOME") + "/.config/fish/config.fish")
	fishConfig := "\nfunction logCmd --on-event fish_prompt\n  set cmd $history[1]\n  opl \"$cmd\"\nend\n"
	zshConfFileName := filepath.Join(os.Getenv("HOME") + "/.zshrc")
	zshConfig := "\npreexec() { opl \"${1}\" }\n"

	switch shell {
	case "fish":
		fmt.Println("[-] Disabling fish environment")
		content, err := os.ReadFile(fishConfFileName)
		if err != nil {
			log.Println("[!] Error opening fish config file:", err)
		}
		fileContent := string(content)
		newContent := strings.Replace(fileContent, fishConfig, "", -1)

		err = os.WriteFile(fishConfFileName, []byte(newContent), os.ModePerm)
		if err != nil {
			log.Println("[!] Error removing fish configuration:", err)
		}
	case "zsh":
		fmt.Println("[-] Disabling zsh environment")
		content, err := os.ReadFile(zshConfFileName)
		if err != nil {
			log.Println("[!] Error opening zsh config file:", err)
		}
		fileContent := string(content)
		newContent := strings.Replace(fileContent, zshConfig, "", -1)

		err = os.WriteFile(zshConfFileName, []byte(newContent), os.ModePerm)
		if err != nil {
			log.Println("[!] Error removing zsh configuration:", err)
		}

	default:
		log.Printf("Shell %s not supported\n", shell)
		os.Exit(1)
	}
}

func LogCommand(entry *LogEntry) {
	// Open the log file or create it if it doesn't exist
	logFileName := filepath.Join(logFolder, entry.Date[:10]+".json")
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("[!] Error opening log file:", err)
		os.Exit(1)
	}
	defer logFile.Close()

	// Read existing entries or initialize as empty array
	var entries []LogEntry
	stat, err := logFile.Stat()
	if err != nil {
		log.Println("[!] Error getting file stat:", err)
		os.Exit(1)
	}

	if stat.Size() > 0 {
		decoder := json.NewDecoder(logFile)
		if err := decoder.Decode(&entries); err != nil && err != io.EOF {
			log.Println("[!] Error decoding log file:", err)
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
		log.Println("[!] Error encoding log entries:", err)
		os.Exit(1)
	}
}

func PrintLogs(path string) {
	header := fmt.Sprintf("%-25s %-20s %-20s %-20s\n", "Date", "IPAddr", "Operator", "Command")
	fmt.Printf("%s%s\n", header, strings.Repeat("-", len(header)))

	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("[!] Error accessing path:", err)
			return nil
		}

		if info.IsDir() {
			return nil
		}

		logFile, err := os.Open(filePath)
		if err != nil {
			log.Println("[!] Error opening log file:", err)
			return nil
		}
		defer logFile.Close()

		var entries []LogEntry
		decoder := json.NewDecoder(logFile)
		if err := decoder.Decode(&entries); err != nil {
			log.Println("[!] Error decoding log file:", err)
			return nil
		}

		for _, entry := range entries {
			fmt.Printf("%-25s %-20s %-20s %-20s\n", entry.Date, entry.IPAddr, entry.Operator, entry.Command)
		}
		fmt.Println()
		return nil
	})
}
