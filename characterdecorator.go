package main

import (
	"bytes"
	_ "embed"
	"fmt"
	_ "image/png"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kharism/hanashi/core"
)

//go:embed assets/img/Knight.png
var KnightImage []byte

//go:embed assets/img/goblin.png
var GoblinImage []byte

//go:embed assets/fonts/PixelOperator8.ttf
var PixelFontTTF []byte

//go:embed assets/img/skeleton.png
var SkeletonImage []byte

//go:embed assets/img/org.png
var OrgImage []byte

//go:embed assets/img/pyro-eyes.png
var PyroEyesImage []byte

//go:embed assets/img/salamandragon.png
var SalamandragonImage []byte

//go:embed assets/img/ladylab.png
var BrandishMaiden []byte

//go:embed shaders/dakka.kage
var DakkaShader []byte

//go:embed assets/img/attack2.png
var Attack2Png []byte

var PixelFont *text.GoTextFaceSource

var knightImg *ebiten.Image
var goblinImg *ebiten.Image
var skeletonImg *ebiten.Image
var orgImage *ebiten.Image
var pyroEyesImg *ebiten.Image
var salamandragonImg *ebiten.Image
var brandishMaidenImg *ebiten.Image
var face *text.GoTextFace
var dakkaShader *ebiten.Shader
var attackImg *core.AnimatedImage

type CharacterDecorator struct {
	Hp          int
	MaxHP       int
	Name        string
	image       *ebiten.Image
	OnClickFunc OnInteractFunction
	OnDefeat    OnDefeatFunc
	Description string
	shader      *ebiten.Shader
}

type CharacterInterface interface {
	// take direct damage ignoring loadout
	TakeDirectDamage(int, *MainScene, Card)
	// take damage but still putting loadout into consideration
	TakeDamage(int, *MainScene, Card)
	Draw(card *ebiten.Image)
	GetHP() int
	SetHP(int)
	GetMaxHP() int
	SetMaxHP(int)
	GetOnDefeat() OnDefeatFunc

	DoBattle(CharacterInterface, *MainScene)
}
type OnInteractFunction func(*MainScene, Card)
type OnDefeatFunc func(*MainScene, Card)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(PixelFontTTF))
	if err != nil {
		log.Fatal(err)
	}
	PixelFont = s
	face = &text.GoTextFace{
		Source: PixelFont,
		Size:   24,
	}
	if knightImg == nil {
		imgReader := bytes.NewReader(KnightImage)
		knightImg, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if goblinImg == nil {
		imgReader := bytes.NewReader(GoblinImage)
		goblinImg, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if skeletonImg == nil {
		imgReader := bytes.NewReader(SkeletonImage)
		skeletonImg, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if orgImage == nil {
		imgReader := bytes.NewReader(OrgImage)
		orgImage, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	dakkaShader, err = ebiten.NewShader(DakkaShader)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
	if pyroEyesImg == nil {
		imgReader := bytes.NewReader(PyroEyesImage)
		pyroEyesImg, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if salamandragonImg == nil {
		imgReader := bytes.NewReader(SalamandragonImage)
		salamandragonImg, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if brandishMaidenImg == nil {
		imgReader := bytes.NewReader(BrandishMaiden)
		brandishMaidenImg, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if attackImg == nil {
		imgReader := bytes.NewReader(Attack2Png)
		attackImgAll, _, _ := ebitenutil.NewImageFromReader(imgReader)
		attackAnimParam := core.NewMovableImageParams().WithMoveParam(core.MoveParam{Sx: 0, Sy: 0}).WithScale(&core.ScaleParam{Sx: 0.25, Sy: 0.25})
		attackImg = &core.AnimatedImage{
			MovableImage:   core.NewMovableImage(attackImgAll, attackAnimParam),
			SubImageStartX: 0,
			SubImageStartY: 0,
			SubImageWidth:  512,
			SubImageHeight: 512,
			FrameCount:     16,
			Modulo:         2,
		}
	}
}
func (d *CharacterDecorator) GetMaxHP() int {
	return d.MaxHP
}
func (d *CharacterDecorator) TakeDirectDamage(dmg int, s *MainScene, source Card) {
	if dmg <= 0 && d.Hp-dmg > d.MaxHP {
		dmg = -(d.MaxHP - d.Hp)
	}
	d.Hp -= dmg
	if d.Hp <= 0 {
		// os.Exit(0)

		d.OnDefeat(s, source)
	}
}
func (d *CharacterDecorator) GetHP() int {
	return d.Hp
}
func (d *CharacterDecorator) TakeDamage(dmg int, s *MainScene, source Card) {
	// the same with take damage
	d.Hp -= dmg
	if d.Hp <= 0 {
		// os.Exit(0)
		d.OnDefeat(s, source)
	}
}
func (d *CharacterDecorator) DoBattle(opp CharacterInterface, scene *MainScene) {
	d.TakeDamage(opp.GetHP(), scene, nil)
	opp.SetHP(0)
}
func NewKnightDecor() CardDecorator {

	return &CharacterDecorator{Hp: 10, MaxHP: 10, image: knightImg, Name: "Knight", OnClickFunc: func(s *MainScene, c Card) {

	}, OnDefeat: func(scene *MainScene, h Card) {
		// os.Exit(0)
		scene.OnDefeat()
	}, Description: "Your Character"}
}
func (d *CharacterDecorator) GetOnDefeat() OnDefeatFunc {
	return d.OnDefeat
}
func (d *CharacterDecorator) SetHP(a int) {
	d.Hp = a
}
func (d *CharacterDecorator) SetMaxHP(a int) {
	d.MaxHP = a
}
func GenerateDoneFunc(s *MainScene, source CharacterInterface, card Card) func() {
	return func() {
		s.Character.DoBattle(source, s)
		// jj := rwdGenerator.GenerateReward(0)
		if source.GetHP() <= 0 {
			source.GetOnDefeat()(s, card)
		}
		s.ShowAtk = false
		s.OnPlayerMove()
	}
}

// return a function that handle
func GenerateCombat() func(*MainScene, Card) {
	return func(s *MainScene, source Card) {
		// posX, posY := source.(*BaseCard).GetPos()
		// idxX, idxY := PixelToIndex(int(posX), int(posY))
		// s.Character.Hp -= damage
		s.ShowAtk = true
		posX, posY := source.(*BaseCard).GetPos()
		attackImg.SetPos(posX-10, posY)
		attackImg.ScaleParam.Sx = 0.3
		attackImg.ScaleParam.Sy = 0.3
		attackImg.Done = GenerateDoneFunc(s, source.(*BaseCard).decorators[0].(CharacterInterface), source)

	}
}
func GenerateReward(tier int) func(*MainScene, Card) {
	return func(s *MainScene, source Card) {
		reward := rwdGenerator.GenerateReward(tier)
		s.MonstersDefeated += 1
		transDecorator := NewTransitionDecorator(source.(*BaseCard).decorators[0], reward, source.(*BaseCard))
		// source.(*BaseCard).decorators[0] = rwdGenerator.GenerateReward(tier)
		source.(*BaseCard).decorators[0] = transDecorator
	}

}
func NewOrgDecor() CardDecorator {
	orgHP := []int{15, 16, 20}
	return &CharacterDecorator{Hp: orgHP[rand.Int()%len(orgHP)], image: orgImage,
		Name:        "Org",
		OnDefeat:    GenerateReward(2),
		OnClickFunc: GenerateCombat(),
		// shader:      dakkaShader,
		Description: "Brute-ish creature"}
}
func NewHopGoblinDecor() CardDecorator {
	goblinHP := []int{6, 7, 8}
	hp := goblinHP[rand.Int()%len(goblinHP)]
	return &CharacterDecorator{Hp: hp, image: goblinImg,
		Name:        "HopGoblin",
		OnDefeat:    GenerateReward(1),
		OnClickFunc: GenerateCombat(),
		shader:      dakkaShader,
		Description: "A tougher goblin"}
}
func NewSalamandragon() CardDecorator {
	return &CharacterDecorator{Hp: 17, image: salamandragonImg,
		Name: "Salamandragon", OnDefeat: func(s *MainScene, source Card) {
			reward := rwdGenerator.GenerateReward(1)
			posX, posY := source.(*BaseCard).GetPos()
			idxX, idxY := PixelToIndex(int(posX), int(posY))
			for i := 0; i < 3; i++ {
				if i == idxY {
					continue
				}
				if v, ok := s.zones[i][idxX].decorators[0].(CharacterInterface); ok {
					v.TakeDirectDamage(4, s, s.zones[i][idxX])
				}
			}
			s.MonstersDefeated += 1
			transDecorator := NewTransitionDecorator(source.(*BaseCard).decorators[0], reward, source.(*BaseCard))
			// source.(*BaseCard).decorators[0] = rwdGenerator.GenerateReward(tier)
			source.(*BaseCard).decorators[0] = transDecorator
		}, OnClickFunc: GenerateCombat(), Description: "inflict 4 damage to\nall cards on its\ncolumn when defeated"}
}
func NewGoblinDecor() CardDecorator {
	goblinHP := []int{2, 3, 4}

	return &CharacterDecorator{Hp: goblinHP[rand.Int()%len(goblinHP)], image: goblinImg,
		Name: "Goblin", OnDefeat: GenerateReward(0), OnClickFunc: GenerateCombat(), Description: "A small goblin"}
}
func NewXOrg() CardDecorator {
	xorgHP := []int{25, 30, 35}
	hp := xorgHP[rand.Int()%len(xorgHP)]
	org := &CharacterDecorator{Hp: hp, image: orgImage,
		Name:        "XOrg",
		OnDefeat:    GenerateReward(1),
		OnClickFunc: GenerateCombat(),
		shader:      dakkaShader,
		Description: "Org but redder\nIt means it is\nThree times as fast"}
	weakness := &RotatingWeaknessDecorator{&WeaknessDecorator{decorator: org, Direction: DIRECTION_ARR_TWO_WEAKNESS[0]}}
	return weakness
}
func NewSkeletonDecor() CardDecorator {
	if SkeletonImage == nil {
		imgReader := bytes.NewReader(SkeletonImage)
		skeletonImg, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	skletonHP := []int{1, 2, 3}

	return &CharacterDecorator{
		Hp:    skletonHP[rand.Int()%len(skletonHP)],
		image: skeletonImg, Name: "Skeltn",
		OnClickFunc: GenerateCombat(),
		OnDefeat:    GenerateReward(0),
		Description: "A small Skeleton",
	}
}
func NewBrandishMaiden() CardDecorator {
	brandMaidenHp := []int{35}
	pyroEyes := &CharacterDecorator{
		Hp:    brandMaidenHp[rand.Int()%len(brandMaidenHp)],
		image: brandishMaidenImg, Name: "Bradish\nMaiden",
		OnClickFunc: GenerateCombat(),
		OnDefeat: func(ms *MainScene, source Card) {
			// reward := rwdGenerator.GenerateReward(1)
			// ms.MonstersDefeated += 1
			// ms.LastDefeatedMiniBoss = "Brandish Maiden"
			// transDecorator := NewTransitionDecorator(source.(*BaseCard).decorators[0], reward, source.(*BaseCard))
			// source.(*BaseCard).decorators[0] = rwdGenerator.GenerateReward(tier)
			// source.(*BaseCard).decorators[0] = transDecorator
			ms.OnVictory()
		},
		Description: "attack adjacent card for\n5 each time you move",
	}
	aoeDecor := &AttackOnPlayerMoveDecorator{Damage: 5, CharacterDecorator: pyroEyes}
	return aoeDecor
}
func NewPyroEyesDecor() CardDecorator {
	pyroEyesHp := []int{25}
	pyroEyes := &CharacterDecorator{
		Hp:    pyroEyesHp[rand.Int()%len(pyroEyesHp)],
		image: pyroEyesImg, Name: "Pyro-Eyes",
		OnClickFunc: GenerateCombat(),
		OnDefeat: func(ms *MainScene, source Card) {
			reward := rwdGenerator.GenerateReward(1)
			ms.MonstersDefeated += 1
			ms.LastDefeatedMiniBoss = "Pyro-eyes"
			transDecorator := NewTransitionDecorator(source.(*BaseCard).decorators[0], reward, source.(*BaseCard))
			// source.(*BaseCard).decorators[0] = rwdGenerator.GenerateReward(tier)
			source.(*BaseCard).decorators[0] = transDecorator
		},
		Description: "attack adjacent card for\n3 each time you move",
	}
	aoeDecor := &AttackOnPlayerMoveDecorator{Damage: 3, CharacterDecorator: pyroEyes}
	return aoeDecor
}
func (k *CharacterDecorator) Update() error {
	return nil
}
func (k *CharacterDecorator) GetType() CardType {
	return CARD_TYPE_CHARACTER
}

func (k *CharacterDecorator) OnClick(state *MainScene, source Card) {
	// return CARD_TYPE_CHARACTER
	k.OnClickFunc(state, source)
}
func (k *CharacterDecorator) GetDescription() string {
	return k.Description
}
func (k *CharacterDecorator) Draw(card *ebiten.Image) {
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(10, 50)

	txtOpt := text.DrawOptions{}
	txtOpt.GeoM.Scale(0.5, 0.5)
	txtOpt.GeoM.Translate(80, 6)
	txtOpt.PrimaryAlign = text.AlignEnd
	txtOpt.ColorScale.ScaleWithColor(RED)
	if k.MaxHP > 0 {
		text.Draw(card, fmt.Sprintf("%d/%d", k.Hp, k.MaxHP), face, &txtOpt)
	} else {
		text.Draw(card, fmt.Sprintf("%d", k.Hp), face, &txtOpt)
	}

	// txtOpt.GeoM.Reset()
	txtOpt = text.DrawOptions{}
	txtOpt.PrimaryAlign = text.AlignCenter
	txtOpt.GeoM.Scale(0.7, 0.7)
	txtOpt.GeoM.Translate(70, 56)
	txtOpt.LineSpacing = 25

	txtOpt.ColorScale.ScaleWithColor(RED)
	txtOpt.GeoM.Scale(0.6, 0.6)
	text.Draw(card, k.Name, face, &txtOpt)
	if k.shader == nil {
		card.DrawImage(k.image, &opt)
	} else {
		opts := &ebiten.DrawRectShaderOptions{}
		// if e.ScaleParam != nil {
		// 	opts.GeoM.Scale(e.ScaleParam.Sx, e.ScaleParam.Sy)
		// }
		opts.GeoM.Translate(10, 50)
		opts.Images[0] = k.image
		bounds := k.image.Bounds()
		card.DrawRectShader(bounds.Dx(), bounds.Dy(), k.shader, opts)
	}

}
