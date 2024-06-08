package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type TransitionDecorator struct {
	Start           CardDecorator
	Finish          CardDecorator
	Card            *BaseCard
	translationDist int
}

func (d *TransitionDecorator) Update() error {
	if d.translationDist < BASE_CARD_WIDTH {
		d.translationDist += 10
	}
	if d.translationDist == BASE_CARD_WIDTH {
		d.Card.decorators[0] = d.Finish
	}
	return nil
}
func NewTransitionDecorator(start, finish CardDecorator, card *BaseCard) CardDecorator {
	return &TransitionDecorator{Start: start, Finish: finish, Card: card}
}
func (d *TransitionDecorator) Draw(card *ebiten.Image) {
	dummyCard1 := ebiten.NewImage(BASE_CARD_WIDTH, BASE_CARD_HEIGHT)
	dummyCard2 := ebiten.NewImage(BASE_CARD_WIDTH, BASE_CARD_HEIGHT)
	d.Start.Draw(dummyCard1)
	d.Finish.Draw(dummyCard2)
	finalImage := ebiten.NewImage(BASE_CARD_WIDTH*2, BASE_CARD_HEIGHT)
	opt1 := ebiten.DrawImageOptions{}
	opt1.GeoM.Translate(BASE_CARD_WIDTH, 0)
	finalImage.DrawImage(dummyCard1, &opt1)
	opt1.GeoM.Reset()
	finalImage.DrawImage(dummyCard2, &opt1)
	opt1.GeoM.Translate(float64(-BASE_CARD_WIDTH+d.translationDist), 0)
	card.DrawImage(finalImage, &opt1)
}
func (d *TransitionDecorator) GetType() CardType {
	return d.Start.GetType()
}
func (d *TransitionDecorator) OnClick(mainScene *MainScene, source Card) {
	fmt.Println("Transition Decorator clicked")
}
func (d *TransitionDecorator) TakeDirectDamage(dmg int, s *MainScene, source Card) {
	if v, ok := d.Start.(*CharacterDecorator); ok {
		v.TakeDirectDamage(dmg, s, source)
	}
}
func (d *TransitionDecorator) TakeDamage(dmg int, s *MainScene, source Card) {
	if v, ok := d.Start.(*CharacterDecorator); ok {
		v.TakeDamage(dmg, s, source)
	}
}
func (d *TransitionDecorator) DoBattle(c CharacterInterface, s *MainScene) {
	if v, ok := d.Start.(*CharacterDecorator); ok {
		v.DoBattle(c, s)
	}
}
func (d *TransitionDecorator) GetDescription() string {
	return d.Start.GetDescription()
}
