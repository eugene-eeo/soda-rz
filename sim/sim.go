package main

import "math/rand"

type Member interface {
	isRagezerker() bool
	buffAttack(int)
	damage() int
	refresh()
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

func (r *Ragezerker) refresh() {
}

func (r *Ragezerker) copy() Member {
	return r
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

func (a *Actor) copy() Member {
	return newActor(
		a.atk,
		a.p_crit,
		a.m_crit,
		a.m_base,
	)
}

type parameters struct {
	party        []Member
	samples      int
	levels       int
	report_every int
	parallelism  int
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

func copyParty(party []Member) []Member {
	p := make([]Member, len(party))
	for i, member := range party {
		p[i] = member.copy()
	}
	return p
}

func refreshParty(party []Member) {
	for _, member := range party {
		member.refresh()
	}
}

type levelStats struct {
	level  int
	damage int
}

func worker(dst chan<- levelStats, params parameters) {
	party := params.party
	levels := params.levels
	report_every := params.report_every
	ragers, actors := filterParty(party)

	for i := 0; i < params.samples; i++ {
		refreshParty(party)
		gen := sim(party, ragers, actors)

		for lvl := 1; lvl <= levels; lvl++ {
			dmg := gen()
			if (lvl % report_every) == 0 {
				dst <- levelStats{lvl, dmg}
			}
		}
	}
}

func aggregate(samples int, sink <-chan levelStats) map[int]map[int]int {
	data := map[int]map[int]int{}
	for i := 0; i < samples; i++ {
		row := <-sink
		counter := data[row.level]
		if counter == nil {
			counter = map[int]int{}
			data[row.level] = counter
		}
		counter[row.damage] += 1
	}
	return data
}

func run(params parameters) map[int]map[int]int {
	n := params.parallelism
	sink := make(chan levelStats, n)

	sample_per_worker := params.samples / n
	remainder := params.samples % n

	for i := 0; i < n; i++ {
		samples := sample_per_worker
		if i == n-1 {
			samples += remainder
		}
		go worker(sink, parameters{
			party:        copyParty(params.party),
			samples:      sample_per_worker,
			levels:       params.levels,
			report_every: params.report_every,
		})
	}

	return aggregate(params.samples, sink)
}
