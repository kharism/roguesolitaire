package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/joelschutz/stagehand"
)

type MainScene struct {
	director      *stagehand.SceneDirector[MyState]
	State         *MyState
	Character     *CharacterDecorator
	CharacterCard *BaseCard
	CharacterPosX int
	CharacterPosY int
	touchIDs      []ebiten.TouchID
	zones         [3][3]*BaseCard
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
func (m *MainScene) Update() error {
	// if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	// 	fmt.Println(inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft))
	// }
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		idxX, idxY := PixelToIndex(mouseX, mouseY)
		if PlayerCanInteractHere(idxX, idxY) {
			m.zones[idxY][idxX].OnClick(m)
		}

	}
	for idx, _ := range m.zones {
		for idx2, _ := range m.zones[idx] {
			err := m.zones[idx][idx2].Update()
			if err != nil {
				return err
			}
		}
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

func (m *MainScene) Draw(screen *ebiten.Image) {
	for idx, _ := range m.zones {
		for idx2, b := range m.zones[idx] {
			if idx == PLAYER_IDX_Y && idx2 == PLAYER_IDX_X {
				continue
			}
			b.Draw(screen)
		}
	}
	m.zones[PLAYER_IDX_Y][PLAYER_IDX_X].Draw(screen)
}

func (s *MainScene) Load(state MyState, director stagehand.SceneController[MyState]) {
	// your load code
	s.director = director.(*stagehand.SceneDirector[MyState]) // This type assertion is important
	s.State = &state

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
				s.Character = pp.(*CharacterDecorator)
				s.zones[idx][idx2] = NewBaseCard([]CardDecorator{pp}).(*BaseCard)
				s.CharacterCard = s.zones[idx][idx2]
			} else {
				i := rand.Int() % 3
				if i == 0 {
					s.zones[idx][idx2] = NewBaseCard([]CardDecorator{NewLightPotionDecorator()}).(*BaseCard)
				} else if i == 1 {
					s.zones[idx][idx2] = NewBaseCard([]CardDecorator{NewSpikeTrapDecorator()}).(*BaseCard)
					// s.zones[idx][idx2] = NewBaseCard([]CardDecorator{NewChestDecorator()}).(*BaseCard)
				} else if i == 2 {
					s.zones[idx][idx2] = NewBaseCard([]CardDecorator{NewSkeletonDecor()}).(*BaseCard)
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
	return MyState{}
}
