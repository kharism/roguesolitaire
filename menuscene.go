package main

import (
	"bytes"
	_ "embed"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/joelschutz/stagehand"
	"github.com/kharism/hanashi/core"
)

//go:embed assets/img/menu_bg.jpg
var menuBg []byte
var MenuBg *ebiten.Image

//go:embed assets/img/btnbg.png
var btnBg []byte
var BtnBg *ebiten.Image

type MenuScene struct {
	director    *stagehand.SceneDirector[MyState]
	StartGame   MenuButton
	musicPlayer *AudioPlayer
}
type MenuButton struct {
	*core.MovableImage
	Label       string
	cursorIn    bool
	onClickFunc func()
}

func (b *MenuButton) Draw(screen *ebiten.Image) {
	if b.cursorIn {
		b.ScaleParam.Sx = 1.8
	} else {
		b.ScaleParam.Sx = 1.7
	}
	btnX, btnY := b.MovableImage.GetPos()
	b.MovableImage.Draw(screen)
	txtOpt := text.DrawOptions{}
	txtOpt.GeoM.Scale(0.8, 0.8)
	txtOpt.GeoM.Translate(btnX+10, btnY+10)

	txtOpt.ColorScale.ScaleWithColor(RED)
	text.Draw(screen, b.Label, face, &txtOpt)
}
func (b *MenuButton) Update() {
	// curX, curY := ebiten.CursorPosition()
	isClicked, curX, curY := IsClickedOrTap()
	if !isClicked {
		curX, curY = ebiten.CursorPosition()
	}
	butPosX, butPosY := b.GetPos()
	width, height := b.GetSize()
	if curX > int(butPosX) && curX < int(butPosX+width) && curY > int(butPosY) && curY < int(butPosY+height) {
		b.cursorIn = true
		// fmt.Println("Cursor In")
	} else {
		b.cursorIn = false
	}
	if isClicked { //inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		// fmt.Println("Click detected")
		// fmt.Println(curX, int(butPosX), butPosX+width, curY, butPosY, butPosY+height)
		if b.cursorIn && b.onClickFunc != nil {
			b.onClickFunc()
		}
	}
}
func init() {
	reader := bytes.NewReader(menuBg)
	MenuBg, _, _ = ebitenutil.NewImageFromReader(reader)
	reader = bytes.NewReader(btnBg)
	BtnBg, _, _ = ebitenutil.NewImageFromReader(reader)
}
func (m *MenuScene) Update() error {
	m.StartGame.Update()
	return nil
}
func (m *MenuScene) Draw(screen *ebiten.Image) {
	imgOpt := ebiten.DrawImageOptions{}
	screen.DrawImage(MenuBg, &imgOpt)

	m.StartGame.Draw(screen)
}
func (s *MenuScene) Load(state MyState, director stagehand.SceneController[MyState]) {
	s.director = director.(*stagehand.SceneDirector[MyState]) // This type assertion is important
	s.StartGame = MenuButton{}
	s.StartGame.MovableImage = core.NewMovableImage(BtnBg, core.NewMovableImageParams())
	s.StartGame.Label = "Start Game"
	if s.musicPlayer == nil {
		var err error
		s.musicPlayer, err = NewAudioPlayer(audioContext, arcadeMusic, typeMP3)
		if err != nil {
			fmt.Println(err.Error())
		}
		s.musicPlayer.audioPlayer.Play()
	} else {
		// s.musicPlayer.audioPlayer.Rewind()
		s.musicPlayer.audioPlayer.Play()
	}

	// s.musicPlayer.audioPlayer.Play()
	s.StartGame.onClickFunc = func() {
		// Layout.musicPlayer = nil
		s.director.ProcessTrigger(TriggerToOPCutscene)
	}
	s.StartGame.SetPos(230, 250)
}
func (s *MenuScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}
func (s *MenuScene) Unload() MyState {
	// your unload code
	s.musicPlayer.audioPlayer.Rewind()
	s.musicPlayer.audioPlayer.Pause()
	return MyState{}
}
