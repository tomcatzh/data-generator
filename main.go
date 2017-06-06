package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/tomcatzh/data-generator/data"
	"github.com/tomcatzh/data-generator/ticket"
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
		os.Exit(2)
	}

	var wg sync.WaitGroup

	ticket, err := ticket.NewGoTicket(runtime.GOMAXPROCS(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when read template: %v", err)
		os.Exit(2)
	}

	for f := range template.Iterate() {
		wg.Add(1)

		ticket.Take()
		go func(f data.FileData) {
			defer ticket.Return()
			defer wg.Done()
			f.Save()
		}(f)
	}

	wg.Wait()
}
