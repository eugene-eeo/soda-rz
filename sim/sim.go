package main

import "math/rand"

type Member interface {
	isRagezerker() bool
	buffAttack(int)
	damage() int
	copy() Member
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

func (r *Ragezerker) copy() Member {
	return r
}

type Actor struct {
	atk    int
	p_crit float64
	m_crit float64
	m_base float64
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

func (a *Actor) copy() Member {
	return &Actor{
		atk:    a.atk,
		p_crit: a.p_crit,
		m_crit: a.m_crit,
		m_base: a.m_base,
	}
}

func sim(party []Member) func() int {
	ragers := []*Ragezerker{}
	for _, member := range party {
		if member.isRagezerker() {
			ragers = append(ragers, member.(*Ragezerker))
		}
	}
	return func() int {
		for _, r := range ragers {
			r.buff(party)
		}
		n := 0
		for _, member := range party {
			n += member.damage()
		}
		return n
	}
}

func copyParty(party []Member) []Member {
	p := []Member{}
	for _, member := range party {
		p = append(p, member.copy())
	}
	return p
}

func run(party []Member, samples int, levels int, report_every int) map[int]map[int]int {
	data := map[int]map[int]int{}
	for i := 0; i < samples; i++ {
		sample := copyParty(party)
		gen := sim(sample)
		for lvl := 1; lvl < levels+1; lvl++ {
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
