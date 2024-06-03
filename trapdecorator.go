package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kharism/hanashi/core"
)

//go:embed assets/img/spike.png
var Spike []byte
var spike *ebiten.Image

//go:embed assets/img/bomb.png
var Bomb []byte
var bomb *ebiten.Image

func init() {
	if spike == nil {
		imgReader := bytes.NewReader(Spike)
		spike, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if bomb == nil {
		imgReader := bytes.NewReader(Bomb)
		bomb, _, _ = ebitenutil.NewImageFromReader(imgReader)
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
	state.Character.TakeDirectDamage(t.Hp, state, source)
	state.zones[idxY][idxX] = state.CharacterCard
	movedCard, newCardIdxY, newCardIdxX := GetMovedCard(state, idxX, idxY)
	oldCardPosX, oldCardPosY := movedCard.GetPos()
	newMoveParam := core.MoveParam{Tx: float64(BORDER_X[PLAYER_IDX_X]), Ty: float64(BORDER_Y[PLAYER_IDX_Y]), Speed: CARD_MOVE_SPEED}
	movedCard.AddAnimation(core.NewMoveAnimationFromParam(newMoveParam))
	state.CurMovingCard = movedCard
	state.zones[PLAYER_IDX_Y][PLAYER_IDX_X] = movedCard
	newCard := generator.GenerateCard(state)
	newCard.(*BaseCard).SetPos(oldCardPosX, oldCardPosY)
	state.zones[newCardIdxY][newCardIdxX] = newCard.(*BaseCard)

	PLAYER_IDX_X = idxX
	PLAYER_IDX_Y = idxY
	state.OnPlayerMove()
}
func NewCrimsonTrapDecorator() CardDecorator {
	Hp := []int{10, 15, 20}
	return &SpikeTrapDecorator{CharacterDecorator: &CharacterDecorator{image: spike,
		shader:      dakkaShader,
		Name:        "Crimson Spike",
		Hp:          Hp[rand.Int()%len(Hp)],
		OnDefeat:    GenerateReward(1),
		Description: "Do direct damage\nequal to its cost\nAn explosion may\nneutralize this"}}
}
func NewSpikeTrapDecorator() CardDecorator {
	return &SpikeTrapDecorator{CharacterDecorator: &CharacterDecorator{image: spike,
		Name:        "Spike",
		Hp:          1,
		OnDefeat:    GenerateReward(1),
		Description: "Do 1 Direct damage\nregardless of loadout\nAn explosion may\nneutralize this"}}
}

type BombDecorator struct {
	TurnToExplode int
	img           *ebiten.Image
}

func NewBombDecorator() CardDecorator {
	return &BombDecorator{TurnToExplode: 3, img: bomb}
}
func (b *BombDecorator) Update() error {
	return nil
}

// draw on the card
func (b *BombDecorator) Draw(card *ebiten.Image) {
	var curSprite *ebiten.Image
	var rect image.Rectangle
	switch b.TurnToExplode {
	case 3:
		rect.Min = image.Point{X: 0, Y: 0}
		rect.Max = image.Point{X: 64, Y: 64}
	case 2:
		rect.Min = image.Point{X: 64, Y: 0}
		rect.Max = image.Point{X: 64 * 2, Y: 64}
	case 1:
		rect.Min = image.Point{X: 64 * 2, Y: 0}
		rect.Max = image.Point{X: 64 * 3, Y: 64}
	}
	curSprite = b.img.SubImage(rect).(*ebiten.Image)
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(10, 50)
	txtOpt := text.DrawOptions{}
	txtOpt.GeoM.Translate(90, 50)
	txtOpt.ColorScale.ScaleWithColor(RED)
	// txtOpt.LineSpacing = 10
	txtOpt.LayoutOptions = text.LayoutOptions{PrimaryAlign: text.AlignCenter, LineSpacing: 20}
	txtOpt.GeoM.Scale(0.5, 0.5)
	text.Draw(card, "Bomb", face, &txtOpt)
	card.DrawImage(curSprite, &opt)
}
func (b *BombDecorator) GetType() CardType {
	return CARD_TYPE_TRAP
}
func (b *BombDecorator) OnClick(mainScene *MainScene, source Card) {
	playerCard := mainScene.zones[PLAYER_IDX_Y][PLAYER_IDX_X]
	posX, posY := source.(*BaseCard).GetPos()
	playerX, playerY := playerCard.GetPos()
	idxX, idxY := PixelToIndex(int(posX), int(posY))
	playerCard.AddAnimation(core.NewMoveAnimationFromParam(core.MoveParam{Tx: posX, Ty: posY, Speed: CARD_MOVE_SPEED}))
	source.(*BaseCard).AddAnimation(core.NewMoveAnimationFromParam(core.MoveParam{
		Tx:    playerX,
		Ty:    playerY,
		Speed: CARD_MOVE_SPEED,
	}))
	mainScene.zones[idxY][idxX], mainScene.zones[PLAYER_IDX_Y][PLAYER_IDX_X] = playerCard, source.(*BaseCard)
	PLAYER_IDX_X = idxX
	PLAYER_IDX_Y = idxY
	mainScene.OnPlayerMove()
}
func (b *BombDecorator) GetDescription() string {
	return fmt.Sprintf("Bomb that will damage\nadjacent cards for 4\nWill explode in %d turn", b.TurnToExplode)
}

// TODO: get index of certain
func GetPos(zones [3][3]*BaseCard, card *BaseCard) (int, int) {
	for idx, _ := range zones {
		for idx2, s := range zones[idx] {
			if s == card {
				return idx2, idx
			}
		}
	}
	return -1, -1
}
func IdxToPixel(x, y int) (int, int) {
	return BORDER_X[x], BORDER_Y[y]
}
func (b *BombDecorator) OnPlayerMove(card Card, s *MainScene) {
	b.TurnToExplode -= 1
	if b.TurnToExplode == 0 {
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
					v.TakeDirectDamage(4, s, s.zones[i][j])
				}
			}
		}
		card.RemoveDecorator(b)
		card.AddDecorator(NewCoinDecorator())
	}
}
