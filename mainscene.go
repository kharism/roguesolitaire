package main

import (
	"bytes"
	"fmt"
	"image"
	"math"
	"math/rand"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/joelschutz/stagehand"
	"github.com/kharism/hanashi/core"
)

//go:embed assets/img/infobg.png
var InfoBg []byte
var infobg *ebiten.Image

//go:embed assets/img/tilebg.png
var tile []byte
var tileImg *ebiten.Image

func init() {
	imgReader := bytes.NewReader(InfoBg)
	infobg, _, _ = ebitenutil.NewImageFromReader(imgReader)
	imgReader = bytes.NewReader(tile)
	tileImg, _, _ = ebitenutil.NewImageFromReader(imgReader)
}

type MainScene struct {
	director      *stagehand.SceneDirector[MyState]
	State         *MyState
	Character     CharacterInterface
	CharacterCard *BaseCard
	CharacterPosX int
	CharacterPosY int
	touchIDs      []ebiten.TouchID
	zones         [3][3]*BaseCard
	CurDesc       string
	CurMovingCard *BaseCard

	isDefeated      bool
	defeatedCounter int // this counter is only used in animation when loosing

	ShowAtk bool

	MonstersDefeated int //
}

func NewMainScene() *MainScene {
	newScene := &MainScene{}
	newScene.zones = [3][3]*BaseCard{}
	return newScene
}

// turn a point in game (x,y) to index of grid ranging from 0-2
// grid border is based on BORDER_X and BORDER_Y
func PixelToIndex(x, y int) (int, int) {
	idxX := -1
	idxY := -1
	for i := 0; i <= 2; i++ {
		if x >= BORDER_X[i] && x <= BORDER_X[i+1] {
			idxX = i
		}
		if y >= BORDER_Y[i] && y <= BORDER_Y[i+1] {
			idxY = i
		}
	}
	return idxX, idxY
}
func PlayerCanInteractHere(idxX, idxY int) bool {
	dist := math.Abs(float64(idxX-PLAYER_IDX_X)) + math.Abs(float64(idxY-PLAYER_IDX_Y))
	return dist == 1
}
func (m *MainScene) OnPlayerMove() {
	fmt.Println("Player moves")
	for idx, _ := range m.zones {
		for idx2, _ := range m.zones[idx] {
			if v, ok := m.zones[idx][idx2].decorators[0].(PlayerMoveListener); ok {
				v.OnPlayerMove(m.zones[idx][idx2], m)
			}
		}
	}
}
func (m *MainScene) Update() error {
	// mouseX, mouseY := ebiten.CursorPosition()
	// if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	// 	fmt.Println(inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft))
	// }
	isClicked, mouseX, mouseY := IsClickedOrTap()
	if isClicked { //inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		idxX, idxY := PixelToIndex(mouseX, mouseY)
		if PlayerCanInteractHere(idxX, idxY) {
			m.zones[idxY][idxX].OnClick(m)
		}

	} else {
		mouseX, mouseY = ebiten.CursorPosition()
	}
	if m.isDefeated {
		m.defeatedCounter++
	}
	if m.defeatedCounter%60 == 59 {
		// m.defeatedCounter = 0
		m.director.ProcessTrigger(TriggerToSum)
	}
	cardIdx, cardIdy := PixelToIndex(mouseX, mouseY)
	if cardIdx >= 0 && cardIdy >= 0 {
		m.CurDesc = m.zones[cardIdy][cardIdx].GetDescription()
	}
	for idx, _ := range m.zones {
		for idx2, _ := range m.zones[idx] {
			err := m.zones[idx][idx2].Update()
			if err != nil {
				return err
			}

		}
	}
	if m.ShowAtk {
		attackImg.Update()
	}
	return nil
}

const (
	BOARD_START_X = 40
	BOARD_START_Y = 40
	MARGIN_X      = 20
	MARGIN_Y      = 30
)

type Pos struct {
	X float64
	Y float64
}

var (
	CARDS_POS    [][]Pos
	BORDER_X     []int
	BORDER_Y     []int
	PLAYER_IDX_X int
	PLAYER_IDX_Y int
)
var (
	bgInfoStartX = BOARD_START_X + (BASE_CARD_WIDTH*SCALE_CARD)*3 + MARGIN_X*2 + 30
	bgInfoStartY = BOARD_START_Y

	bgInfo2StartX = bgInfoStartX
	bgInfo2StartY = bgInfoStartY + 300
)

func (m *MainScene) DrawInfoBg(screen *ebiten.Image) {
	//top parts
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Scale(1.2, 1)

	opts.GeoM.Translate(bgInfoStartX, float64(bgInfoStartY))
	screen.DrawImage(infobg, &opts)
	opts.GeoM.Reset()
	// opts.GeoM.Rotate(math.Pi)
	opts.GeoM.Scale(1.2, -1)
	opts.GeoM.Translate(bgInfoStartX, float64(bgInfoStartY+35+250)+10)
	screen.DrawImage(infobg, &opts)
	// center parts
	midPart := infobg.SubImage(image.Rect(0, 20, 251, 30))
	opts.GeoM.Reset()
	opts.GeoM.Scale(1.2, 25)
	opts.GeoM.Translate(bgInfoStartX, float64(bgInfoStartY)+35)
	screen.DrawImage(midPart.(*ebiten.Image), &opts)

}
func (m *MainScene) OnDefeat() {
	for idx1, _ := range m.zones {
		for idx2, _ := range m.zones[idx1] {
			scaleAnim := core.ScaleAnimation{Tsx: 0.1, Tsy: 0.1, SpeedX: -0.01, SpeedY: -0.01}
			m.zones[idx1][idx2].AddAnimation(&scaleAnim)
		}
	}
	m.isDefeated = true

}
func (m *MainScene) DrawInfoBg2(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Scale(1.2, 1)

	opts.GeoM.Translate(bgInfo2StartX, float64(bgInfo2StartY))
	screen.DrawImage(infobg, &opts)
	opts.GeoM.Reset()

	opts.GeoM.Scale(1.2, -1)
	opts.GeoM.Translate(bgInfo2StartX, float64(bgInfo2StartY+35+70))
	screen.DrawImage(infobg, &opts)
	midPart := infobg.SubImage(image.Rect(0, 20, 251, 30))
	opts.GeoM.Reset()
	opts.GeoM.Scale(1.2, 5)
	opts.GeoM.Translate(bgInfo2StartX, float64(bgInfo2StartY)+35)
	screen.DrawImage(midPart.(*ebiten.Image), &opts)

}
func (m *MainScene) DrawBg(screen *ebiten.Image) {
	posX := 0
	posY := 0
	opt := ebiten.DrawImageOptions{}
	for i := 0; i < 10; i++ {
		posX = 0
		for j := 0; j < 10; j++ {
			opt.GeoM.Reset()
			opt.GeoM.Translate(float64(posX), float64(posY))
			screen.DrawImage(tileImg, &opt)
			posX += 64
		}
		posY += 64
	}
}
func (m *MainScene) DrawDesc(screen *ebiten.Image) {
	txtOpt := text.DrawOptions{}
	txtOpt.GeoM.Scale(0.7, 0.7)
	txtOpt.GeoM.Translate(bgInfoStartX+5, float64(bgInfoStartY)+35)
	txtOpt.LineSpacing = 24
	txtOpt.ColorScale.ScaleWithColor(RED)
	text.Draw(screen, m.CurDesc, face, &txtOpt)
}

func (m *MainScene) Draw(screen *ebiten.Image) {
	m.DrawBg(screen)
	m.DrawInfoBg(screen)
	for idx, _ := range m.zones {
		for idx2, b := range m.zones[idx] {
			if idx == PLAYER_IDX_Y && idx2 == PLAYER_IDX_X {
				continue
			}
			if m.zones[idx][idx2] == m.CurMovingCard {
				continue
			}

			b.Draw(screen)
		}
	}
	if m.CurDesc != "" {
		m.DrawDesc(screen)
	}
	m.zones[PLAYER_IDX_Y][PLAYER_IDX_X].Draw(screen)
	if m.CurMovingCard != nil {
		m.CurMovingCard.Draw(screen)
	}

	m.DrawInfoBg2(screen)
	if m.ShowAtk {
		attackImg.Draw(screen)
	}
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(bgInfoStartX, float64(bgInfoStartY+35+250)+10)
	screen.DrawImage(coinImg, &opt)
	txtOpt := text.DrawOptions{}
	txtOpt.GeoM.Translate(bgInfoStartX+60, float64(bgInfoStartY+35+250)+30)
	txtOpt.ColorScale.ScaleWithColor(RED)
	text.Draw(screen, fmt.Sprintf("%d", m.State.Coin), face, &txtOpt)
}

func (s *MainScene) Load(state MyState, director stagehand.SceneController[MyState]) {
	// your load code
	s.director = director.(*stagehand.SceneDirector[MyState]) // This type assertion is important
	s.State = &state
	s.isDefeated = false
	s.MonstersDefeated = 0
	s.CurMovingCard = nil
	BORDER_X = make([]int, 4)
	BORDER_Y = make([]int, 4)
	// s.zones[1][1]
	for idx, _ := range s.zones {
		BORDER_Y[idx] = BOARD_START_Y + BASE_CARD_HEIGHT*SCALE_CARD*idx + idx*MARGIN_Y
		for idx2, _ := range s.zones[idx] {
			BORDER_X[idx2] = BOARD_START_X + BASE_CARD_WIDTH*SCALE_CARD*idx2 + idx2*MARGIN_X
			xPos := BOARD_START_X + BASE_CARD_WIDTH*SCALE_CARD*idx2 + idx2*MARGIN_X
			yPos := BOARD_START_Y + BASE_CARD_HEIGHT*SCALE_CARD*idx + idx*MARGIN_Y
			if idx == 1 && idx2 == 1 {
				pp := NewKnightDecor()
				SwordedKnight := NewSwordChDecorator(pp.(*CharacterDecorator), 5)
				s.Character = SwordedKnight.(*SwordChDecorator)
				s.zones[idx][idx2] = NewBaseCard([]CardDecorator{SwordedKnight}).(*BaseCard)
				s.CharacterCard = s.zones[idx][idx2]
			} else if idx == 0 && idx2 == 0 {
				org := NewOrgDecor()
				direction := 1
				org = NewWeaknessDecorator(org, byte(direction))
				s.zones[idx][idx2] = NewBaseCard([]CardDecorator{org}).(*BaseCard)
			} else {
				i := rand.Int() % 3
				if i == 0 {
					s.zones[idx][idx2] = NewBaseCard([]CardDecorator{NewSkeletonDecor()}).(*BaseCard)
				} else if i == 1 {
					// s.zones[idx][idx2] = NewBaseCard([]CardDecorator{NewBombDecorator()}).(*BaseCard)
					// s.zones[idx][idx2] = NewBaseCard([]CardDecorator{NewSpikeTrapDecorator()}).(*BaseCard)
					s.zones[idx][idx2] = NewBaseCard([]CardDecorator{NewChestDecorator()}).(*BaseCard)
				} else if i == 2 {
					// HopDecor := NewHopGoblinDecor()
					// weakness := NewWeaknessDecorator(HopDecor, DIRECTION_UP|DIRECTION_DOWN)
					s.zones[idx][idx2] = NewBaseCard([]CardDecorator{NewGoblinDecor()}).(*BaseCard)
					// s.zones[idx][idx2] = NewBaseCard([]CardDecorator{weakness}).(*BaseCard)
					// s.zones[idx][idx2] = NewBaseCard([]CardDecorator{NewSwordDecorator()}).(*BaseCard)
				}

			}
			PLAYER_IDX_X = 1
			PLAYER_IDX_Y = 1
			s.zones[idx][idx2].SetPos(float64(xPos), float64(yPos))
		}
		BORDER_X[3] = BOARD_START_X + BASE_CARD_WIDTH*SCALE_CARD*3 + 3*MARGIN_X
	}
	BORDER_Y[3] = BOARD_START_Y + BASE_CARD_HEIGHT*SCALE_CARD*3 + 3*MARGIN_Y
}
func (s *MainScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}
func (s *MainScene) Unload() MyState {
	// your unload code
	return *s.State
}
