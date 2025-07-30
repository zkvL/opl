package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	logger "github.com/zkvL/opl/internal/log"
)

var wg sync.WaitGroup

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "opl",
	Short: "Yet another operator logging tool for Red Teamers",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		if activity, _ := cmd.Flags().GetString("act"); activity != "" {
			wg.Add(1)
			go func() {
				defer wg.Done()

				var logEntry *logger.LogEntry

				if ips, _ := cmd.Flags().GetString("ip"); ips != "" {
					ipsArray := strings.Split(ips, ",")
					logEntry = logger.NewLogEntry(activity, ipsArray)
				} else {
					ip, err := logger.GetPublicIP()
					if err != nil {
						log.Println("[!] Error getting the public IP")
					}
					logEntry = logger.NewLogEntry(activity, []string{ip.String()})
				}
				logger.LogActivity(logEntry, debug)
			}()
			wg.Wait()
		} else if show, _ := cmd.Flags().GetBool("show"); show {
			format, _ := cmd.Flags().GetString("format")
			location, _ := cmd.Flags().GetString("location")
			logger.ShowLogs(location, format, debug)
		} else {
			cmd.Help()
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	homeDir, _ := os.UserHomeDir()
	logFolder := filepath.Join(homeDir, ".oplogs")

	rootCmd.Flags().StringP("act", "a", "", "Logs an activity")
	rootCmd.Flags().StringP("ip", "i", "", "Comma-separated IPs to log from where the activity was performed. If not set, public IP will be logged")

	rootCmd.Flags().BoolP("show", "s", false, "Print logs from the default location or any specified with -l")
	rootCmd.Flags().StringP("location", "l", logFolder, "Individual file or folder name to print logs from")
	rootCmd.Flags().StringP("format", "f", "", "Print logs in specified format (currently supported md and xlsx)")

	rootCmd.Flags().BoolP("debug", "d", false, "Print debug error messages")
	rootCmd.Flags().SortFlags = false
}
