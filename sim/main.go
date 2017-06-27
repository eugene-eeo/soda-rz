package main

import "bufio"
import "os"
import "encoding/json"

type row struct {
	Level       int     `json:"l"`
	Damage      int     `json:"d"`
	Probability float64 `json:"p"`
}

func main() {
	conf, err := readConfig()
	if err != nil {
		panic(err)
	}
	params := paramsFromConfig(conf)
	data := run(
		params.party,
		params.samples,
		params.levels,
		params.report_every,
	)
	buffered_stdout := bufio.NewWriter(os.Stdout)
	encoder := json.NewEncoder(buffered_stdout)
	N := float64(params.samples)
	for level, counter := range data {
		for damage, count := range counter {
			err := encoder.Encode(row{
				Level:       level,
				Damage:      damage,
				Probability: float64(count) / N,
			})
			if err != nil {
				panic(err)
			}
		}
	}
	buffered_stdout.Flush()
}
