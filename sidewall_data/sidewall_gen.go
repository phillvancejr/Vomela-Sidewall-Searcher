package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

const (
	outFile         = "sidewall.json"
	dataDir         = "sidewall_data/"
	versionFilePath = dataDir + "version.txt"
)

var (
	/*
		filenames containing color matches
		files are in two partitions, a newline separated list of the sidewall colors on top, then the word VOMELA on a line and then the vomela colors matching the sidewall colors
	*/
	targets = map[string]string{
		"Amerimax Metal":                 "amerimax.csv",
		"Foremost Metal":                 "foremost_metal.csv",
		"Foremost Fiberglass":            "foremost_fiber.csv",
		"Crane Fiberglass":               "crane_fiber.csv",
		"Lamiplast / Lamilux Fiberglass": "lami_fiber.csv",
	}
	// this holds the final matches that becomes a json object
	matches = make(map[string]map[string]string)
)

func main() {
	fmt.Println("Generate sidewal.json")

	for brand, file := range targets {
		inFile := dataDir + file
		fmt.Println("********** " + brand + " **********")
		fmt.Printf("opening %s...\n", inFile)
		csvFile, e := os.Open(inFile)

		if e != nil {
			log.Fatalf("Could not read csv file %s because: %v\n", file, e)
		}

		reader := csv.NewReader(csvFile)
		// skip header
		if _, e := reader.Read(); e != nil {
			panic(e)
		}

		fmt.Printf("reading %s...\n", inFile)
		count := 0
		for {
			colorMatch, e := reader.Read()

			if e == io.EOF {
				break
			} else if e != nil {
				log.Fatal("csv reading error: ", e)
			}

			if len(colorMatch) != 2 {
				log.Fatalf("incorrect number of fields on CSV line containing %s\n", colorMatch[0])
			}

			sidewall := colorMatch[0]
			vomela := colorMatch[1]

			if _, exists := matches[vomela]; !exists {
				matches[vomela] = make(map[string]string)
			}

			matches[vomela][brand] = sidewall
			count++
		}

		fmt.Printf("read %d matches from %s\n\n", count, inFile)
	}
	fmt.Printf("total match length: %d\n", len(matches))
	fmt.Println("writing sidewall.json matches")
	if jsonFile, e := os.Create(dataDir + outFile); e != nil {
	} else {
		encoder := json.NewEncoder(jsonFile)
		encoder.Encode(matches)
		exec.Command("cp", dataDir+outFile, outFile).Run()
	}
}
