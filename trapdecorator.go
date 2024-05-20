package main

import (
	"bytes"
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kharism/hanashi/core"
)

//go:embed assets/img/spike.png
var Spike []byte
var spike *ebiten.Image

func init() {
	if spike == nil {
		imgReader := bytes.NewReader(Spike)
		spike, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}

}

type SpikeTrapDecorator struct {
	*CharacterDecorator
}

func (t *SpikeTrapDecorator) Update() error {
	return nil
}

func (t *SpikeTrapDecorator) GetType() CardType {
	return CARD_TYPE_TRAP
}
func (t *SpikeTrapDecorator) OnClick(state *MainScene, source Card) {
	// fmt.Println("Click Spike")
	posX, posY := source.(*BaseCard).GetPos()
	idxX, idxY := PixelToIndex(int(posX), int(posY))
	state.CharacterCard.AddAnimation(core.NewMoveAnimationFromParam(core.MoveParam{Tx: posX, Ty: posY, Speed: CARD_MOVE_SPEED}))
	// k.OnAccquire(state)
	state.Character.TakeDirectDamage(1)
	state.zones[idxY][idxX] = state.CharacterCard
	newCard := generator.GenerateCard(state)
	newCard.(*BaseCard).SetPos(float64(BORDER_X[PLAYER_IDX_X]), float64(BORDER_Y[PLAYER_IDX_Y]))
	state.zones[PLAYER_IDX_Y][PLAYER_IDX_X] = newCard.(*BaseCard)
	PLAYER_IDX_X = idxX
	PLAYER_IDX_Y = idxY
}

func NewSpikeTrapDecorator() CardDecorator {
	return &SpikeTrapDecorator{CharacterDecorator: &CharacterDecorator{image: spike, Name: "Spike", Hp: 1, Description: "Do 1 Direct damage\nregardless of loadout"}}
}
