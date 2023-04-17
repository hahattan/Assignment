package metrics

import (
	"context"
	"log"
	"sync"
	"time"
)

type Reporter struct {
	success   int
	fail      int
	metricsCh chan bool
	freq      time.Duration
	ctx       context.Context
	wg        *sync.WaitGroup
}

func NewReporter(ctx context.Context, wg *sync.WaitGroup, freq time.Duration, metricsCh chan bool) *Reporter {
	return &Reporter{metricsCh: metricsCh, freq: freq, ctx: ctx, wg: wg}
}

func (r *Reporter) Run() {
	r.wg.Add(1)
	go r.count()
	r.wg.Add(1)
	go r.report()
}

func (r *Reporter) count() {
	defer r.wg.Done()

	for b := range r.metricsCh {
		if b {
			r.success++
		} else {
			r.fail++
		}
	}
}

func (r *Reporter) report() {
	defer r.wg.Done()

	ticker := time.NewTicker(r.freq)
	defer ticker.Stop()

	for {
		select {
		case <-r.ctx.Done():
			return
		case <-ticker.C:
			// Report statistics about the number of succeeded calls and failed calls periodically
			log.Printf("Succeeded calls: %v\n", r.success)
			log.Printf("Failed calls: %v\n", r.fail)
		}
	}

}
