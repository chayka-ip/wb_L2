package orchannel

import (
	"reflect"
	"sync"
)

//Or traces for moment when one of the channels passed will be closed.
//Channels considered to be "done-channels" which purpose is indicate
//that job was finished. Because of that reason the data they send is ignored.
//Result channel is closed by default.
func Or(channels ...<-chan interface{}) <-chan interface{} {
	bGoroutne := !true
	if bGoroutne {
		return orGoroutineApproach(channels...)
	}
	return orReflectionApproach(channels...)
}

func orGoroutineApproach(channels ...<-chan any) <-chan any {
	outChan := make(chan any)

	var syncExit sync.Once
	var exitFunc = func() {
		close(outChan)
	}

	if len(channels) < 1 {
		exitFunc()
	}

	write := func(ch <-chan any) {
		for v := range ch {
			_ = v
		}
		syncExit.Do(exitFunc)
	}

	for _, ch := range channels {
		go write(ch)
	}

	<-outChan
	return outChan
}

func orReflectionApproach(channels ...<-chan any) <-chan any {
	if len(channels) < 1 {
		var ch = make(chan any)
		close(ch)
		return ch
	}

	var cases []reflect.SelectCase
	for _, ch := range channels {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}

	for {
		if i, _, ok := reflect.Select(cases); !ok {
			return channels[i]
		}
	}
}
