package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	fd, err := os.OpenFile("log.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fd.Close()

	pid := os.Getpid()
	fmt.Printf("running as pid %d\n", pid)
	fmt.Fprintf(fd, "running as pid %d\n", pid)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch)
	val, ok := <-ch
	if !ok {
		fmt.Println("channel closed")
		fmt.Fprintln(fd, "channel closed")
	} else {
		fmt.Println(val)
		fmt.Fprintln(fd, val)

		time.Sleep(10 * time.Second)

		fmt.Println("10s elapsed")
		fmt.Fprintln(fd, "10s elapsed")
	}
}
