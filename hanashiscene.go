package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
	"github.com/kharism/hanashi/core"
)

type HanashiScene struct {
	scene      *core.Scene
	skipButton *MenuButton
	director   *stagehand.SceneDirector[MyState]
}

func (m *HanashiScene) Update() error {
	m.skipButton.Update()
	e := m.scene.Update()
	if e != nil {
		return e
	}

	return nil
}
func (m *HanashiScene) Draw(screen *ebiten.Image) {
	m.scene.Draw(screen)
	m.skipButton.Draw(screen)
}

func (s *HanashiScene) Load(state MyState, director stagehand.SceneController[MyState]) {
	s.director = director.(*stagehand.SceneDirector[MyState]) // This type assertion is important
	s.skipButton = &MenuButton{
		MovableImage: core.NewMovableImage(BtnBg, core.NewMovableImageParams()),
		onClickFunc: func() {
			s.director.ProcessTrigger(TriggerToMain)
		},
	}
	s.skipButton.Label = "Skip"
	s.skipButton.SetPos(400, 25)
}
func (s *HanashiScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}
func (s *HanashiScene) Unload() MyState {
	// your unload code
	return MyState{}
}
