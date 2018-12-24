package game

import (
	"day24/main/army"
	"sort"
)

type Game struct {
	armies [][]*army.Army
	Stalemate bool
}

func (g Game) IsFinished() bool {
	return g.isDead(g.armies[0]) || g.isDead(g.armies[1])
}

func (g *Game) PlayRound() {
	g.Stalemate = true
	targets := g.targetSelection()
	order := g.getAttackOrder(targets)
	for _, a := range order {
		if a.NumUnits > 0 {
			g.Stalemate = targets[a].AttackBy(a) && g.Stalemate
		}
	}
}

func (g Game) targetSelection() map[*army.Army]*army.Army {
	selectionOrder := g.getSelectionOrder()
	return g.selectTargets(selectionOrder)
}
func (g Game) selectTargets(armies []*army.Army) map[*army.Army]*army.Army {
	available := make(map[*army.Army]bool)
	targets := make(map[*army.Army]*army.Army)
	for _, a := range armies {
		available[a] = true
	}

	for _, a := range armies {
		t := g.getTarget(a, available)
		if t != nil {
			targets[a] = t
			available[t] = false
		}
	}
	return targets
}

func (g Game) getSelectionOrder() []*army.Army {
	order := append(g.StillAlive(g.armies[0]), g.StillAlive(g.armies[1])...)
	sort.Sort(army.ByEffectivePower(order))
	return order
}

func (g Game) isDead(armies []*army.Army) bool {
	for _, a := range armies {
		if a.NumUnits > 0 {
			return false
		}
	}
	return true
}
func (g Game) StillAlive(armies []*army.Army) []*army.Army {
	res := make([]*army.Army, 0)
	for _, a := range armies {
		if a.NumUnits > 0 {
			res = append(res, a)
		}
	}
	return res
}

func (g Game) getTarget(a *army.Army, available map[*army.Army]bool) *army.Army {
	maxDamage := -1
	effectivePower := -1
	initiative := -1
	var target *army.Army
	for _, enemy := range g.armies[1-a.Side] {
		if available[enemy] && enemy.NumUnits > 0 && a.NumUnits > 0 {
			damage := a.Damage(enemy)
			if damage > maxDamage || (damage == maxDamage && enemy.EffectivePower() > effectivePower) || (damage == maxDamage && enemy.EffectivePower() == effectivePower && enemy.Initiative > initiative) {
				target = enemy
				maxDamage = damage
				effectivePower = enemy.EffectivePower()
				initiative = enemy.Initiative
			}
		}
	}
	if maxDamage == 0 {
		return nil
	}
	return target
}
func (g *Game) getAttackOrder(armies map[*army.Army]*army.Army) []*army.Army {
	order := make([]*army.Army, 0)
	for k := range armies {
		if k.NumUnits > 0 {
			order = append(order, k)
		}
	}
	sort.Sort(army.ByAttackOrder(order))
	return order
}

func New(lines []string, boost int) *Game {
	armies := make([][]*army.Army, 2)
	armies[0] = make([]*army.Army, 0)
	armies[1] = make([]*army.Army, 0)

	side := -1
	number := 1
	for _, line := range lines {
		if line == "Immune System:" {
			side = 0
			number = 1
		} else if line == "Infection:" {
			side = 1
			number = 1
		} else if len(line) > 1 {
			armies[side] = append(armies[side], army.New(line, side, number))
			number++
		}
	}


	g := Game{armies, false}
	for _, a := range armies[0] {
		a.Attack = a.Attack + boost
	}
	return &g
}

func (g Game) WinnerScore() int {
	if g.isDead(g.armies[0]) {
		return g.numUnits(g.armies[1])
	} else {
		return g.numUnits(g.armies[0])
	}
}
func (g Game) numUnits(armies []*army.Army) int {
	score := 0
	for _, a := range armies {
		score += a.NumUnits
	}
	return score
}
func (g *Game) ImmuneWins() bool {
	return g.isDead(g.armies[1])
}
