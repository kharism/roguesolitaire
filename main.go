package main

import (
	"bytes"
	"io"
	"log"
	"time"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/joelschutz/stagehand"
	"github.com/kharism/hanashi/core"
)

type Game struct{}

const (
	TriggerToMain stagehand.SceneTransitionTrigger = iota
	TriggerToMenu
	TriggerToSum
	TriggerToOPCutscene
	TriggerToEnding1
	TriggerToEnding2
)
const (
	sampleRate = 48000
)

var knight CardDecorator
var card Card

// var musicPlayer *AudioPlayer
var musicPlayerCh chan *AudioPlayer
var audioContext *audio.Context

type musicType int

const (
	typeOgg musicType = iota
	typeMP3
)

func (t musicType) String() string {
	switch t {
	case typeOgg:
		return "Ogg"
	case typeMP3:
		return "MP3"
	default:
		panic("not reached")
	}
}

type AudioPlayer struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	current      time.Duration
	total        time.Duration
	seBytes      []byte
	seCh         chan []byte

	volume128 int

	musicType musicType
}

func (p *AudioPlayer) Close() error {
	return p.audioPlayer.Close()
}
func (p *AudioPlayer) update() error {
	select {
	case p.seBytes = <-p.seCh:
		close(p.seCh)
		p.seCh = nil
	default:
	}
	p.playSEIfNeeded()
	return nil
}
func (p *AudioPlayer) shouldPlaySE() bool {
	if p.seBytes == nil {
		// Bytes for the SE is not loaded yet.
		return false
	}
	return false
}

func (p *AudioPlayer) playSEIfNeeded() {
	// if !p.shouldPlaySE() {
	// 	return
	// }
	// sePlayer := p.audioContext.NewPlayerFromBytes(p.seBytes)
	// sePlayer.Play()

}

//go:embed assets/music/battle-time-178551.mp3
var battleMusic []byte

//go:embed assets/music/pixel-fight-8-bit-arcade-music-background-music-for-video-208775.mp3
var arcadeMusic []byte

func NewAudioPlayer(audioContext *audio.Context, audioByte []byte, musicType musicType) (*AudioPlayer, error) {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}
	const bytesPerSample = 4 // TODO: This should be defined in audio package

	var s audioStream
	// audio, err := os.Open(audioPath)
	// if err != nil {
	// 	return nil, err
	// }
	// defer audio.Close()
	switch musicType {

	case typeMP3:
		var err error
		s, err = mp3.DecodeWithoutResampling(bytes.NewReader(audioByte))
		if err != nil {
			return nil, err
		}
	default:
		panic("not reached")
	}

	p, err := audioContext.NewPlayer(s)
	if err != nil {
		return nil, err
	}

	player := &AudioPlayer{
		audioContext: audioContext,
		audioPlayer:  p,
		total:        time.Second * time.Duration(s.Length()) / bytesPerSample / sampleRate,
		volume128:    32,
		seCh:         make(chan []byte),
		seBytes:      []byte{},
		musicType:    musicType,
	}
	if player.total == 0 {
		player.total = 1
	}

	// player.audioPlayer.Play()

	return player, nil
}

// which SE we should play
var seSignal chan string

type MyState struct {
	PlayerCharacter CardDecorator
	Coin            int
	Victory         bool
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

type LayouterImpl struct {
	musicPlayer *AudioPlayer
}

var Layout *LayouterImpl

func (l *LayouterImpl) GetLayout() (width, height int) {
	return 640, 480
}
func (l *LayouterImpl) GetNamePosition() (x, y int) {
	return 0, 512 - 150
}
func (l *LayouterImpl) GetTextPosition() (x, y int) {
	return 0, 512 - 120
}

func main() {
	audioContext = audio.NewContext(sampleRate)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Rogue Solitaire")
	scene1 := &MainScene{}
	menuScene := &MenuScene{}
	Layout = &LayouterImpl{}

	cutScene1 := Scene1(Layout)
	endingScene1 := Ending1(Layout)
	endingScene2 := Ending2(Layout)

	HanashiScene1 := &HanashiScene{scene: cutScene1}
	cutScene1.Done = func() {
		HanashiScene1.director.ProcessTrigger(TriggerToMain)
	}

	HanashiScene2 := &HanashiScene{scene: endingScene1}
	HanashiScene2.SkipButton = &MenuButton{
		MovableImage: core.NewMovableImage(BtnBg, core.NewMovableImageParams()),
		onClickFunc: func() {
			HanashiScene2.director.ProcessTrigger(TriggerToSum)
		},
	}
	endingScene1.Done = func() {
		HanashiScene2.director.ProcessTrigger(TriggerToSum)
	}
	HanashiScene3 := &HanashiScene{scene: endingScene2}

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
			// stagehand.Directive[MyState]{Dest: scene1, Trigger: TriggerToMain, Transition: trans},
			stagehand.Directive[MyState]{Dest: HanashiScene1, Trigger: TriggerToOPCutscene, Transition: trans3},
		},
		HanashiScene1: []stagehand.Directive[MyState]{
			stagehand.Directive[MyState]{Dest: scene1, Trigger: TriggerToMain, Transition: trans},
		},
		HanashiScene2: []stagehand.Directive[MyState]{
			stagehand.Directive[MyState]{Dest: summary, Trigger: TriggerToSum, Transition: trans3},
		},
		HanashiScene3: []stagehand.Directive[MyState]{
			stagehand.Directive[MyState]{Dest: summary, Trigger: TriggerToSum, Transition: trans3},
		},
		scene1: []stagehand.Directive[MyState]{
			stagehand.Directive[MyState]{Dest: summary, Trigger: TriggerToSum, Transition: trans3},
			stagehand.Directive[MyState]{Dest: HanashiScene2, Trigger: TriggerToEnding1, Transition: trans3},
			stagehand.Directive[MyState]{Dest: HanashiScene3, Trigger: TriggerToEnding2, Transition: trans3},
		},
		summary: []stagehand.Directive[MyState]{
			stagehand.Directive[MyState]{Dest: menuScene, Trigger: TriggerToMenu, Transition: trans2},
			stagehand.Directive[MyState]{Dest: scene1, Trigger: TriggerToMain, Transition: trans3},
		},
	}
	manager := stagehand.NewSceneDirector[MyState](menuScene, state, ruleSet)
	// musicPlayer, _ = NewAudioPlayer(audioContext, arcadeMusic, typeMP3)

	if err := ebiten.RunGame(manager); err != nil {
		log.Fatal(err)
	}
	// knight = NewKnightDecor()
	// card = NewBaseCard([]CardDecorator{knight})
	// if err := ebiten.RunGame(&Game{}); err != nil {
	// 	log.Fatal(err)
	// }
}
