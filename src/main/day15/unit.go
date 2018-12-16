package day15

type Unit struct {
	Side byte
	HP int
	PosX int
	PosY int
	turnDone bool
}

func NewUnit(side byte, posX int, posY int) *Unit {
	return &Unit{side, 200, posX, posY, false}
}

func (u * Unit) Activate(g *Grid) bool {
	if u.turnDone {
		return true
	}

	attackTargetX, attackTargetY := g.GetAttackTarget(u.PosX, u.PosY)
	if attackTargetX != -1 {
		g.Attack(attackTargetX, attackTargetY)
		u.turnDone = true
		return true
	}

	foundTarget := g.MoveUnit(u.PosX, u.PosY)
	attackTargetX, attackTargetY = g.GetAttackTarget(u.PosX, u.PosY)
	if attackTargetX != -1 {
		g.Attack(attackTargetX, attackTargetY)
		u.turnDone = true
		return true
	}

	u.turnDone = true
	return foundTarget
}

func (u * Unit) CleanUp() {
	u.turnDone = false
}