package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	
	"mermaid2drawio/internal/mermaid"
	"mermaid2drawio/internal/drawio"
)

func main() {
	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose error output")
	flag.Parse()
	
	if err := run(verbose); err != nil {
		if verbose {
			log.Printf("Error: %v", err)
		}
		os.Exit(1)
	}
}

func run(verbose bool) error {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}
	
	diagram, err := mermaid.ParseDiagram(string(input))
	if err != nil {
		return fmt.Errorf("parsing Mermaid diagram: %w", err)
	}
	
	xmlOutput, err := drawio.GenerateDrawIOXML(diagram)
	if err != nil {
		return fmt.Errorf("generating draw.io XML: %w", err)
	}
	
	fmt.Print(xmlOutput)
	return nil
}
