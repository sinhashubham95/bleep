package bleep

import (
	"os"
	"os/signal"
	"sync"

	"github.com/google/uuid"
)

var defaultBleep = New()

// Action is what will be performed on an os signal
type Action func(os.Signal)

// Bleep is the type for the handler
type Bleep struct {
	actions map[string]Action
	mu      *sync.RWMutex
}

// New is used to create a new os signal handler
func New() *Bleep {
	return &Bleep{
		actions: make(map[string]Action),
		mu:      &sync.RWMutex{},
	}
}

// Add is used to add an action to the default bleep instance to be performed on an os signal
func Add(action Action) string {
	return defaultBleep.Add(action)
}

// Remove is used to remove an action from the default bleep instance
func Remove(key string) Action {
	return defaultBleep.Remove(key)
}

// Reset is used to reset the default handler instance
func Reset() map[string]Action {
	return defaultBleep.Reset()
}

// Actions is used to get the set of actions part of the default bleep instance
func Actions() map[string]Action {
	return defaultBleep.Actions()
}

// Listen is used to listen for the provided OS signals
// If none are provided, then it will listen for any OS signal
func Listen() {
	defaultBleep.Listen()
}

// Add is used to add an action to be performed on an os signal
func (b *Bleep) Add(action Action) string {
	b.mu.Lock()
	defer b.mu.Unlock()
	key := uuid.New().String()
	b.actions[key] = action
	return key
}

// Remove is used to remove an action added previously
func (b *Bleep) Remove(key string) Action {
	b.mu.Lock()
	defer b.mu.Unlock()
	if action, ok := b.actions[key]; ok {
		delete(b.actions, key)
		return action
	}
	return nil
}

// Reset is used to reset the current handler instance
func (b *Bleep) Reset() map[string]Action {
	b.mu.Lock()
	defer b.mu.Unlock()
	actions := make(map[string]Action)
	for k, a := range b.actions {
		delete(b.actions, k)
		actions[k] = a
	}
	return actions
}

// Actions is used to get the set of actions part of the current bleep instance
func (b *Bleep) Actions() map[string]Action {
	b.mu.RLock()
	defer b.mu.RUnlock()
	actions := make(map[string]Action)
	for k, a := range b.actions {
		actions[k] = a
	}
	return actions
}

// Listen is used to listen for the provided OS signals
// If none are provided, then it will listen for any OS signal
func (b *Bleep) Listen(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	s := <-ch
	b.mu.RLock()
	defer b.mu.RUnlock()
	wg := sync.WaitGroup{}
	for _, a := range b.actions {
		wg.Add(1)
		action := a
		go func(s os.Signal) {
			defer wg.Done()
			action(s)
		}(s)
	}
	wg.Wait()
}
