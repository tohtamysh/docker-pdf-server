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
	http.HandleFunc("/", handle)

	fmt.Printf("Starting server...\n")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Sorry, only POST method are support.")
		return
	}

	url := r.FormValue("url")
	html := r.FormValue("html")

	pdfParams := PDFParams{OrientationLandscape: false}
	if r.FormValue("landscape") == "1" {
		pdfParams.OrientationLandscape = true
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte

	if url == "" && html == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Params error")
		return
	} else if url != "" {
		buf = urlToPDF(ctx, url, pdfParams)
	} else if html != "" {
		buf = htmlToPDF(ctx, html, pdfParams)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/pdf")
	w.Write(buf)
	return
}

func urlToPDF(ctx context.Context, url string, pdfParams PDFParams) []byte {
	var buf []byte

	if err := chromedp.Run(ctx, getURL(url, pdfParams, &buf)); err != nil {
		log.Fatal(err)
	}
	return buf
}

func htmlToPDF(ctx context.Context, html string, pdfParams PDFParams) []byte {
	var buf []byte
	if err := chromedp.Run(ctx, getHTML(html, pdfParams, &buf)); err != nil {
		log.Fatal(err)
	}

	return buf
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
