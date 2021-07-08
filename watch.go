package main

import (
	"time"
)

func Watch(dir string, changes chan Change) {
	oldState, err := NewSnapshot(dir)
	ExitIf(err)

	for {
		newState, err := NewSnapshot(dir)
		ExitIf(err)
		diff, err := newState.Diff(oldState)
		ExitIf(err)

		for _, change := range diff {
			changes <- change
		}

		oldState = newState
		time.Sleep(100 * time.Millisecond)
	}
}
