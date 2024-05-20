package main

import (
	"bytes"
	_ "embed"
	"fmt"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed assets/img/Knight.png
var KnightImage []byte

//go:embed assets/img/goblin.png
var GoblinImage []byte

//go:embed assets/fonts/PixelOperator8.ttf
var PixelFontTTF []byte

//go:embed assets/img/skeleton.png
var SkeletonImage []byte

var PixelFont *text.GoTextFaceSource

var knightImg *ebiten.Image
var goblinImg *ebiten.Image
var skeletonImg *ebiten.Image
var face *text.GoTextFace

type CharacterDecorator struct {
	Hp          int
	Name        string
	image       *ebiten.Image
	OnClickFunc OnInteractFunction
	OnDefeat    OnDefeatFunc
	Description string
}

type CharacterInterface interface {
	// take direct damage ignoring loadout
	TakeDirectDamage(int)
	// take damage but still putting loadout into consideration
	TakeDamage(int)
	Draw(card *ebiten.Image)

	DoBattle(*CharacterDecorator, *MainScene)
}
type OnInteractFunction func(*MainScene, Card)
type OnDefeatFunc func(*MainScene, Card)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(PixelFontTTF))
	if err != nil {
		log.Fatal(err)
	}
	PixelFont = s
	face = &text.GoTextFace{
		Source: PixelFont,
		Size:   24,
	}
	if knightImg == nil {
		imgReader := bytes.NewReader(KnightImage)
		knightImg, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if goblinImg == nil {
		imgReader := bytes.NewReader(GoblinImage)
		goblinImg, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if skeletonImg == nil {
		imgReader := bytes.NewReader(SkeletonImage)
		skeletonImg, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
}
func (d *CharacterDecorator) TakeDirectDamage(dmg int) {
	d.Hp -= dmg
	if d.Hp <= 1 {
		os.Exit(0)
	}
}
func (d *CharacterDecorator) TakeDamage(dmg int) {
	// the same with take damage
	d.Hp -= dmg
	if d.Hp <= 1 {
		os.Exit(0)
	}
}
func (d *CharacterDecorator) DoBattle(opp *CharacterDecorator, scene *MainScene) {
	d.TakeDamage(opp.Hp)
	opp.Hp = 0
}
func NewKnightDecor() CardDecorator {

	return &CharacterDecorator{Hp: 10, image: knightImg, Name: "Knight", OnClickFunc: func(s *MainScene, c Card) {

	}, Description: "Your Character"}
}

// return a function that handle
func GenerateCombat(damage int) func(*MainScene, Card) {
	return func(s *MainScene, source Card) {
		// posX, posY := source.(*BaseCard).GetPos()
		// idxX, idxY := PixelToIndex(int(posX), int(posY))
		// s.Character.Hp -= damage
		s.Character.DoBattle(source.(*BaseCard).decorators[0].(*CharacterDecorator), s)
		// jj := rwdGenerator.GenerateReward(0)
		if source.(*BaseCard).decorators[0].(*CharacterDecorator).Hp <= 0 {
			source.(*BaseCard).decorators[0].(*CharacterDecorator).OnDefeat(s, source)
		}

	}
}
func GenerateReward(tier int) func(*MainScene, Card) {
	return func(s *MainScene, source Card) {
		source.(*BaseCard).decorators[0] = rwdGenerator.GenerateReward(tier)
	}

}
func NewGoblinDecor() CardDecorator {

	return &CharacterDecorator{Hp: 3, image: goblinImg,
		Name: "Goblin", OnDefeat: GenerateReward(0), OnClickFunc: GenerateCombat(3), Description: "A small goblin"}
}
func NewSkeletonDecor() CardDecorator {
	if goblinImg == nil {
		imgReader := bytes.NewReader(SkeletonImage)
		skeletonImg, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	return &CharacterDecorator{
		Hp: 1, image: skeletonImg, Name: "Skeltn",
		OnClickFunc: GenerateCombat(1),
		OnDefeat:    GenerateReward(0),
		Description: "A small Skeleton",
	}
}
func (k *CharacterDecorator) Update() error {
	return nil
}
func (k *CharacterDecorator) GetType() CardType {
	return CARD_TYPE_CHARACTER
}

func (k *CharacterDecorator) OnClick(state *MainScene, source Card) {
	// return CARD_TYPE_CHARACTER
	k.OnClickFunc(state, source)
}
func (k *CharacterDecorator) GetDescription() string {
	return k.Description
}
func (k *CharacterDecorator) Draw(card *ebiten.Image) {
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(10, 50)
	txtOpt := text.DrawOptions{}
	txtOpt.GeoM.Scale(0.7, 0.7)
	txtOpt.GeoM.Translate(55, 6)

	txtOpt.ColorScale.ScaleWithColor(RED)
	text.Draw(card, fmt.Sprintf("%d", k.Hp), face, &txtOpt)
	// txtOpt.GeoM.Reset()
	txtOpt.GeoM.Translate(-30, 50)
	txtOpt.ColorScale.ScaleWithColor(RED)
	txtOpt.GeoM.Scale(0.6, 0.6)
	text.Draw(card, k.Name, face, &txtOpt)
	card.DrawImage(k.image, &opt)
}
