package main

import (
	"bytes"
	"fmt"
	"math"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// weakness will give the card some weakness the player can exploit
type WeaknessDecorator struct {
	decorator CardDecorator
	Direction byte
}

const (
	// use bitwise operator to get combination of these
	DIRECTION_UP    = 0b0001
	DIRECTION_RIGHT = 0b0010
	DIRECTION_DOWN  = 0b0100
	DIRECTION_LEFT  = 0b1000
)

func init() {
	reader := bytes.NewReader(arrow)
	arrowImg, _, _ = ebitenutil.NewImageFromReader(reader)
}

func NewWeaknessDecorator(decorator CardDecorator, direction byte) CardDecorator {
	return &WeaknessDecorator{decorator: decorator, Direction: direction}
}

func (d *WeaknessDecorator) Update() error {
	return nil
}

//go:embed assets/img/arrow_weakness.png
var arrow []byte
var arrowImg *ebiten.Image

// draw on the card
func (d *WeaknessDecorator) Draw(card *ebiten.Image) {
	d.decorator.Draw(card)
	rect := card.Bounds()

	if d.Direction&DIRECTION_DOWN != 0 {
		opt := ebiten.DrawImageOptions{}
		opt.GeoM.Rotate(math.Pi)
		opt.GeoM.Translate(55, float64(rect.Max.Y-10))
		card.DrawImage(arrowImg, &opt)
	}
	if d.Direction&DIRECTION_UP != 0 {
		opt := ebiten.DrawImageOptions{}
		opt.GeoM.Translate(35, 0)
		card.DrawImage(arrowImg, &opt)
	}
	if d.Direction&DIRECTION_LEFT != 0 {
		opt := ebiten.DrawImageOptions{}
		opt.GeoM.Rotate(-math.Pi / 2)
		opt.GeoM.Translate(5, 80)
		card.DrawImage(arrowImg, &opt)
	}
	if d.Direction&DIRECTION_RIGHT != 0 {
		opt := ebiten.DrawImageOptions{}
		opt.GeoM.Rotate(math.Pi / 2)
		opt.GeoM.Translate(80, 60)
		card.DrawImage(arrowImg, &opt)
	}

}
func (d *WeaknessDecorator) GetType() CardType {
	return d.decorator.GetType()
}
func (d *WeaknessDecorator) OnClick(mainScene *MainScene, source Card) {
	posX, posY := source.(*BaseCard).GetPos()
	idxX, idxY := PixelToIndex(int(posX), int(posY))
	//player come from above
	if PLAYER_IDX_X == idxX && PLAYER_IDX_Y < idxY && int(d.Direction&DIRECTION_UP) > 0 {
		if opp, ok := d.decorator.(*CharacterDecorator); ok {
			opp.OnDefeat(mainScene, source)
		}
	} else if PLAYER_IDX_X == idxX && PLAYER_IDX_Y > idxY && int(d.Direction&DIRECTION_DOWN) > 0 {
		//player come from below
		if opp, ok := d.decorator.(*CharacterDecorator); ok {
			opp.OnDefeat(mainScene, source)
		}
	} else if PLAYER_IDX_Y == idxY && PLAYER_IDX_X < idxX && int(d.Direction&DIRECTION_LEFT) > 0 {
		//player come from left
		if opp, ok := d.decorator.(*CharacterDecorator); ok {
			opp.OnDefeat(mainScene, source)
		}
	} else if PLAYER_IDX_Y == idxY && PLAYER_IDX_X > idxX && int(d.Direction&DIRECTION_RIGHT) > 0 {
		//player come from right
		if opp, ok := d.decorator.(*CharacterDecorator); ok {
			opp.OnDefeat(mainScene, source)
		}
	} else {
		source.(*BaseCard).decorators[0] = d.decorator
		d.decorator.OnClick(mainScene, source)
		if d.decorator == source.(*BaseCard).decorators[0] {
			source.(*BaseCard).decorators[0] = d
		}
	}

}
func (d *WeaknessDecorator) GetDescription() string {
	return d.decorator.GetDescription() + fmt.Sprintf("\n %b", d.Direction)
}
