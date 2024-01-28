package context_test

import (
	"ctx/pkg/context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithDeadline(t *testing.T) {
	ctx := context.Background()
	ctxWithDeadline, cancel := context.WithDeadline(ctx, time.Now().Add(50*time.Millisecond))
	now := time.Now()
	<-ctxWithDeadline.Done()
	cancel()
	assert.LessOrEqual(t, time.Since(now), 1*time.Second)
}

func TestWithDeadlineFirstCancel(t *testing.T) {
	ctx := context.Background()
	ctxWithDeadline, cancel := context.WithDeadline(ctx, time.Now().Add(50*time.Second))
	now := time.Now()
	cancel()
	<-ctxWithDeadline.Done()
	assert.LessOrEqual(t, time.Since(now), 1*time.Second)
}

func TestWithCancel(t *testing.T) {
	ctx := context.Background()
	ctxWithCancel, cancel := context.WithCancel(ctx)
	now := time.Now()
	cancel()
	<-ctxWithCancel.Done()

	assert.LessOrEqual(t, time.Since(now), 1*time.Second)
}

func TestWithParentDeadline(t *testing.T) {
	ctx := context.Background()
	ctxWithDeadline, parentCancel := context.WithDeadline(ctx, time.Now().Add(1*time.Second))
	now := time.Now()
	ctxWithTimeout, cancel := context.WithTimeout(ctxWithDeadline, 1*time.Minute)
	<-ctxWithTimeout.Done()
	<-ctxWithDeadline.Done()
	parentCancel()
	cancel()

	assert.LessOrEqual(t, time.Since(now), 2*time.Second)
}

func TestWithParentCancel(t *testing.T) {
	ctx := context.Background()
	ctxWithDeadline, parentCancel := context.WithDeadline(ctx, time.Now().Add(10*time.Second))
	now := time.Now()
	ctxWithTimeout, cancel := context.WithTimeout(ctxWithDeadline, 1*time.Minute)
	parentCancel()
	<-ctxWithDeadline.Done()
	<-ctxWithTimeout.Done()
	cancel()

	assert.LessOrEqual(t, time.Since(now), 1*time.Second)
}
