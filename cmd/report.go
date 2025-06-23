/*
Copyright © 2025 Alireza Arzehgar <alirezaarzehgar82@gmail.com>
*/
package cmd

import (
	"log"

	"github.com/alirezaarzehgar/iflashc/internal/reporter"
	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Report bug or feature request with creating new github issue",
	Run: func(cmd *cobra.Command, args []string) {
		if !reporter.Open() {
			log.Println("Failed to open issue tracker. Please report your bug through ", reporter.URL)
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}
