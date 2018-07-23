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
