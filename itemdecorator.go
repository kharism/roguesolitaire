package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	_ "image/png"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kharism/hanashi/core"
)

//go:embed assets/img/coin.png
var CoinImage []byte
var coinImg *ebiten.Image

//go:embed assets/img/Potion1.png
var Potion1 []byte
var potion1 *ebiten.Image

//go:embed assets/img/Potion2.png
var Potion2 []byte
var potion2 *ebiten.Image

//go:embed assets/img/meat.png
var Meat []byte
var meat *ebiten.Image

const meat_img_length = 3

type OnAccquireFunc func(*MainScene)
type ItemDecorator struct {
	Name        string
	image       *ebiten.Image
	OnAccquire  OnAccquireFunc
	Description string
}

func GenerateAddCoinFunc(Amount int) OnAccquireFunc {
	return func(state *MainScene) {
		fmt.Println("Get coin")
		state.State.Coin += Amount
	}
}
func (k *ItemDecorator) OnClick(state *MainScene, source Card) {
	posX, posY := source.(*BaseCard).GetPos()
	idxX, idxY := PixelToIndex(int(posX), int(posY))
	state.CharacterCard.AddAnimation(core.NewMoveAnimationFromParam(core.MoveParam{Tx: posX, Ty: posY, Speed: CARD_MOVE_SPEED}))
	k.OnAccquire(state)
	state.zones[idxY][idxX] = state.CharacterCard
	newCard := generator.GenerateCard()
	newCard.(*BaseCard).SetPos(float64(BORDER_X[PLAYER_IDX_X]), float64(BORDER_Y[PLAYER_IDX_Y]))
	state.zones[PLAYER_IDX_Y][PLAYER_IDX_X] = newCard.(*BaseCard)
	PLAYER_IDX_X = idxX
	PLAYER_IDX_Y = idxY
}
func (k *ItemDecorator) Update() error {
	return nil
}
func (k *ItemDecorator) GetType() CardType {
	return CARD_TYPE_ITEM
}
func init() {
	if coinImg == nil {
		imgReader := bytes.NewReader(CoinImage)
		coinImg, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if potion1 == nil {
		imgReader := bytes.NewReader(Potion1)
		potion1, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if potion2 == nil {
		imgReader := bytes.NewReader(Potion2)
		potion2, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	imgReader := bytes.NewReader(Meat)
	meat, _, _ = ebitenutil.NewImageFromReader(imgReader)
}
func NewCoinDecorator() CardDecorator {
	// if coinImg == nil {
	// 	imgReader := bytes.NewReader(CoinImage)
	// 	coinImg, _, _ = ebitenutil.NewImageFromReader(imgReader)
	// }
	return &ItemDecorator{image: coinImg, Name: "Coin", OnAccquire: GenerateAddCoinFunc(1), Description: "Add 1 gold coin"}
}
func (c *ItemDecorator) GetDescription() string {
	return c.Description
}
func NewMedPotionDecorator() CardDecorator {
	// if potion1 == nil {
	// 	imgReader := bytes.NewReader(Potion1)
	// 	potion1, _, _ = ebitenutil.NewImageFromReader(imgReader)
	// }
	return &ItemDecorator{image: potion1, Name: "Red\nPotion(m)", OnAccquire: func(s *MainScene) {
		// s.Character.TakeDamage(-3)
		s.Character.TakeDirectDamage(-4)
	}, Description: "Recover 4 HP"}
}
func NewLightPotionDecorator() CardDecorator {
	// if potion1 == nil {
	// 	imgReader := bytes.NewReader(Potion1)
	// 	potion1, _, _ = ebitenutil.NewImageFromReader(imgReader)
	// }
	return &ItemDecorator{image: potion1, Name: "Red\nPotion(s)", OnAccquire: func(s *MainScene) {
		s.Character.TakeDirectDamage(-2)
	}, Description: "recover 2 HP"}
}
func NewMeat() CardDecorator {
	idx := rand.Int() % meat_img_length

	SliceIdxImg := meat.SubImage(image.Rect(64*idx, 0, 64*(idx+1), 64))
	return &ItemDecorator{image: SliceIdxImg.(*ebiten.Image), Name: "Meat", OnAccquire: func(s *MainScene) {
		s.Character.TakeDirectDamage(-1)
	}, Description: "Recover 1 HP"}
}
func (k *ItemDecorator) Draw(card *ebiten.Image) {
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(10, 50)
	txtOpt := text.DrawOptions{}
	// txtOpt.GeoM.Translate(50, 0)
	// txtOpt.ColorScale.ScaleWithColor(RED)
	// text.Draw(card, fmt.Sprintf("%d", k.Hp), face, &txtOpt)
	// txtOpt.GeoM.Reset()
	txtOpt.GeoM.Translate(90, 50)
	txtOpt.ColorScale.ScaleWithColor(RED)
	// txtOpt.LineSpacing = 10
	txtOpt.LayoutOptions = text.LayoutOptions{PrimaryAlign: text.AlignCenter, LineSpacing: 20}
	txtOpt.GeoM.Scale(0.5, 0.5)
	text.Draw(card, k.Name, face, &txtOpt)
	card.DrawImage(k.image, &opt)
}
