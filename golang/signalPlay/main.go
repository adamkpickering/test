package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func doStep1() error {
	fmt.Println("doing step 1")
	time.Sleep(3 * time.Second)
	fmt.Println("step 1 done")
	return nil
}

func doStep2() error {
	fmt.Println("doing step 2")
	time.Sleep(4 * time.Second)
	fmt.Println("step 2 done")
	// return fmt.Errorf("hello step %d", 2)
	return nil
}

func doStep3() error {
	fmt.Println("doing step 3")
	time.Sleep(2 * time.Second)
	fmt.Println("step 3 done")
	return nil
}

func checkContextBetween(ctx context.Context, taskChan <-chan func() error, errChan chan<- error) {
	for task := range taskChan {
		select {
		case <-ctx.Done():
			close(errChan)
			return
		default:
			if err := task(); err != nil {
				errChan <- err
				close(errChan)
				return
			}
		}
	}
	close(errChan)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	taskChan := make(chan func() error, 10)
	errChan := make(chan error)
	go checkContextBetween(ctx, taskChan, errChan)
	taskChan <- doStep1
	taskChan <- doStep2
	taskChan <- doStep3
	close(taskChan)
	err, ok := <-errChan
	fmt.Printf("err: %s\nok: %t", err, ok)
	if ok {
		fmt.Printf("error: %s\n", err)
	}
}
