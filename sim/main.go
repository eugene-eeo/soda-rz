package main

import "bufio"
import "os"
import "encoding/json"

type outputRow struct {
	Level       int     `json:"l"`
	Damage      int     `json:"d"`
	Probability float64 `json:"p"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	conf, err := readConfig()
	check(err)
	params := paramsFromConfig(conf)
	data := run(params)
	stdout := bufio.NewWriter(os.Stdout)
	encoder := json.NewEncoder(stdout)
	N := float64(params.samples)
	for level, counter := range data {
		for damage, count := range counter {
			check(encoder.Encode(outputRow{
				Level:       level,
				Damage:      damage,
				Probability: float64(count) / N,
			}))
		}
	}
	check(stdout.Flush())
}
