package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// this decorator will inflict direct damage to nearby cards each time a player
// move.
type AttackOnPlayerMoveDecorator struct {
	CharacterDecorator *CharacterDecorator
	Damage             int
}

func (b *AttackOnPlayerMoveDecorator) TakeDirectDamage(dmg int, scene *MainScene, c Card) {
	b.CharacterDecorator.TakeDirectDamage(dmg, scene, c)
}
func (b *AttackOnPlayerMoveDecorator) TakeDamage(dmg int, scene *MainScene, c Card) {
	b.CharacterDecorator.TakeDamage(dmg, scene, c)
}
func (b *AttackOnPlayerMoveDecorator) Draw(card *ebiten.Image) {
	b.CharacterDecorator.Draw(card)
}
func (b *AttackOnPlayerMoveDecorator) GetHP() int {
	return b.CharacterDecorator.Hp
}
func (b *AttackOnPlayerMoveDecorator) GetMaxHP() int {
	return b.CharacterDecorator.MaxHP
}
func (d *AttackOnPlayerMoveDecorator) GetType() CardType {
	return d.CharacterDecorator.GetType()
}
func (b *AttackOnPlayerMoveDecorator) DoBattle(c CharacterInterface, scene *MainScene) {

}
func (b *AttackOnPlayerMoveDecorator) GetOnDefeat() OnDefeatFunc {
	return b.CharacterDecorator.GetOnDefeat()
}
func (b *AttackOnPlayerMoveDecorator) SetHP(a int) {
	b.CharacterDecorator.SetHP(a)
}
func (b *AttackOnPlayerMoveDecorator) SetMaxHP(a int) {
	b.CharacterDecorator.SetMaxHP(a)
}
func (b *AttackOnPlayerMoveDecorator) Update() error {
	return nil
}
func (b *AttackOnPlayerMoveDecorator) OnClick(mainScene *MainScene, source Card) {
	b.CharacterDecorator.OnClick(mainScene, source)
}
func (b *AttackOnPlayerMoveDecorator) GetDescription() string {
	return b.CharacterDecorator.GetDescription()
}
func (b *AttackOnPlayerMoveDecorator) OnPlayerMove(card Card, s *MainScene) {
	cardX, cardY := GetPos(s.zones, card.(*BaseCard))
	posX, posY := IdxToPixel(cardX, cardY)
	idxX, idxY := PixelToIndex(int(posX), int(posY))
	adjacent := []*BaseCard{}
	for i := idxY - 1; i <= idxY+1; i++ {
		if i < 0 || i > 2 {
			continue
		}
		for j := idxX - 1; j <= idxX+1; j++ {
			if j < 0 || j > 2 {
				continue
			}
			if i == idxY && j == idxX {
				continue
			}
			if math.Abs(float64(idxY-i))+math.Abs(float64(idxX-j)) > 1 {
				continue
			}
			adjacent = append(adjacent, s.zones[i][j])
			if i == PLAYER_IDX_Y && j == PLAYER_IDX_X {
				fmt.Println("SSDDS")
			}
			// if s.zones[i][j].decorators[0].(CharacterInterface) == s.Character {
			// 	fmt.Println("DSDSDSD")
			// }
			if v, ok := s.zones[i][j].decorators[0].(CharacterInterface); ok {
				fmt.Println(s.zones[i][j].decorators[0].GetDescription())
				v.TakeDirectDamage(b.Damage, s, s.zones[i][j])
			}
		}
	}
}
