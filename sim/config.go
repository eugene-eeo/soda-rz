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
}

type parameters struct {
	party        []Member
	samples      int
	levels       int
	report_every int
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

func paramsFromConfig(conf *config) parameters {
	party := []Member{}
	for _, member := range conf.Party {
		party = append(party, &Actor{
			atk:    member.Atk,
			m_base: member.Mbase,
			m_crit: member.Mcrit,
			p_crit: member.Pcrit,
		})
	}
	if 5-len(party) > 0 {
		for i := 0; i < 5-len(party); i++ {
			party = append(party, &Ragezerker{})
		}
	}
	return parameters{
		party:        party,
		samples:      conf.Samples,
		levels:       conf.Levels,
		report_every: conf.ReportEvery,
	}
}
