package main

import "math/rand"

type Member interface {
	isRagezerker() bool
	buffAttack(int)
	damage() int
	refresh()
}

type Ragezerker struct{}

func (r *Ragezerker) buff(members []Member) {
	index := rand.Intn(len(members))
	members[index].buffAttack(2)
}

func (r *Ragezerker) isRagezerker() bool {
	return true
}

func (r *Ragezerker) buffAttack(i int) {
}

func (r *Ragezerker) damage() int {
	return 0
}

func (r *Ragezerker) refresh() {
}

type Actor struct {
	atk    int
	b_atk  int
	p_crit float64
	m_crit float64
	m_base float64
}

func newActor(atk int, p_crit, m_crit, m_base float64) *Actor {
	return &Actor{
		atk:    atk,
		b_atk:  atk,
		p_crit: p_crit,
		m_crit: m_crit,
		m_base: m_base,
	}
}

func (a *Actor) isRagezerker() bool {
	return false
}

func (a *Actor) buffAttack(i int) {
	a.atk += i
}

func (a *Actor) damage() int {
	atk := float64(a.atk) * a.m_base
	if rand.Float64() <= a.p_crit {
		return int(atk * a.m_crit)
	}
	return int(atk)
}

func (a *Actor) refresh() {
	a.atk = a.b_atk
}

func sim(party []Member, ragers []*Ragezerker, actors []*Actor) func() int {
	return func() int {
		for _, r := range ragers {
			r.buff(party)
		}
		n := 0
		for _, member := range actors {
			n += member.damage()
		}
		return n
	}
}

func filterParty(party []Member) ([]*Ragezerker, []*Actor) {
	r := []*Ragezerker{}
	a := []*Actor{}
	for _, member := range party {
		if member.isRagezerker() {
			r = append(r, member.(*Ragezerker))
		} else {
			a = append(a, member.(*Actor))
		}
	}
	return r, a
}

func refreshParty(party []Member) {
	for _, member := range party {
		member.refresh()
	}
}

func run(party []Member, samples int, levels int, report_every int) map[int]map[int]int {
	data := map[int]map[int]int{}
	ragers, actors := filterParty(party)

	for i := 0; i < samples; i++ {
		refreshParty(party)
		gen := sim(party, ragers, actors)

		for lvl := 1; lvl <= levels; lvl++ {
			dmg := gen()
			if (lvl % report_every) == 0 {
				counter := data[lvl]
				if counter == nil {
					counter = map[int]int{}
					data[lvl] = counter
				}
				counter[dmg] += 1
			}
		}
	}

	return data
}
