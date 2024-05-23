package main

import "math/rand"

var generator cardGenerator
var rwdGenerator rewardGenerator

type cardGenerator struct {
}

func init() {
	generator = cardGenerator{}
	rwdGenerator = rewardGenerator{}
}
func (g *cardGenerator) GenerateCard(ms *MainScene) Card {
	baseCard := NewBaseCard([]CardDecorator{})
	p := rand.Int() % 6
	switch p {
	case 0:
		baseCard.AddDecorator(NewCoinDecorator())
	case 1:
		var decor CardDecorator
		if ms.State.Coin <= 9 {
			decor = NewGoblinDecor()
			// baseCard.AddDecorator()
		} else {
			decor = NewHopGoblinDecor()
			// baseCard.AddDecorator()
			h := rand.Int() % 10
			if h <= 3 {
				j := []byte{1, 2, 4, 8}
				decor = NewWeaknessDecorator(decor, j[rand.Int()%len(j)])
			}
		}

		// decor = NewWeaknessDecorator(decor, DIRECTION_UP)
		baseCard.AddDecorator(decor)
	case 2:
		baseCard.AddDecorator(NewSkeletonDecor())
	case 3:
		baseCard.AddDecorator(NewSwordDecorator())
	case 4:
		baseCard.AddDecorator(NewChestDecorator())
	case 5:
		decorators := []CardDecorator{NewBombDecorator(), NewSpikeTrapDecorator()}
		decor := decorators[rand.Int()%len(decorators)]
		baseCard.AddDecorator(decor)
	}
	return baseCard
}

type rewardGenerator struct {
	seed int
}

func (r *rewardGenerator) GenerateReward(tierLevel int) CardDecorator {
	switch tierLevel {
	case 0:
		id := rand.Int() % 3
		switch id {
		case 0:
			return NewCoinDecorator()
		case 1:
			return NewLightPotionDecorator()
		case 2:
			return NewMeat()
		}
	case 1:
		id := rand.Int() % 2
		switch id {
		case 0:
			return NewLightPotionDecorator()
		case 1:
			return NewMedPotionDecorator()
		}

	default:
		return NewCoinDecorator()
	}
	return NewCoinDecorator()
}
