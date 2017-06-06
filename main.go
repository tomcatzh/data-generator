package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/tomcatzh/data-generator/data"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	args := os.Args[1:]
	templateFile := "./template.json"
	if len(args) > 0 {
		templateFile = args[0]
	}

	template, err := data.NewTemplate(templateFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when read template: %v", err)
	}

	var wg sync.WaitGroup

	for f := range template.Iterate() {
		wg.Add(1)

		go func(f data.FileData) {
			defer wg.Done()

			f.Save()
		}(f)
	}

	wg.Wait()
}
