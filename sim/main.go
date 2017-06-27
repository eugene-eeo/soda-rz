package main

import "bufio"
import "os"
import "encoding/json"

type row struct {
	Level  int `json:"level"`
	Damage int `json:"damage"`
	Count  int `json:"count"`
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
	for level, counter := range data {
		for damage, count := range counter {
			err := encoder.Encode(row{
				Level:  level,
				Damage: damage,
				Count:  count,
			})
			if err != nil {
				panic(err)
			}
		}
	}
	buffered_stdout.Flush()
}
