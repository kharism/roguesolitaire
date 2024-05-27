package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	// "github.com/hajimehoshi/ebiten/v2/internal/ui"
)

var (
	TouchIDs []ebiten.TouchID
	TouchPos map[ebiten.TouchID]*Pos
)

func init() {
	TouchIDs = []ebiten.TouchID{}
	TouchPos = map[ebiten.TouchID]*Pos{}
}

// return whether a click or tap is happened, and its location if it happened
func IsClickedOrTap() (bool, int, int) {
	posX := -1
	posY := -1
	mouseReleased := inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0)
	if mouseReleased {
		posX, posY = ebiten.CursorPosition()
		return true, posX, posY
	}
	// fmt.Println("Check", TouchIDs)
	for _, id := range TouchIDs {
		// fmt.Println("Check", id, inpututil.TouchPressDuration(id))
		if inpututil.IsTouchJustReleased(id) {

			posX, posY = int(TouchPos[id].X), int(TouchPos[id].Y)
			fmt.Println("touch released", id, posX, posY)

			return true, posX, posY
		}
	}
	TouchIDs = inpututil.AppendJustPressedTouchIDs(TouchIDs[:0])
	for _, id := range TouchIDs {
		x, y := ebiten.TouchPosition(id)
		TouchPos[id] = &Pos{
			X: float64(x),
			Y: float64(y),
		}
	}
	TouchIDs = ebiten.AppendTouchIDs(TouchIDs[:0])
	// fmt.Println(TouchIDs)
	return false, posX, posY
}
