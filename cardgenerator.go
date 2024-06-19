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
func (g *cardGenerator) GenerateTreasure(ms *MainScene) CardDecorator {
	if ms.LastDefeatedMiniBoss == "Pyro-Eyes" {
		return NewDiamondDecorator()
	}
	return NewCoinDecorator()
}
func (g *cardGenerator) GenerateOrg(ms *MainScene) CardDecorator {
	if ms.LastDefeatedMiniBoss == "Pyro-Eyes" {
		kk := rand.Int() % 2
		if kk == 0 {
			org := NewOrgDecor()
			direction := []byte{1, 2, 4, 8}
			org = NewWeaknessDecorator(org, byte(direction[rand.Int()%len(direction)]))
			return org
		} else if kk == 1 {
			return NewSalamandragon()
		}

	}

	if ms.MonstersDefeated < 10 {
		// baseCard.AddDecorator(NewSkeletonDecor())
		return NewSkeletonDecor()
	} else if ms.MonstersDefeated < 20 {
		org := NewOrgDecor()
		direction := rand.Int()%15 + 1
		org = NewWeaknessDecorator(org, byte(direction))
		// baseCard.AddDecorator(org)
		return org
	} else {
		return NewXOrg()

	}
}
func (g *cardGenerator) GenerateCard(ms *MainScene) Card {
	baseCard := NewBaseCard([]CardDecorator{})

	// unique boss monsters
	if ms.Character.GetMaxHP() >= 20 && ms.Character.GetMaxHP() <= 30 {
		if _, ok := ms.GeneratedBoss["Pyro-Eyes"]; !ok {
			pyroEyes := NewPyroEyesDecor()
			baseCard.AddDecorator(pyroEyes)
			ms.GeneratedBoss["Pyro-Eyes"] = true
			return baseCard
		}
	}
	maxHp := ms.Character.GetMaxHP()
	if maxHp >= 30 && maxHp <= 45 {
		_, ok := ms.GeneratedBoss["Brandish-maiden"]
		if ms.LastDefeatedMiniBoss == "Pyro-Eyes" && !ok {
			brandish := NewBrandishMaiden()
			baseCard.AddDecorator(brandish)
			ms.GeneratedBoss["Brandish-maiden"] = true
			return baseCard
		}
	}
	p := rand.Int() % 6
	switch p {
	case 0:
		dd := g.GenerateTreasure(ms)
		baseCard.AddDecorator(dd)
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
		decor := g.GenerateOrg(ms)
		baseCard.AddDecorator(decor)
	case 3:
		baseCard.AddDecorator(NewSwordDecorator())
	case 4:
		baseCard.AddDecorator(NewChestDecorator())
	case 5:
		decorators := []CardDecorator{NewBombDecorator(), NewSpikeTrapDecorator(), NewSpikeTrapDecorator()}
		if ms.Character.GetHP() >= 20 {
			decorators[1] = NewCrimsonTrapDecorator()
			decorators[2] = NewCrimsonTrapDecorator()
		}
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
		id := rand.Int() % 5
		switch id {
		case 0:
			return NewCoinDecorator()
		case 1:
			NewMeat()
			// return NewLightPotionDecorator()
		case 2:
			return NewMeat()
		case 3:
			return NewMeat()
		case 4:
			return NewHpUpDecorator()
		}
	case 1:
		id := rand.Int() % 3
		switch id {
		case 0:
			return NewLightPotionDecorator()
		case 1:
			return NewMedPotionDecorator()
		case 2:
			return NewHpUpDecorator()
		}
	case 2:
		return NewSwordDecorator()
	default:
		return NewCoinDecorator()
	}
	return NewCoinDecorator()
}
