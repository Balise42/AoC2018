package army

import (
	"regexp"
	"strconv"
	"strings"
)

type Army struct {
	Side       int // 0 for the immune system, 1 for the infection
	NumUnits   int
	HitPoints  int
	Immunities []string
	Weaknesses []string
	Attack     int
	AttackType string
	Initiative int
	Number     int
}

func (a Army) EffectivePower() int {
	return a.NumUnits * a.Attack
}
func (a *Army) Damage(enemy *Army) int {
	if enemy.isImmune(a) {
		return 0
	} else if enemy.isWeak(a) {
		return a.Attack * a.NumUnits * 2
	} else {
		return a.Attack * a.NumUnits
	}
}

func (a Army) isImmune(attacker *Army) bool {
	return a.containsAttack(a.Immunities, attacker.AttackType)
}

func (a Army) isWeak(attacker *Army) bool {
	return a.containsAttack(a.Weaknesses, attacker.AttackType)
}

func (Army) containsAttack(attacks []string, attack string) bool {
	for _, v := range attacks {
		if v == attack {
			return true
		}
	}
	return false
}
func (a *Army) AttackBy(attacker *Army) bool{
	damage := attacker.Damage(a)
	numHitUnits := damage / a.HitPoints
	if numHitUnits > a.NumUnits {
		numHitUnits = a.NumUnits
		a.NumUnits = 0
	} else {
		a.NumUnits -= numHitUnits
	}
	return numHitUnits == 0
	//fmt.Println("Side ", attacker.Side, " group ", attacker.Number, " attacks side ", a.Side, " group ", a.Number, " and kills ", numHitUnits, " units")
}

func (a *Army) GetNilArmy() *Army {
	return nil
}

var armyRegex = regexp.MustCompile(`(\d+) units each with (\d+) hit points (\(.*\) )?with an attack that does (\d+) ([a-z]+) damage at initiative (\d+)`)

func New(line string, side int, number int) *Army {
	toks := armyRegex.FindStringSubmatch(line)
	numUnits, _ := strconv.Atoi(toks[1])
	hitPoints, _ := strconv.Atoi(toks[2])
	weakImmune := toks[3]
	attack, _ := strconv.Atoi(toks[4])
	attackType := toks[5]
	initiative, _ := strconv.Atoi(toks[6])

	army := Army{
		Side:       side,
		NumUnits:   numUnits,
		HitPoints:  hitPoints,
		Immunities: make([]string, 0),
		Weaknesses: make([]string, 0),
		Attack:     attack,
		AttackType: attackType,
		Initiative: initiative,
		Number:     number,
	}

	if weakImmune != "" {
		weakImmune = weakImmune[1:len(weakImmune)-2]
		toks = strings.Split(weakImmune, "; ")
		for _, tok := range toks {
			if strings.HasPrefix(tok, "immune to") {
				tok = strings.Replace(tok, "immune to ", "", 1)
				army.Immunities = strings.Split(tok, ", ")
			} else {
				tok = strings.Replace(tok, "weak to ", "", 1)
				army.Weaknesses = strings.Split(tok, ", ")
			}
		}
	}

	return &army
}

type ByEffectivePower []*Army

func (s ByEffectivePower) Len() int {
	return len(s)
}

func (s ByEffectivePower) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByEffectivePower) Less(i, j int) bool {
	if s[i].EffectivePower() > s[j].EffectivePower() {
		return true
	} else if s[j].EffectivePower() > s[i].EffectivePower() {
		return false
	}
	return s[i].Initiative > s[j].Initiative
}

type ByAttackOrder []*Army

func (s ByAttackOrder) Len() int {
	return len(s)
}

func (s ByAttackOrder) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByAttackOrder) Less(i, j int) bool {
	return s[i].Initiative > s[j].Initiative
}
