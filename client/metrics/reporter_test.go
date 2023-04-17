package metrics

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReporter_Run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	metricsCh := make(chan bool, 2)

	r := NewReporter(ctx, wg, time.Second, metricsCh)
	r.Run()

	metricsCh <- true
	metricsCh <- false

	cancel()
	close(metricsCh)
	wg.Wait()

	assert.Equal(t, r.success, 1)
	assert.Equal(t, r.fail, 1)
}
