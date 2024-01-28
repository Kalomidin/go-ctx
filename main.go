package main

import (
	"ctx/pkg/context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	ctxWithCancel, cancel := context.WithCancel(ctx)
	now := time.Now()
	cancel()
	<-ctxWithCancel.Done()
	fmt.Println("timer finished", time.Since(now))
}
