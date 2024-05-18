package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

type Game struct{}

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
	state := MyState{
		PlayerCharacter: NewKnightDecor(),
	}

	ruleSet := make(map[stagehand.Scene[MyState]][]stagehand.Directive[MyState])

	manager := stagehand.NewSceneDirector[MyState](scene1, state, ruleSet)

	if err := ebiten.RunGame(manager); err != nil {
		log.Fatal(err)
	}
	// knight = NewKnightDecor()
	// card = NewBaseCard([]CardDecorator{knight})
	// if err := ebiten.RunGame(&Game{}); err != nil {
	// 	log.Fatal(err)
	// }
}
