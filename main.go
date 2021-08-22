package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
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
	fmt.Println("   [Block Definition]")
	fmt.Println("   Source: ", block.Source)
	fmt.Println("   Type  : ", block.Type)
	fmt.Println("   Rows  :")
	for i, row := range block.Rows {
		fmt.Println("  ", i, ":", row)
	}
	fmt.Println("   Fields:")
	for i, field := range block.Fields {
		fmt.Println("  ", i, ":(", field.Name, ":", field.Value, ")")
	}
}

func main() {
	var blocks []Block
	var keyValueStore [][]string

	var inputFilePath string
	var outputFilePath string
	var keyValuePath string
	var verboseOutput bool

	// flags declaration using flag package
	flag.StringVar(&inputFilePath, "i", "source.txt", "path of the input file. Default is source.txt")
	flag.StringVar(&outputFilePath, "o", "output.txt", "path of the output file. Default is output.txt")
	flag.StringVar(&keyValuePath, "k", "keyvalue.csv", "path of the key-value store file as comma seperated file with two columns. Default is keyvalue.csv")
	flag.BoolVar(&verboseOutput, "v", false, "activate verbose output to console")

	c := color.New(color.FgWhite).Add(color.Bold)
	flag.Usage = func() {
		fmt.Println("cst - the command shell template parser - ALPHA version 0.001\n") // redundant newline ok
		fmt.Println("   cst -i <input filepath> -o <output filepath> [-k <key-value file path>] [-v]\n")
		c.Println("DESCRIPTION")
		fmt.Println("cst is a command line template parser that uses comma seperated files (csv) as it's data source.")
		fmt.Println("It is being used to transform data, present in csv files, into structured text files,")
		fmt.Println("such as HTML, JSON, XML or yaml files or any other structured file format.")
		c.Println("ARGUMENTS")
		fmt.Println("The program needs an input and an output file path.")
		fmt.Println("All needed data file paths are being defined in the template file itself in so called block definitions")
		fmt.Println("The input file defines the template for the output file.")
		fmt.Println("All parsed template definitions are processed and written to the output file.\n")
		fmt.Println("-i specifies the path of the input file")
		fmt.Println("-o specifies the path of the output file")
		c.Println("OPTIONS")
		fmt.Println("-k if specified, cst uses the specified file path as the key value store file")
		fmt.Println("   The key value store file is a csv file with two columns, where the first column")
		fmt.Println("   defines the key and the second the value. Everything is a string.")
		fmt.Println("   if -k is not specified, cst uses the definition in the source file. It has the following format:")
		fmt.Println("   (store;source:<filepath>)")
		fmt.Println("   This needs to be defined in its own text line with opening and closing brackets and the keyword 'store'")
		fmt.Println("-v Activate verbose output to the console.")
		fmt.Println("   This prints out all actions of the parser and the templating engine so that you can debug your template.")
	}
	flag.Parse()

	// TODO: make all variables available through command arguments
	//       verbosity
	//		 save block definitions as JSON
	//       read block definition from definition file or JSON file
	//       Optional: token delimiter like '{{}}'

	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	c.Println("[Parser] Start on input file: ", inputFilePath)

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
					panic(1)
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
						panic(1)
					}
				}
				blocks = append(blocks, block)
				currentBlock = &blocks[len(blocks)-1]
			} else {
				// error in defining a block header, report error and panic
				fmt.Println("Error in defining a blocks header. Row[", rowIndex, "]:", row)
				panic(1)
			}
		} else if strings.Compare(row, "(block-end)") == 0 { // found block-end
			if currentBlock != nil {
				currentBlock.completed = true
				PrintBlock(currentBlock)
				currentBlock = nil
			} else {
				fmt.Println("Error (block-end) without (block-start). Row[", rowIndex, "]:", row)
				panic(1)
			}
		} else if strings.HasPrefix(row, "(store;") {
			if strings.HasSuffix(row, ")") {
				if currentBlock != nil && currentBlock.completed == false {
					fmt.Println("Error (store) definition started without closing current block definition. Row[", rowIndex, "]:", row)
					panic(1)
				}
				// key-value store definition found
				block := Block{}
				block.Header = row
				block.Type = "store"

				strHeader := strings.TrimPrefix(row, "(store;") // Todo: DUPLICATE CODE - MAKE FUNCTION!!
				strHeader = strings.TrimSuffix(strHeader, ")")
				fields := strings.Split(strHeader, ";")
				for i, strField := range fields { // range over name:value definitions and add them to the structs Fields array
					fieldTokens := strings.Split(strField, ":")
					if len(fieldTokens) == 2 {
						field := Field{Name: fieldTokens[0], Value: fieldTokens[1]}
						if field.Name == "source" {
							block.Source = field.Value
						}
						block.Fields = append(block.Fields, field)
					} else {
						fmt.Println("Error in block header definition. Cannot read field definition, field index [", i, "] Row[", rowIndex, "]:", row)
						panic(1)
					}
				}
				blocks = append(blocks, block)
				PrintBlock(&block)
				currentBlock = nil
				// finished with store block
			} else {
				fmt.Println("Error (store) definition started but not closed. Closing bracket is needed! Row[", rowIndex, "]:", row)
				panic(1)
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
	file.Close()
	c.Println("[Parser] Done\n")

	// TODO: Save all blocks in a JSON file
	//       Add the possibility to read back that JSON so that there is an option to create blocks from a JSON definition

	// file, _ := json.MarshalIndent(blocks, "", " ")
	// _ = ioutil.WriteFile("test.json", file, 0644)

	// Now loop over the blocks and work them out
	// write into a buffer that then is going to be appended to an output file

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatal(err)
	}

	c.Println("[Template Engine] Start on output file: ", outputFilePath)

	for i := 0; i < len(blocks); i++ {
		// verbose output
		fmt.Print("Parsing block #", i, " type=", blocks[i].Type, " rows=", blocks[i].Rows)

		if blocks[i].Type == "string" {
			// 1.) basic string from row in source file
			outputStr := blocks[i].Rows[0] + "\n"

			// 2.) Key-Value Templating
			//     {{key}}
			for _, key := range keyValueStore {
				// really dumb method now to loop over all keys but it will do for the first implementation
				// Todo: Performance enhancement here!!
				token := "{{" + key[0] + "}}"
				outputStr = strings.Replace(outputStr, token, key[1], -1) // replace {index} with the column nr=index
			}

			outputFile.WriteString(outputStr)
			fmt.Println(" -> written to output file.")
		} else if blocks[i].Type == "store" {
			// key-store definition
			// read the csv file and use it as a key value store for all tokens
			if blocks[i].Source != "" {
				sourceFile, err := os.Open(blocks[i].Source)
				if err != nil {
					fmt.Println("Error while opening the source file: ", err)
					panic(1)
				}
				reader := csv.NewReader(sourceFile)
				reader.Comma = ','
				reader.Comment = '#'

				keyValueStore, err = reader.ReadAll()
				if err != nil {
					fmt.Println("Error: Could not read the csv file. Description: ", err)
					panic(1)
				}
				// key value store is present in the variable keyValueStore

				// verbose output
				fmt.Println(" -> Key-Value store successfully read: ", keyValueStore)
			} else {
				fmt.Println("Error: Cannot read key value store fields as the source definition is not set. Check the (store ..) definition in the defintiion file..")
				panic(1)
			}

		} else if blocks[i].Type == "foreach" {
			// foreach block definition
			// 1.) open source file
			fmt.Println("")

			block := &blocks[i]
			if block.Source != "" {
				sourceFile, err := os.Open(block.Source)
				if err != nil {
					fmt.Println("Error while opening the source file: ", err)
					panic(1)
				}
				reader := csv.NewReader(sourceFile)
				reader.Comma = ','
				reader.Comment = '#'

				records, _ := reader.ReadAll()
				//fmt.Println("Records in file ", block.Source, ": ", records)

				for recordIndex, record := range records {
					// per record ...
					// loop over all Rows of the current definition block
					for rowIndex, blockRow := range block.Rows {
						// 1.) column-index Templating
						//     {column-index}

						// replace in Row string all occurrences of columns in the current csv record
						outputStr := blockRow
						for indexColumn, column := range record {
							// for every column, replace the {column-index} with its content in the current record
							token := "{" + strconv.Itoa(indexColumn) + "}"
							outputStr = strings.Replace(outputStr, token, column, -1) // replace {index} with the column nr=index
							if indexColumn == (len(record) - 1) {
								// verbose output
								fmt.Println("   Replacing for record #", recordIndex, " with record='", record, "' in row #", rowIndex, "='"+blockRow, " output='", outputStr, "'")
								// last column, add newline character
								outputStr = outputStr + "\n"
							}
						}

						// 2.) Key-Value Templating
						//     {{key}}
						for _, key := range keyValueStore {
							// really dumb method, currently, to loop over all keys but it will do for the first implementation
							// Todo: Performance enhancement here!!
							token := "{{" + key[0] + "}}"
							outputStr = strings.Replace(outputStr, token, key[1], -1) // replace {index} with the column nr=index
						}

						// write the Row to the output file
						outputFile.WriteString(outputStr)
					}
				}
				fmt.Println("   -> written to output file.")
				sourceFile.Close()
			} else {
				fmt.Println("Error: Foreach block does not have a source definition. Please add the source of the csv file in the \"source\" field. Definition:", block.Fields)
				panic(1)
			}
		}
	}
}
