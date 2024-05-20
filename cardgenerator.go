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
func (g *cardGenerator) GenerateCard() Card {
	baseCard := NewBaseCard([]CardDecorator{})
	p := rand.Int() % 4
	switch p {
	case 0:
		baseCard.AddDecorator(NewCoinDecorator())
	case 1:
		baseCard.AddDecorator(NewGoblinDecor())
	case 2:
		baseCard.AddDecorator(NewSkeletonDecor())
	case 3:
		baseCard.AddDecorator(NewSwordDecorator())
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
