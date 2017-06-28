package main

import "os"
import "encoding/json"

type memberDef struct {
	Atk   int     `json:"atk"`
	Pcrit float64 `json:"p_crit"`
	Mcrit float64 `json:"m_crit"`
	Mbase float64 `json:"m_base"`
}

type config struct {
	Party       []memberDef `json:"party"`
	Samples     int         `json:"samples"`
	Levels      int         `json:"levels"`
	ReportEvery int         `json:"report_every"`
	Parallelism int         `json:"parallelism"`
}

func readConfig() (*config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	conf := new(config)
	dec := json.NewDecoder(file)
	err = dec.Decode(conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func paramsFromConfig(conf *config) parameters {
	party := []Member{}
	for _, member := range conf.Party {
		party = append(party, newActor(
			member.Atk,
			member.Pcrit,
			member.Mcrit,
			member.Mbase,
		))
	}
	if 5-len(party) > 0 {
		r := &Ragezerker{}
		for i := 0; i < 5-len(party); i++ {
			party = append(party, r)
		}
	}
	return parameters{
		party:        party,
		samples:      conf.Samples,
		levels:       conf.Levels,
		report_every: conf.ReportEvery,
		parallelism:  max(conf.Parallelism, 1),
	}
}
