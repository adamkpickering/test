package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fd, err := os.OpenFile("log.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fd.Close()

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	val, ok := <-ctx.Done()
	if ok {
		fmt.Printf("val: %v\n", val)
		fmt.Fprintf(fd, "val: %v\n", val)
		os.Exit(1)
	} else {
		msg := fmt.Sprintf("context done at %s", time.Now().Format(time.RFC3339))
		fmt.Println(msg)
		fmt.Fprintln(fd, msg)
	}
	time.Sleep(60 * time.Second)
	msg := fmt.Sprintf("exiting at %s", time.Now().Format(time.RFC3339))
	fmt.Println(msg)
	fmt.Fprintln(fd, msg)
}
