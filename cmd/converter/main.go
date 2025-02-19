package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/xamma/yck/internal/comparer"
)

func main() {
	sourceFile := flag.String("source", "", "Path to source YAML file")
	compareFile := flag.String("compare", "", "Path to compare YAML file")
	valueMismatch := flag.Bool("value-mismatch", false, "Show value mismatches")

	flag.Parse()

	if *sourceFile == "" || *compareFile == "" {
		fmt.Println("Usage: yck -source=<source.yaml> -compare=<compare.yaml> [-value-mismatch]")
		os.Exit(1)
	}

	sourceData, err := comparer.LoadYAML(*sourceFile)
	if err != nil {
		log.Fatalf("Error loading source file: %v", err)
	}
	compareData, err := comparer.LoadYAML(*compareFile)
	if err != nil {
		log.Fatalf("Error loading compare file: %v", err)
	}

	comparer := comparer.YamlKeyComparer{
		SourceFilePath:     *sourceFile,
		CompareFilePath:    *compareFile,
		ValueMismatchEnabled: *valueMismatch,
		SourceData:         sourceData,
		CompareData:        compareData,
	}

	comparer.CompareMaps(comparer.SourceData, comparer.CompareData, "")
	// fmt.Println("\n--- Parsed Source YAML ---")
	// comparer.PrintYAMLDebug(sourceData, "")

	// fmt.Println("\n--- Parsed Compare YAML ---")
	// comparer.PrintYAMLDebug(compareData, "")

}
