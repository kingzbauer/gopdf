package main // import "browserless"

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"browserless/cdp"
	"browserless/server"
)

// type PDFWrapper struct {
// 	params *page.PrintToPDFParams
// 	data   []byte
// }

// func (p *PDFWrapper) Do(ctx context.Context, h cdp.Executor) error {
// 	var err error
// 	p.data, err = p.params.Do(ctx, h)
// 	return err
// }

// func mains() {
// 	ctx := context.Background()

// 	runnerOpts := []runner.CommandLineOption{
// 		runner.Flag("headless", true),
// 		runner.Flag("no-sandbox", true),
// 	}
// 	cdp, err := chromedp.New(
// 		ctx,
// 		chromedp.WithRunnerOptions(runnerOpts...),
// 	)
// 	if err != nil {
// 		fmt.Println("ERROR: ", err)
// 		os.Exit(1)
// 	}

// 	printToPDF := page.PrintToPDF()
// 	printToPDF = printToPDF.
// 		WithLandscape(false).
// 		WithPrintBackground(true).
// 		WithHeaderTemplate("Jack")
// 	pdf := &PDFWrapper{params: printToPDF}

// 	tasks := chromedp.Tasks{
// 		chromedp.Navigate("https://centrixt.mobisite.co.ke/invoice/1/doc/render/"),
// 		chromedp.WaitVisible("div", chromedp.ByQueryAll),
// 		pdf,
// 	}

// 	fmt.Println(cdp.Run(ctx, tasks))
// 	ioutil.WriteFile("test.pdf", pdf.data, 0777)
// 	fmt.Println(cdp.Shutdown(ctx))
// }

func main() {
	defer func() {
		fmt.Println("Shutting down...")
		fmt.Println(cdp.ShutdownCDP(context.Background()))
	}()
	s := server.InitServer(":8089")
	go func() {
		fmt.Println(s.ListenAndServe())
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	sig := <-c
	fmt.Println("Received", sig)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	fmt.Println(cdp.ShutdownCDP(ctx))
	fmt.Println("Shutting down server...")
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	s.Shutdown(ctx)
	fmt.Println("Server out...")
}
