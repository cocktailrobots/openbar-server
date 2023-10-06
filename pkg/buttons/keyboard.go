package buttons

import (
	"context"
	"fmt"
	"github.com/cocktailrobots/openbar-server/pkg/util"
	"golang.design/x/hotkey"
)

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
	down    *util.SafeBoolArray
}

func NewKeyboardButtons(ctx context.Context, numButtons int) (*KeyboardButtons, error) {
	if numButtons > len(keycodes) {
		numButtons = len(keycodes)
	}

	hotkeys := make([]*hotkey.Hotkey, numButtons)
	down := util.NewSafeBoolArray(numButtons)

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
