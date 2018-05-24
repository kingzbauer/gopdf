package cdp

import (
	"context"
	"sync"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/client"
	"github.com/chromedp/chromedp/runner"
	"github.com/fatih/color"
)

var (
	CDP    *chromedp.CDP
	lock   sync.Mutex
	Client *client.Client
	Pool   *chromedp.Pool
)

type PDFOptions struct {
	URL             string  `json:"url"`
	PrintBackground bool    `json:"print_background"`
	Landscape       bool    `json:"landscape"`
	PaperSize       string  `json:"paper_size"`
	Scale           float64 `json:"scale"`
	// TargetElement is the element browser will wait to be visible before
	// it attempts to generate pdf
	TargetElement     string `json:"target_element"`
	TargetElementType string `json:"target_element_type"` // e.g class or id
}

type PaperSize struct {
	Width, Height float64
}

var PaperSizes map[string]PaperSize

func init() {
	PaperSizes = map[string]PaperSize{
		"A0": {33.1, 46.8},
		"A1": {23.4, 33.1},
		"A2": {16.5, 23.4},
		"A3": {11.7, 16.5},
		"A4": {8.3, 11.7},
		"A5": {5.8, 8.3},
		"A6": {4.1, 5.8},
		"A7": {2.9, 4.1},
	}
}

var log = color.New(color.FgRed)
var whiteLog = color.New(color.FgWhite)

func logger(str string, values ...interface{}) {
	whiteLog.Printf(str+"\n", values)
}

func GenPDF(ctx context.Context, pdfOptions *PDFOptions) ([]byte, error) {
	// get the paper size
	paperSize := PaperSizes["A4"]
	if size, ok := PaperSizes[pdfOptions.PaperSize]; ok {
		paperSize = size
	}

	printToPDF := page.PrintToPDF()
	printToPDF = printToPDF.
		WithLandscape(pdfOptions.Landscape).
		WithPrintBackground(pdfOptions.PrintBackground).
		WithPaperWidth(paperSize.Width).
		WithPaperHeight(paperSize.Height).
		WithScale(pdfOptions.Scale)

	// set the query option
	queryOption := chromedp.ByQueryAll
	if pdfOptions.TargetElementType == "id" {
		queryOption = chromedp.ByID
	}

	p, err := Pool.Allocate(ctx, runnerOpts()...)
	defer p.Release()
	if err != nil {
		return nil, err
	}
	var data []byte
	tasks := chromedp.Tasks{
		chromedp.Navigate(pdfOptions.URL),
		chromedp.WaitVisible(pdfOptions.TargetElement, queryOption),
		pdfWrapperFunc(&data, printToPDF),
	}
	err = p.Run(ctx, tasks)
	return data, err
}

// InitCDP creates a Pool resource for chrome runner instances
func InitCDP(ctx context.Context) (err error) {
	Pool, err = chromedp.NewPool()
	return
}

// ShutdownCDP shuts down Pool releasing the resources
func ShutdownCDP(ctx context.Context) error {
	return Pool.Shutdown()
}

func runnerOpts() []runner.CommandLineOption {
	return []runner.CommandLineOption{
		runner.ExecPath("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"),
		runner.Flag("headless", true),
		runner.Flag("no-sandbox", true),
	}
}

// ActionFunc implements the chromedp.Action interface
type ActionFunc func(context.Context, cdp.Executor) error

func (af ActionFunc) Do(ctx context.Context, h cdp.Executor) error {
	return af(ctx, h)
}

func pdfWrapperFunc(result *[]byte, printParams *page.PrintToPDFParams) chromedp.Action {
	return ActionFunc(func(ctx context.Context, h cdp.Executor) error {
		data, err := printParams.Do(ctx, h)
		if err == nil {
			*result = data
		}
		return err
	})
}
