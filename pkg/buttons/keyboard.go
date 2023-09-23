//go:build !linux && !windows

package buttons

import (
	"context"
	"fmt"
	"golang.design/x/hotkey"
	"sync"
)

type safeBoolArray struct {
	mu   *sync.Mutex
	vals []bool
}

func NewSafeBoolArray(size int) *safeBoolArray {
	return &safeBoolArray{
		mu:   &sync.Mutex{},
		vals: make([]bool, size),
	}
}

func (s *safeBoolArray) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.vals {
		s.vals[i] = false
	}
}

func (s *safeBoolArray) Get(idx int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.vals[idx]
}

func (s *safeBoolArray) Set(idx int, val bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.vals[idx] = val
}

func (s *safeBoolArray) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.vals)
}

var keycodes = []hotkey.Key{
	hotkey.Key0,
	hotkey.Key1,
	hotkey.Key2,
	hotkey.Key3,
	hotkey.Key4,
	hotkey.Key5,
	hotkey.Key6,
	hotkey.Key7,
	hotkey.Key8,
	hotkey.Key9,
	hotkey.KeyA,
	hotkey.KeyB,
	hotkey.KeyC,
	hotkey.KeyD,
	hotkey.KeyE,
	hotkey.KeyF,
	hotkey.KeyG,
	hotkey.KeyH,
	hotkey.KeyI,
	hotkey.KeyJ,
	hotkey.KeyK,
	hotkey.KeyL,
	hotkey.KeyM,
	hotkey.KeyN,
	hotkey.KeyO,
	hotkey.KeyP,
	hotkey.KeyQ,
	hotkey.KeyR,
	hotkey.KeyS,
	hotkey.KeyT,
	hotkey.KeyU,
	hotkey.KeyV,
	hotkey.KeyW,
	hotkey.KeyX,
	hotkey.KeyY,
	hotkey.KeyZ,
}

type KeyboardButtons struct {
	hotkeys []*hotkey.Hotkey
	down    *safeBoolArray
}

func NewKeyboardButtons(ctx context.Context, numButtons int) (*KeyboardButtons, error) {
	if numButtons > len(keycodes) {
		numButtons = len(keycodes)
	}

	hotkeys := make([]*hotkey.Hotkey, numButtons)
	down := NewSafeBoolArray(numButtons)

	for i := 0; i < numButtons; i++ {
		hotkeys[i] = hotkey.New(nil, keycodes[i])

		err := hotkeys[i].Register()
		if err != nil {
			return nil, fmt.Errorf("hotkey: failed to register hotkey: %v", err)
		}
	}

	return &KeyboardButtons{
		hotkeys: hotkeys,
		down:    down,
	}, nil
}

func (k *KeyboardButtons) NumButtons() int {
	return k.down.Len()
}

func (k *KeyboardButtons) IsPressed(idx int) bool {
	return k.down.Get(idx)
}

func (k *KeyboardButtons) Update() error {
	for i, hk := range k.hotkeys {
		select {
		case <-hk.Keydown():
			k.down.Set(i, true)
		case <-hk.Keyup():
			k.down.Set(i, false)
		default:
		}
	}
	return nil
}

func (k *KeyboardButtons) Close() error {
	var err error
	for _, hk := range k.hotkeys {
		unrErr := hk.Unregister()
		if unrErr != nil && err == nil {
			err = unrErr
		}
	}

	return err
}
