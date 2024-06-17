package main

import (
	"bytes"
	_ "embed"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kharism/hanashi/core"
)

//go:embed assets/img/background1.png
var bgimg1 []byte
var bg1 *ebiten.Image

//go:embed assets/img/bradamante.png
var bradamanteImg []byte
var bradamante *ebiten.Image

//go:embed assets/img/cutscene2.png
var bgimg2 []byte
var bg2 *ebiten.Image

func ReadImgFromBytes(raw []byte) *ebiten.Image {
	reader := bytes.NewReader(raw)
	h, _, _ := ebitenutil.NewImageFromReader(reader)
	return h
}
func Scene1(layouter core.GetLayouter) *core.Scene {
	scene := core.NewScene()
	scene.SetLayouter(layouter)

	bg1 = ReadImgFromBytes(bgimg1)
	// knight = ReadImgFromBytes(knightImg)
	bradamante = ReadImgFromBytes(bradamanteImg)

	bg2 = ReadImgFromBytes(bgimg2)
	scene.Characters = []*core.Character{
		core.NewCharacterImage("Charlie", knightImg),
		core.NewCharacterImage("Isolde", bradamante),
	}
	scene.FontFace = face
	scene.Events = []core.Event{
		&core.ComplexEvent{Events: []core.Event{
			core.NewBgChangeEvent(bg1, core.MoveParam{Sx: 0, Sy: 200, Tx: 0, Ty: 0, Speed: 3}, nil),
			core.NewCharacterAddEvent("Charlie", &core.MoveParam{440, 440, 440, 240, 3}, &core.ScaleParam{-1.4, 1.4, 32, 32}),
			core.NewCharacterAddEvent("Isolde", &core.MoveParam{140, 440, 140, 240, 3}, &core.ScaleParam{1.4, 1.4, 32, 32}),
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterMoveEvent("Charlie", core.MoveParam{440, 240, 340, 240, 5}),
			core.NewDialogueEvent("Charlie", "Lady Isolde...", face),
		}},
		core.NewDialogueEvent("Charlie", "May I take your hand in marriage?", face),
		&core.ComplexEvent{Events: []core.Event{
			&core.CharacterComplexMoveEvent{Name: "Isolde", AnimationQueue: []core.Animation{
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 220, Speed: 7}),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 240, Speed: 7}).SetSleepPost(1 * time.Second),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 220, Speed: 7}),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 240, Speed: 7}),
			}},
			core.NewDialogueEvent("Isolde", "THAT'S THE 5th TIME YOU PROPOSE,\nTHIS WEEK ALONE!!!!", face),
		}},
		&core.ComplexEvent{Events: []core.Event{
			&core.CharacterComplexMoveEvent{Name: "Charlie", AnimationQueue: []core.Animation{
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 340, Ty: 220, Speed: 7}),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 340, Ty: 240, Speed: 7}).SetSleepPost(1 * time.Second),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 340, Ty: 220, Speed: 7}),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 340, Ty: 240, Speed: 7}),
			}},
			core.NewDialogueEvent("Charlie", "I believe you will change your\nmind perhaps", face),
		}},
		core.NewDialogueEvent("Isolde", "(*How on earth I shoo him\nfor good??*)", face),
		&core.ComplexEvent{Events: []core.Event{
			&core.CharacterComplexMoveEvent{Name: "Isolde", AnimationQueue: []core.Animation{
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 220, Speed: 7}),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 240, Speed: 7}).SetSleepPost(1 * time.Second),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 220, Speed: 7}),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 240, Speed: 7}),
			}},
			core.NewDialogueEvent("Isolde", "Fine, then here's the deal", face),
		}},
		&core.ComplexEvent{Events: []core.Event{
			&core.CharacterComplexMoveEvent{Name: "Isolde", AnimationQueue: []core.Animation{
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 220, Speed: 7}),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 240, Speed: 7}).SetSleepPost(1 * time.Second),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 220, Speed: 7}),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 240, Speed: 7}),
			}},
			core.NewDialogueEvent("Isolde", "If you can gather 80 coins\nbefore next week", face),
		}},
		&core.ComplexEvent{Events: []core.Event{
			&core.CharacterComplexMoveEvent{Name: "Isolde", AnimationQueue: []core.Animation{
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 220, Speed: 7}),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 240, Speed: 7}).SetSleepPost(1 * time.Second),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 220, Speed: 7}),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 240, Speed: 7}),
			}},
			core.NewDialogueEvent("Isolde", "Then, I'll be your bride", face),
		}},
		&core.ComplexEvent{Events: []core.Event{
			&core.CharacterComplexMoveEvent{Name: "Isolde", AnimationQueue: []core.Animation{
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 220, Speed: 7}),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 240, Speed: 7}).SetSleepPost(1 * time.Second),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 220, Speed: 7}),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 240, Speed: 7}),
			}},
			core.NewDialogueEvent("Isolde", "But if you can't, just leave me\nalone", face),
		}},
		&core.ComplexEvent{Events: []core.Event{
			&core.CharacterComplexMoveEvent{Name: "Isolde", AnimationQueue: []core.Animation{
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 220, Speed: 7}),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 240, Speed: 7}).SetSleepPost(1 * time.Second),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 220, Speed: 7}),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 240, Speed: 7}),
			}},
			core.NewDialogueEvent("Isolde", "The brandish maiden has announced\na gauntlet on her castle", face),
		}},
		&core.ComplexEvent{Events: []core.Event{
			&core.CharacterComplexMoveEvent{Name: "Isolde", AnimationQueue: []core.Animation{
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 220, Speed: 7}),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 240, Speed: 7}).SetSleepPost(1 * time.Second),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 220, Speed: 7}),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 140, Ty: 240, Speed: 7}),
			}},
			core.NewDialogueEvent("Isolde", "You might get bountyfull chest\nfull of gold", face),
		}},
		&core.ComplexEvent{Events: []core.Event{
			&core.CharacterComplexMoveEvent{Name: "Charlie", AnimationQueue: []core.Animation{
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 340, Ty: 220, Speed: 7}),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 340, Ty: 240, Speed: 7}).SetSleepPost(1 * time.Second),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 340, Ty: 220, Speed: 7}),
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 340, Ty: 240, Speed: 7}),
			}},
			core.NewDialogueEvent("Charlie", "Anything for you my\nlove", face),
		}},
		&core.ComplexEvent{Events: []core.Event{
			&core.CharacterComplexMoveEvent{Name: "Charlie", AnimationQueue: []core.Animation{
				&core.ScaleAnimation{Tsy: 1.4, Tsx: 1.4, SpeedX: 0.1, SpeedY: 0.1, CenterX: 32, CenterY: 32},
				core.NewMoveAnimationFromParam(core.MoveParam{Tx: 700, Ty: 240, Speed: 10}),
			}},
			core.NewDialogueEvent("Charlie", "OFF I GO!!!!", face),
		}},
		// core.NewCharacterMoveEvent("Charlie", core.MoveParam{Tx: 700, Ty: 240, Speed: 10}),
		&core.ComplexEvent{
			Events: []core.Event{
				core.NewCharacterRemoveEvent("Isolde"),
				core.NewBgChangeEvent(bg2, core.MoveParam{Sx: 0, Sy: -200, Tx: 0, Ty: 0, Speed: 0.5}, nil),
				core.NewDialogueEvent("", "Here, our heroes, Charlie went to\nbrandish maiden's castle", face),
			},
		},
		core.NewDialogueEvent("", "What kind of obstacle lays's\nbetween him and isolde?", face),
	}

	scene.TxtBg = ebiten.NewImage(640, 200)
	scene.TxtBg.Fill(color.Black)
	scene.Events[0].Execute(scene)
	return scene
}
