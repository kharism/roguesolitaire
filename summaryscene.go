package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/joelschutz/stagehand"
	"github.com/kharism/hanashi/core"
)

type SummaryScene struct {
	director           *stagehand.SceneDirector[MyState]
	knightAnim         *ebiten.Image
	state              MyState
	counter            int
	visibleCoinCounter int
	BackToMenu         MenuButton
	PlayAgain          MenuButton
}

//go:embed assets/img/Knight_death.png
var KnightDieAnim []byte
var knightdieanim *ebiten.Image

func (s *SummaryScene) Update() error {
	s.counter += 1
	if s.counter == 10 {
		rect := image.Rectangle{Min: image.Point{64, 0}, Max: image.Point{64 * 2, 64}}
		s.knightAnim = knightdieanim.SubImage(rect).(*ebiten.Image)
	}
	if s.counter == 20 {
		rect := image.Rectangle{Min: image.Point{64 * 2, 0}, Max: image.Point{64 * 3, 64}}
		s.knightAnim = knightdieanim.SubImage(rect).(*ebiten.Image)
	}
	if s.counter == 30 {
		rect := image.Rectangle{Min: image.Point{64 * 3, 0}, Max: image.Point{64 * 4, 64}}
		s.knightAnim = knightdieanim.SubImage(rect).(*ebiten.Image)
	}
	if s.counter >= 60 {
		if s.visibleCoinCounter < s.state.Coin {
			diff := s.state.Coin - s.visibleCoinCounter
			if diff < 5 {
				s.visibleCoinCounter += diff
			} else {
				s.visibleCoinCounter += 5
			}

		}
	}
	s.BackToMenu.Update()
	s.PlayAgain.Update()
	return nil
}
func (s *SummaryScene) DrawBgTile(screen *ebiten.Image) {
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
func (m *SummaryScene) DrawInfoBg2(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Scale(2, 1)

	opts.GeoM.Translate(80, 20)
	screen.DrawImage(infobg, &opts)
	opts.GeoM.Reset()

	opts.GeoM.Scale(2, -1)
	opts.GeoM.Translate(80, 400)
	screen.DrawImage(infobg, &opts)
	midPart := infobg.SubImage(image.Rect(0, 20, 251, 30))
	opts.GeoM.Reset()
	opts.GeoM.Scale(2, 33)
	opts.GeoM.Translate(80, 40)
	screen.DrawImage(midPart.(*ebiten.Image), &opts)

}
func (s *SummaryScene) Draw(screen *ebiten.Image) {
	s.DrawBgTile(screen)
	s.DrawInfoBg2(screen)
	drwOpt := ebiten.DrawImageOptions{}
	drwOpt.GeoM.Translate(140, 90)
	screen.DrawImage(s.knightAnim, &drwOpt)
	drwOpt.GeoM.Reset()
	drwOpt.GeoM.Translate(210, 90)
	screen.DrawImage(coinImg, &drwOpt)
	txtDrawOpt := text.DrawOptions{}
	txtDrawOpt.PrimaryAlign = text.AlignStart
	txtDrawOpt.GeoM.Translate(260, 110)
	text.Draw(screen, fmt.Sprintf("%d", s.visibleCoinCounter), face, &txtDrawOpt)

	s.BackToMenu.Draw(screen)
	s.PlayAgain.Draw(screen)
}

func (s *SummaryScene) Load(state MyState, director stagehand.SceneController[MyState]) {
	s.director = director.(*stagehand.SceneDirector[MyState]) // This type assertion is important
	reader := bytes.NewReader(KnightDieAnim)
	knightdieanim, _, _ = ebitenutil.NewImageFromReader(reader)
	rect := image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{64, 64}}
	s.knightAnim = knightdieanim.SubImage(rect).(*ebiten.Image)
	s.state = state
	s.counter = 0
	s.visibleCoinCounter = 0

	s.BackToMenu = MenuButton{}
	s.BackToMenu.MovableImage = core.NewMovableImage(BtnBg, core.NewMovableImageParams())
	s.BackToMenu.Label = "To Menu"
	s.BackToMenu.SetPos(230, 250)
	s.BackToMenu.onClickFunc = func() {
		s.director.ProcessTrigger(TriggerToMenu)
	}

	s.PlayAgain = MenuButton{}
	s.PlayAgain.MovableImage = core.NewMovableImage(BtnBg, core.NewMovableImageParams())
	s.PlayAgain.Label = "Play Again"
	s.PlayAgain.SetPos(230, 290)
	s.PlayAgain.onClickFunc = func() {
		s.director.ProcessTrigger(TriggerToMain)
	}

}
func (s *SummaryScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}
func (s *SummaryScene) Unload() MyState {
	// your unload code
	return MyState{}
}
