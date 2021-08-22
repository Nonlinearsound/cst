package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("data.csv")
	if err != nil {
		fmt.Println(err)
	}

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.Comment = '#'

	records, _ := reader.ReadAll()

	for _, v := range records {
		fmt.Println(v[0], ":", v[1])
	}

}
