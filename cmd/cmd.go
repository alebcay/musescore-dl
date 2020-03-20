package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/alebcay/musescore-dl/internal/pkg"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	pdf "github.com/pdfcpu/pdfcpu/pkg/api"
)

var dest string

var rootCmd = &cobra.Command{
	Use:   "musescore-dl",
	Short: "musescore-to-PDF downloader",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		s.Suffix = " Creating temporary directory"
		tmp, err := ioutil.TempDir("", "musescore-dl-")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(tmp)

		err = msdl.SetupChrome(tmp)
		if err != nil {
			panic(err)
		}

		s.Suffix = " Getting score information"
		id, secret := msdl.GetScoreIDSecret(url)
		if id == "" || secret == "" {
			panic("bad score ID/secret")
		}

		pages, err := msdl.GetNumberOfPages(url)
		if id == "" || secret == "" {
		panic("bad page count")
		}

  		s.Suffix = " Downloading score"
		err = msdl.GetPages(id, secret, tmp, s, pages)
  		if err != nil {
  			panic(err)
  		}

  		s.Suffix = " Merging PDF files"
  		var pdfs []string
  		for i := 0; i < pages; i++ {
  		    pdfs = append(pdfs, fmt.Sprintf("%s/score_%d.pdf", tmp, i))
  		}

  		if dest == "" {
  			dest = msdl.GetScoreTitle(url) + ".pdf"
  		}

  		pdf.MergeFile(pdfs, dest, nil)
  		s.Stop()

  		result, err := filepath.Abs(dest)
  		if err != nil {
  			panic(err)
  		}

  		fmt.Printf("Wrote score to \"%s\"\n", result)
	},
}

func Execute() {
	rootCmd.Flags().StringVarP(&dest, "output", "o", "", "name of output file")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
