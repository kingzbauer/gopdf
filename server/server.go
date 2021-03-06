package server

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"browserless/cdp"

	"github.com/getsentry/raven-go"
	"github.com/go-http-utils/logger"
)

var (
	server         *http.Server
	Mux            *http.ServeMux
	DefaultTimeout time.Duration = time.Second * 20
)

func init() {
	if err := cdp.InitCDP(context.TODO()); err != nil {
		panic(err)
	}

	//
}

func InitServer(addr string) *http.Server {
	Mux = http.NewServeMux()
	server = &http.Server{
		Handler: logger.Handler(Mux, os.Stdout, logger.CombineLoggerType),
		Addr:    addr,
	}

	initViews(Mux)

	return server
}

func initViews(mux *http.ServeMux) {
	mux.Handle("/generate/pdf/", WrapHandlerFunc(
		RequirePOST,
		raven.RecoveryHandler(GeneratePDF),
	))
}

// GeneratePDF receives a pdf request from a particular url, generates, and
// returns the pdf file content
func GeneratePDF(w http.ResponseWriter, r *http.Request) {
	opts := &cdp.PDFOptions{}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(opts); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), DefaultTimeout)
	defer cancel()
	data, err := cdp.GenPDF(ctx, opts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment")
	w.Write(data)
}
