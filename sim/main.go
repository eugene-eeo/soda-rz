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
	data := run(params)
	stdout := bufio.NewWriter(os.Stdout)
	encoder := json.NewEncoder(stdout)
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
	stdout.Flush()
}
