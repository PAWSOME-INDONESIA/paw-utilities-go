package browser

import (
	"context"

	"github.com/PAWSOME-INDONESIA/paw-utilities-go/logs"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type (
	Browser interface {
		RenderPDF(parent context.Context, path string) ([]byte, error)
	}

	browser struct {
		log logs.Logger
	}
)

func (b *browser) RenderPDF(parent context.Context, path string) ([]byte, error) {
	var buf []byte

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		parent,
		chromedp.WithLogf(b.log.Printf),
	)
	defer cancel()

	// render html to pdf file
	err := chromedp.Run(ctx,
		chromedp.Navigate("file:///"+path),
		chromedp.ActionFunc(func(ctx context.Context) error {
			result, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}
			buf = result
			return nil
		}),
	)

	if err != nil {
		return nil, err
	}
	return buf, nil
}

func NewBrowser(log logs.Logger) Browser {
	return &browser{log}
}
