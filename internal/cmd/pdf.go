package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/alebcay/musescore-dl/internal/pkg"
	"github.com/briandowns/spinner"
	pdf "github.com/pdfcpu/pdfcpu/pkg/api"
)

func DownloadPDF(url string, dest string) {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	defer s.Stop()

	s.Suffix = " Creating temporary directory"
	tmp, err := ioutil.TempDir("", "musescore-dl-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)

	s.Suffix = " Getting score information"
	id, secret := msdl.GetScoreIDSecret(url)
	if id == "" || secret == "" {
		panic("bad score ID/secret")
	}

	err = msdl.SetupChrome(tmp)
	if err != nil {
		panic(err)
	}

	pages, err := msdl.GetNumberOfPages(url)
	if id == "" || secret == "" {
		panic("bad page count")
	}

	s.Suffix = " Downloading score PDF"
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

	err = pdf.MergeFile(pdfs, dest, nil)

	if err != nil {
		panic(err)
	}
}
