package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var dest string
var format string

var rootCmd = &cobra.Command{
	Use:   "musescore-dl",
	Short: "musescore-to-PDF downloader",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		switch format {
		case "pdf":
			DownloadPDF(url, dest)
		case "midi":
			DownloadMIDI(url, dest)
		case "mxl":
			DownloadMXL(url, dest)
		default:
			fmt.Println("invalid format specified")
			os.Exit(1)
		}
	},
}

func Execute() {
	rootCmd.Flags().StringVarP(&dest, "output", "o", "", "name of output file")
	rootCmd.Flags().StringVarP(&format, "format", "f", "pdf", "format of download (pdf, midi, or mxl)")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
