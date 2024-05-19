package main

import (
	"bytes"
	"image/color"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kharism/hanashi/core"
)

type CardType int

const (
	CARD_TYPE_CHARACTER CardType = iota
	CARD_TYPE_ITEM
	CARD_TYPE_TRAP
	CARD_TYPE_OTHER
)

//go:embed assets/img/card_bg.png
var BASE_IMG []byte
var baseCard *ebiten.Image

func init() {
	imgReader := bytes.NewReader(BASE_IMG)
	baseCard, _, _ = ebitenutil.NewImageFromReader(imgReader)
}

// the card interface
type Card interface {
	Update() error
	Draw(screen *ebiten.Image)
	// get/build the image of the card. Probably need to update MovableImage lib of hanashi
	GetImage() *ebiten.Image
	GetType() CardType
	AddDecorator(CardDecorator)
	RemoveDecorator(CardDecorator)
	OnClick(*MainScene)
	GetDescription() string
}

// decorate card, add stuff on top of base card
type CardDecorator interface {
	Update() error
	// draw on the card
	Draw(card *ebiten.Image)
	GetType() CardType
	OnClick(mainScene *MainScene, source Card)
	GetDescription() string
}
type PlayerMoveListener interface {
	OnPlayerMove(Card, MainScene)
}

type BaseCard struct {
	*core.MovableImage
	decorators []CardDecorator
}

func (c *BaseCard) Update() error {
	c.MovableImage.SetImage(c.GetImage())
	for _, decor := range c.decorators {
		decor.Update()
	}
	c.MovableImage.Update()
	return nil
}

const (
	BASE_CARD_WIDTH  = 90
	BASE_CARD_HEIGHT = 140
	SCALE_CARD       = 0.8
	CARD_MOVE_SPEED  = 5
)

var (
	BaseColor = color.RGBA{
		R: 235,
		G: 203,
		B: 174,
		A: 255,
	}
	RED = color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	}
	BLUE = color.RGBA{
		R: 0,
		G: 0,
		B: 255,
		A: 255,
	}
)

func NewBaseCard(decorators []CardDecorator) Card {
	c := &BaseCard{}
	c.decorators = decorators
	param := core.MovableImageParams{}
	c.MovableImage = core.NewMovableImage(c.GetImage(), &param)
	return c
}

func (c *BaseCard) GetType() CardType {
	return c.decorators[0].GetType()
}
func (c *BaseCard) GetImage() *ebiten.Image {
	base := ebiten.NewImageFromImage(baseCard)
	// base := ebiten.NewImage(BASE_CARD_WIDTH, BASE_CARD_HEIGHT)
	// base.Fill(BaseColor)
	for _, decor := range c.decorators {
		decor.Draw(base)
	}
	return base
}

func (c *BaseCard) Draw(screen *ebiten.Image) {
	c.MovableImage.Draw(screen)
}

func (c *BaseCard) AddDecorator(decor CardDecorator) {
	c.decorators = append(c.decorators, decor)
}
func (c *BaseCard) RemoveDecorator(decor CardDecorator) {
	tempDecor := []CardDecorator{}
	for _, dec := range c.decorators {
		if dec == decor {
			continue
		}
		tempDecor = append(tempDecor, dec)
	}
	c.decorators = tempDecor
}
func (c *BaseCard) OnClick(state *MainScene) {
	if len(c.decorators) > 0 {
		c.decorators[0].OnClick(state, c)
	}
}
func (c *BaseCard) GetDescription() string {
	output := ""
	for _, v := range c.decorators {
		output += v.GetDescription() + "\n"
	}
	return output
}
