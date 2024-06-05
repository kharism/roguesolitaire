package main

import (
	"bytes"
	"fmt"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// To make long story short
// SwordDecorator is the decorator used to decorate BaseCard
// SwordChDecorator is the decorator used to decorate BaseCard which already has characterdecorator

type SwordDecorator struct {
	*ItemDecorator
}
type SwordChDecorator struct {
	// *CharacterDecorator
	chInterface CharacterInterface
	durability  int
}

//go:embed assets/img/sword.png
var sword []byte
var swordImg *ebiten.Image

func init() {
	imgReader := bytes.NewReader(sword)
	swordImg, _, _ = ebitenutil.NewImageFromReader(imgReader)
}

func NewSwordDecorator() CardDecorator {
	imgrotated := ebiten.NewImage(64, 64)
	// imgrotated.Fill(color.White)
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Rotate(0.8)
	opts.GeoM.Translate(37, 0)
	opts.GeoM.Scale(1.5, 1.5)
	imgrotated.DrawImage(swordImg, &opts)
	j := &ItemDecorator{image: imgrotated, Name: "Sword", OnAccquire: func(m *MainScene) {
		if _, ok := m.Character.(*CharacterDecorator); ok {
			m.Character = NewSwordChDecorator(m.Character, 5).(*SwordChDecorator)
			m.zones[PLAYER_IDX_Y][PLAYER_IDX_X].decorators[0] = m.Character.(*SwordChDecorator)
		} else if vv, ok := m.Character.(*SwordChDecorator); ok {
			vv.durability += 5
			m.Character = vv
		}

	}, Description: "Gain 10 combat"}
	return &SwordDecorator{ItemDecorator: j}
}
func (d *SwordChDecorator) GetHP() int {
	return d.chInterface.GetHP()
}
func (d *SwordChDecorator) GetMaxHP() int {
	return d.chInterface.GetMaxHP()
}
func (d *SwordDecorator) Draw(card *ebiten.Image) {
	d.ItemDecorator.Draw(card)

}
func NewSwordChDecorator(ch CharacterInterface, durability int) CardDecorator {
	h := &SwordChDecorator{chInterface: ch, durability: durability}
	return h
}

func (h *SwordChDecorator) Draw(card *ebiten.Image) {

	h.chInterface.Draw(card)
	imgOpt := ebiten.DrawImageOptions{}
	imgOpt.GeoM.Scale(1.1, 1.5)
	imgOpt.GeoM.Translate(0, 25)
	card.DrawImage(swordImg, &imgOpt)
	txtOpt := text.DrawOptions{}
	txtOpt.GeoM.Scale(0.7, 0.7)
	txtOpt.GeoM.Translate(6, 6)

	txtOpt.ColorScale.ScaleWithColor(BLUE)
	text.Draw(card, fmt.Sprintf("%d", h.durability), face, &txtOpt)
}
func (c *SwordChDecorator) DoBattle(opp *CharacterDecorator, scene *MainScene) {
	if opp.Hp <= c.durability {
		c.durability -= opp.Hp
		opp.Hp = 0
	} else {
		opp.Hp -= c.durability
		c.durability = 0
	}
	if c.durability == 0 {
		scene.Character = c.chInterface
		scene.zones[PLAYER_IDX_Y][PLAYER_IDX_X].decorators[0] = scene.Character.(*CharacterDecorator)
	}

}
func (c *SwordChDecorator) GetType() CardType {
	return CARD_TYPE_CHARACTER
}
func (c *SwordChDecorator) GetDescription() string {
	return ""
}
func (c *SwordChDecorator) Update() error {
	return nil
}
func (c *SwordChDecorator) OnClick(s *MainScene, card Card) {

}
func (c *SwordChDecorator) TakeDamage(dmg int, s *MainScene, card Card) {

	// c.chInterface.TakeDamage()
}
func (c *SwordChDecorator) TakeDirectDamage(dmg int, s *MainScene, card Card) {
	c.chInterface.TakeDirectDamage(dmg, s, card)
}
