package orchannel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func sig(after time.Duration) <-chan any {
	c := make(chan any)
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func getChannelsFromDurations(d []time.Duration) []<-chan any {
	var out []<-chan any
	for _, t := range d {
		out = append(out, sig(t))
	}
	return out
}

var testCases = []struct {
	name            string
	durations       []time.Duration
	maxExpectedTime time.Duration
}{
	{
		name:            "no chan",
		durations:       []time.Duration{},
		maxExpectedTime: 2 * time.Millisecond,
	},
	{
		name: "single chan",
		durations: []time.Duration{
			1 * time.Second,
		},
		maxExpectedTime: time.Second + 100*time.Millisecond,
	},
	{
		name: "basic",
		durations: []time.Duration{
			1 * time.Second,
			2 * time.Second,
			3 * time.Second,
		},
		maxExpectedTime: time.Second + 100*time.Millisecond,
	},
	{
		name: "basic 2",
		durations: []time.Duration{
			2 * time.Hour,
			5 * time.Minute,
			1 * time.Second,
			1 * time.Hour,
			1 * time.Minute,
		},
		maxExpectedTime: time.Second + 100*time.Millisecond,
	},
}

func TestOrGoroutineApproach(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			chans := getChannelsFromDurations(tc.durations)
			start := time.Now()
			<-orGoroutineApproach(chans...)
			execTime := time.Since(start)
			assert.GreaterOrEqual(t, tc.maxExpectedTime, execTime)
		})
	}
}

func TestOrReflectApproach(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			start := time.Now()
			chans := getChannelsFromDurations(tc.durations)
			<-orReflectionApproach(chans...)
			execTime := time.Since(start)
			assert.GreaterOrEqual(t, tc.maxExpectedTime, execTime)
		})
	}
}
