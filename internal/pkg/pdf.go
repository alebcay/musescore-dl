package msdl

import (
	"bytes"
	"context"
	"io"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func GeneratePDF(url string, dest string, width float64, height float64) error {
	paper_width := (width / 96.0) + 2.0
	paper_height := (height / 96.0) + 2.0

	// create context
	ctx = GetChromeContext()

	var pdfReader io.Reader
	err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitReady("svg"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPaperWidth(paper_width).
				WithPaperHeight(paper_height).
				WithMarginTop(1.0).
				WithMarginBottom(1.0).
				WithMarginLeft(1.0).
				WithMarginRight(1.0).
				WithPageRanges("1").
				Do(ctx)
			if err != nil {
				return err
			}
			pdfReader = bytes.NewBuffer(buf)
			return nil
		}),
	})

	if err != nil {
		return err
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, pdfReader)
	if err != nil {
		return err
	}

	return nil
}
