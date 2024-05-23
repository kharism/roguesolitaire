package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

type Game struct{}

const (
	TriggerToMain stagehand.SceneTransitionTrigger = iota
	TriggerToMenu
	TriggerToSum
)

var knight CardDecorator
var card Card

type MyState struct {
	PlayerCharacter CardDecorator
	Coin            int
	MainScene       *MainScene
}

func (g *Game) Update() error {

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// card.(*BaseCard).MovableImage.ScaleParam = &core.ScaleParam{Sx: 0.8, Sy: 0.8}
	// card.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	scene1 := &MainScene{}
	menuScene := &MenuScene{}
	state := MyState{
		PlayerCharacter: NewKnightDecor(),
	}
	// state.Coin = 117
	summary := &SummaryScene{}
	trans := stagehand.NewSlideTransition[MyState](stagehand.LeftToRight, 0.05)
	trans2 := stagehand.NewSlideTransition[MyState](stagehand.RightToLeft, 0.05)
	trans3 := stagehand.NewFadeTransition[MyState](0.3)
	// ruleSet := make(map[stagehand.Scene[MyState]][]stagehand.Directive[MyState])
	ruleSet := map[stagehand.Scene[MyState]][]stagehand.Directive[MyState]{
		menuScene: []stagehand.Directive[MyState]{
			stagehand.Directive[MyState]{Dest: scene1, Trigger: TriggerToMain, Transition: trans},
		},
		scene1: []stagehand.Directive[MyState]{
			stagehand.Directive[MyState]{Dest: summary, Trigger: TriggerToSum, Transition: trans3},
		},
		summary: []stagehand.Directive[MyState]{
			stagehand.Directive[MyState]{Dest: menuScene, Trigger: TriggerToMenu, Transition: trans2},
			stagehand.Directive[MyState]{Dest: scene1, Trigger: TriggerToMain, Transition: trans3},
		},
	}
	manager := stagehand.NewSceneDirector[MyState](menuScene, state, ruleSet)

	if err := ebiten.RunGame(manager); err != nil {
		log.Fatal(err)
	}
	// knight = NewKnightDecor()
	// card = NewBaseCard([]CardDecorator{knight})
	// if err := ebiten.RunGame(&Game{}); err != nil {
	// 	log.Fatal(err)
	// }
}
