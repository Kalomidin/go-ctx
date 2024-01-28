package main

import (
	"ctx/pkg/context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	ctxWithDeadline, cancel := context.WithDeadline(ctx, time.Now().Add(50*time.Second))
	now := time.Now()
	cancel()
	<-ctxWithDeadline.Done()
	fmt.Println("timer canceled", time.Since(now))
}
