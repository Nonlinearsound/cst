package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Field struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Block struct {
	Type      string   `json:"type"`
	Source    string   `json:"source"`
	Rows      []string `json:"rows"`
	Fields    []Field  `json:"fields"`
	Header    string   `json:"header"`
	completed bool     `json:"completed"`
}

func PrintBlock(block *Block) {
	fmt.Println("[Block Definition]")
	fmt.Println(" Source: ", block.Source)
	fmt.Println(" Type  : ", block.Type)
	fmt.Println(" Rows  :")
	for i, row := range block.Rows {
		fmt.Println("  ", i, ":", row)
	}
	fmt.Println(" Fields:")
	for i, field := range block.Fields {
		fmt.Println("  ", i, ":(", field.Name, ":", field.Value, ")")
	}
}

func main() {
	var blocks []Block

	file, err := os.Open("source.tr")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var currentBlock *Block
	var rowIndex int32 = 0
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		fmt.Println("Scanning row index ", rowIndex)
		row := scanner.Text()

		if strings.HasPrefix(row, "(block-start;") {
			if strings.HasSuffix(row, ")") { // found header of a block
				if currentBlock != nil && currentBlock.completed == false {
					fmt.Println("Error (block-start) without closing current block definition. Row[", rowIndex, "]:", row)
					panic(-1)
				}

				block := Block{}
				block.Header = row

				strHeader := strings.TrimPrefix(row, "(block-start;")
				strHeader = strings.TrimSuffix(strHeader, ")")
				fields := strings.Split(strHeader, ";")
				for i, strField := range fields { // range over name:value definitions and add them to the structs Fields array
					fieldTokens := strings.Split(strField, ":")
					if len(fieldTokens) == 2 {
						field := Field{Name: fieldTokens[0], Value: fieldTokens[1]}
						if field.Name == "source" {
							block.Source = field.Value
						} else if field.Name == "type" {
							block.Type = field.Value
						}
						block.Fields = append(block.Fields, field)
					} else {
						fmt.Println("Error in block header definition. Cannot read field definition, field index [", i, "] Row[", rowIndex, "]:", row)
						panic(-1)
					}
				}
				blocks = append(blocks, block)
				currentBlock = &blocks[len(blocks)-1]
			} else {
				// error in defining a block header, report error and panic
				fmt.Println("Error in defining a blocks header. Row[", rowIndex, "]:", row)
				panic(-1)
			}
		} else if strings.Compare(row, "(block-end)") == 0 { // found block-end
			if currentBlock != nil {
				currentBlock.completed = true
				fmt.Println("Mew block definition:")
				PrintBlock(currentBlock)
				currentBlock = nil
			} else {
				fmt.Println("Error (block-end) without (block-start). Row[", rowIndex, "]:", row)
				panic(-1)
			}
		} else {
			// anything else is a template definition of the block
			// it will be added to the Rows array of the currentBlock
			if currentBlock == nil {
				block := Block{}
				block.Rows = append(block.Rows, row)
				block.Type = "string"
				block.completed = true
				blocks = append(blocks, block)
				currentBlock = &blocks[len(blocks)-1]
				PrintBlock(currentBlock)
				currentBlock = nil
			} else {
				currentBlock.Rows = append(currentBlock.Rows, row)
			}
		}
		rowIndex++
	}

	// file, _ := json.MarshalIndent(blocks, "", " ")
	// _ = ioutil.WriteFile("test.json", file, 0644)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
