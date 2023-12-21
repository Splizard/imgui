package main

import (
	"fmt"

	G "github.com/AllenDang/giu"
	I "github.com/AllenDang/imgui-go"
)

//simple widget that allows entering 2 hexadecimal characters

type InputHexByte struct {
	id        string
	hi, lo    byte //should be <16
	at        int  //0 or 1 (hi or low nibble)
	cbCancel  func()
	cbSuccess func(b byte)
}

func (ih *InputHexByte) Dispose() {
	//empty
}

func InputHex(id string, cancel func(), success func(b byte)) G.Widget {
	hRaw := G.Context.GetState(id)
	var h *InputHexByte
	if hRaw != nil {
		h = hRaw.(*InputHexByte)
		h.cbCancel = cancel
		h.cbSuccess = success
	} else {
		h = &InputHexByte{
			cbCancel:  cancel,
			cbSuccess: success,
		}
	}
	G.Context.SetState(id, h)
	return h
}

//value input keys only
var inputKeys = map[G.Key]byte{
	G.Key0: 0,
	G.Key1: 1,
	G.Key2: 2,
	G.Key3: 3,
	G.Key4: 4,
	G.Key5: 5,
	G.Key6: 6,
	G.Key7: 7,
	G.Key8: 8,
	G.Key9: 9,
	G.KeyA: 0xA,
	G.KeyB: 0xB,
	G.KeyC: 0xC,
	G.KeyD: 0xD,
	G.KeyE: 0xE,
	G.KeyF: 0xF,
}

func (ih *InputHexByte) Build() {
	var txt string
	switch ih.at {
	case 0:
		txt = "__ "
	case 1:
		txt = fmt.Sprintf("%1X_ ", ih.hi)
	case 2:
		txt = fmt.Sprintf("%1X%1X ", ih.hi, ih.lo)
	}

	I.Text(txt)

	//handle input keys
	if G.IsKeyPressed(G.KeyBackspace) {
		ih.at--
		if ih.at < 0 {
			ih.cancel()
			return
		}
	}

	if G.IsKeyPressed(G.KeyEscape) {
		ih.cancel()
		return
	}

	var input byte = 255
	for k, v := range inputKeys {
		if G.IsKeyPressed(k) {
			input = v
		}
	}

	if input != 255 {
		if ih.at == 0 {
			ih.hi = input
		} else {
			ih.lo = input
		}
		ih.at++
		if ih.at >= 2 {
			ih.success()
		}
	}

	G.Context.SetState(ih.id, ih)
}

func (ih *InputHexByte) cancel() {
	if ih.cbCancel != nil {
		ih.cbCancel()
	}
	ih.reset()
}

func (ih *InputHexByte) reset() {
	ih.hi = 0
	ih.lo = 0
	ih.at = 0
}

func (ih *InputHexByte) success() {
	if ih.cbSuccess != nil {
		ih.cbSuccess(ih.hi<<4 | ih.lo)
		ih.reset()
	}
}
