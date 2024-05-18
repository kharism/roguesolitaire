package main

import (
	"bytes"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	_ "embed"
)

//go:embed assets/img/chest_anim.png
var Chest []byte
var chestImg *ebiten.Image

type ChestDecorator struct {
	Name         string
	image        *ebiten.Image
	isAnim       bool
	sourceCard   Card
	counter      int
	dropResource func()
}

func init() {
	if chestImg == nil {
		imgReader := bytes.NewReader(Chest)
		chestImg, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
}
func NewChestDecorator() CardDecorator {
	rect := image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{64, 64}}
	firstFrame := chestImg.SubImage(rect)
	return &ChestDecorator{image: firstFrame.(*ebiten.Image)}
}

func (c *ChestDecorator) Update() error {
	if c.isAnim {
		c.counter += 1
	}
	if c.counter/5 == 2 {
		c.dropResource()
	}
	return nil
}
func (c *ChestDecorator) Draw(card *ebiten.Image) {
	if c.isAnim {
		if c.counter/5 == 1 {
			rect := image.Rectangle{Min: image.Point{64, 0}, Max: image.Point{128, 64}}
			newFrame := chestImg.SubImage(rect)
			c.image = newFrame.(*ebiten.Image)
		} else if c.counter/5 == 2 {
			rect := image.Rectangle{Min: image.Point{128, 0}, Max: image.Point{128 + 64, 64}}
			newFrame := chestImg.SubImage(rect)
			c.image = newFrame.(*ebiten.Image)
		}
	}
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(10, 50)
	txtOpt := text.DrawOptions{}
	// txtOpt.GeoM.Translate(50, 0)
	// txtOpt.ColorScale.ScaleWithColor(RED)
	// text.Draw(card, fmt.Sprintf("%d", k.Hp), face, &txtOpt)
	// txtOpt.GeoM.Reset()
	txtOpt.GeoM.Translate(20, 50)
	txtOpt.ColorScale.ScaleWithColor(RED)
	txtOpt.GeoM.Scale(0.6, 0.6)
	text.Draw(card, c.Name, face, &txtOpt)
	card.DrawImage(c.image, &opt)
}
func (c *ChestDecorator) GetType() CardType {
	return CARD_TYPE_ITEM
}
func (c *ChestDecorator) OnClick(mainScene *MainScene, source Card) {
	fmt.Println("Chest clicked")
	c.isAnim = true
	c.dropResource = func() {
		source.RemoveDecorator(c)
		source.AddDecorator(NewLightPotionDecorator())
	}
}
