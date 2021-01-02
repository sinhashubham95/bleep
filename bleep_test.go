package bleep_test

import (
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/sinhashubham95/bleep"
)

func Test(t *testing.T) {
	defer bleep.Reset()

	mu := sync.Mutex{}
	testData := make(map[int]int)

	// add some actions
	bleep.Add(func(s os.Signal) {
		mu.Lock()
		defer mu.Unlock()
		if s == syscall.SIGTERM {
			testData[1]++
		}
	})
	bleep.Add(func(s os.Signal) {
		mu.Lock()
		defer mu.Unlock()
		if s == syscall.SIGTERM {
			testData[2]++
		}
	})
	bleep.Add(func(s os.Signal) {
		mu.Lock()
		defer mu.Unlock()
		if s == syscall.SIGABRT {
			testData[3]++
		}
	})
	k4 := bleep.Add(func(s os.Signal) {
		mu.Lock()
		defer mu.Unlock()
		if s == syscall.SIGTERM {
			testData[4]++
		}
	})

	// remove an action
	a4 := bleep.Remove(k4)

	// remove a non existing action
	a := bleep.Remove("sample")
	if a != nil {
		t.Errorf("Action should not be existing but it does.")
	}

	actions := bleep.Actions()
	if len(actions) != 3 {
		t.Errorf("Invalid number of actions expected %d found %d.", 3, len(actions))
	}

	if len(testData) != 0 {
		t.Errorf("Test data is not empty %+v.", testData)
	}

	// start a go routine that will produce a OS signal
	go func() {
		// wait for 100ms
		time.Sleep(100 * time.Millisecond)
		// send a signal
		p, err := os.FindProcess(os.Getpid())
		if err != nil {
			panic(err.Error())
		}
		p.Signal(syscall.SIGTERM)
	}()

	// now listen to the signal
	bleep.Listen()

	if len(testData) != 2 {
		t.Errorf("Test data invalid expected length %d found %d.", 2, len(testData))
	}

	if testData[1] != 1 || testData[2] != 1 || testData[3] != 0 || testData[4] != 0 {
		t.Errorf("Invalid test data %+v.", testData)
	}

	a4(syscall.SIGTERM)
	if testData[4] != 1 {
		t.Errorf("Invalid removed action.")
	}
}
