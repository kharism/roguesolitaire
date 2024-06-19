package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/joelschutz/stagehand"
	"github.com/kharism/hanashi/core"
)

type HanashiScene struct {
	scene      *core.Scene
	State      MyState
	SkipButton *MenuButton
	director   *stagehand.SceneDirector[MyState]
}

func (m *HanashiScene) Update() error {
	m.SkipButton.Update()
	e := m.scene.Update()
	if e != nil {
		return e
	}

	return nil
}
func (m *HanashiScene) Draw(screen *ebiten.Image) {
	m.scene.Draw(screen)
	m.SkipButton.Draw(screen)
	txt := "click to continue"
	txtOpt := text.DrawOptions{}
	txtOpt.ColorScale.ScaleWithColor(RED)
	txtOpt.GeoM.Scale(0.5, 0.5)
	text.Draw(screen, txt, face, &txtOpt)
}

func (s *HanashiScene) Load(state MyState, director stagehand.SceneController[MyState]) {
	s.director = director.(*stagehand.SceneDirector[MyState]) // This type assertion is important
	s.State = state
	if s.SkipButton == nil {
		s.SkipButton = &MenuButton{
			MovableImage: core.NewMovableImage(BtnBg, core.NewMovableImageParams()),
			onClickFunc: func() {
				s.director.ProcessTrigger(TriggerToMain)
			},
		}
	}

	s.SkipButton.Label = "Skip"
	s.SkipButton.SetPos(400, 25)
}
func (s *HanashiScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}
func (s *HanashiScene) Unload() MyState {
	// your unload code

	return s.State
}
