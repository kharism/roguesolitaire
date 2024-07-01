package main

import (
	"bytes"
	"fmt"
	"math"
	"strings"

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

var (
	//direction array containing weakness to one direction
	DIRECTION_ARR_ONE_WEAKNESS = []byte{DIRECTION_UP, DIRECTION_RIGHT, DIRECTION_DOWN, DIRECTION_LEFT}
	//direction array containing weakness to two direction, both direction are opposing to one another
	DIRECTION_ARR_TWO_WEAKNESS = []byte{DIRECTION_UP | DIRECTION_DOWN, DIRECTION_LEFT | DIRECTION_RIGHT}
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
	WeaknessList := []string{}
	dist := []byte{1, 2, 4, 8}
	direction := []string{"UP", "RIGHT", "DOWN", "LEFT"}
	for idx, v := range dist {
		if d.Direction&v != 0 {
			WeaknessList = append(WeaknessList, direction[idx])
		}
	}
	allWeakness := strings.Join(WeaknessList, "\n-")
	if len(allWeakness) == 0 {
		return d.decorator.GetDescription()
	}
	return d.decorator.GetDescription() + "\nWeakness:\n-" + allWeakness
}
func (d *WeaknessDecorator) TakeDirectDamage(dmg int, s *MainScene, card Card) {
	if c, ok := d.decorator.(CharacterInterface); ok {
		c.TakeDirectDamage(dmg, s, card)
	}
}

// take damage but still putting loadout into consideration
func (d *WeaknessDecorator) TakeDamage(dmg int, s *MainScene, card Card) {
	if c, ok := d.decorator.(CharacterInterface); ok {
		c.TakeDamage(dmg, s, card)
	}
}
func (d *WeaknessDecorator) GetHP() int {
	return d.decorator.(CharacterInterface).GetHP()
}
func (d *WeaknessDecorator) SetHP(hp int) {
	d.decorator.(CharacterInterface).SetHP(hp)
}
func (d *WeaknessDecorator) GetMaxHP() int {
	return d.decorator.(CharacterInterface).GetMaxHP()
}
func (d *WeaknessDecorator) SetMaxHP(hp int) {
	d.decorator.(CharacterInterface).SetMaxHP(hp)
}
func (d *WeaknessDecorator) GetOnDefeat() OnDefeatFunc {
	return d.decorator.(CharacterInterface).GetOnDefeat()
}

func (d *WeaknessDecorator) DoBattle(CharacterInterface, *MainScene) {}

type RotatingWeaknessDecorator struct {
	*WeaknessDecorator
}

// only works for 4 byte number
func rotateByte(b byte) byte {
	h := b << 1
	excess := h & 0b11110000
	excess = excess >> 4
	return (h | excess) & 0b00001111
}
func NewRotatingWeaknessDecorator(decorator CardDecorator, direction byte) CardDecorator {
	weaknessDecor := NewWeaknessDecorator(decorator, direction)
	return &RotatingWeaknessDecorator{WeaknessDecorator: weaknessDecor.(*WeaknessDecorator)}
}
func (d *RotatingWeaknessDecorator) TakeDirectDamage(dmg int, s *MainScene, card Card) {
	if c, ok := d.decorator.(CharacterInterface); ok {
		c.TakeDirectDamage(dmg, s, card)
	}
}
func (d *RotatingWeaknessDecorator) Draw(card *ebiten.Image) {
	d.WeaknessDecorator.Draw(card)
}

// take damage but still putting loadout into consideration
func (d *RotatingWeaknessDecorator) TakeDamage(dmg int, s *MainScene, card Card) {
	if c, ok := d.decorator.(CharacterInterface); ok {
		c.TakeDamage(dmg, s, card)
	}
}
func (d *RotatingWeaknessDecorator) OnClick(mainScene *MainScene, source Card) {
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
func (d *RotatingWeaknessDecorator) DoBattle(*CharacterDecorator, *MainScene) {}
func (t *RotatingWeaknessDecorator) OnPlayerMove(c Card, s *MainScene) {
	curDirection := t.WeaknessDecorator.Direction
	// bits.RotateLeft(uint(curDirection), 1)
	curDirection = rotateByte(curDirection)
	// newDecor := NewWeaknessDecorator(t.WeaknessDecorator.decorator, curDirection)
	rota := NewRotatingWeaknessDecorator(t.WeaknessDecorator.decorator, curDirection)
	fmt.Printf("%b\n", curDirection)
	newTrans := NewTransitionDecorator(t.decorator, rota, c.(*BaseCard))
	c.(*BaseCard).decorators[0] = newTrans
	// c.RemoveDecorator(t)
	// c.AddDecorator(newWeakness)
}
