package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"log"
	"net/http"
)

type PDFParams struct {
	OrientationLandscape bool
}

func main() {
	http.HandleFunc("/url", urlToPDF)

	http.HandleFunc("/html", htmlToPDF)

	fmt.Printf("Starting server...\n")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func urlToPDF(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")

	pdfParams := PDFParams{OrientationLandscape: false}
	if r.FormValue("landscape") == "1" {
		pdfParams.OrientationLandscape = true
	}

	if url != "" {
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		// capture pdf
		var buf []byte
		if err := chromedp.Run(ctx, getURL(url, pdfParams, &buf)); err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/pdf")
		w.Write(buf)
		return
	}
}

func htmlToPDF(w http.ResponseWriter, r *http.Request) {
	html := r.FormValue("html")

	pdfParams := PDFParams{OrientationLandscape: false}
	if r.FormValue("landscape") == "1" {
		pdfParams.OrientationLandscape = true
	}

	if html != "" {
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		// capture pdf
		var buf []byte
		if err := chromedp.Run(ctx, getHTML(html, pdfParams, &buf)); err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/pdf")
		w.Write(buf)
		return
	}
}

func getURL(url string, params PDFParams, res *[]byte) chromedp.Tasks {
	fmt.Println("Go to ", url)
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithLandscape(params.OrientationLandscape).WithPrintBackground(false).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}

func getHTML(html string, params PDFParams, res *[]byte) chromedp.Tasks {
	fmt.Println("Parse html")
	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithLandscape(params.OrientationLandscape).WithPrintBackground(false).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
