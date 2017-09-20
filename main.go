package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/tomcatzh/data-generator/data"
	"github.com/tomcatzh/data-generator/ticket"
)

func main() {
	args := os.Args[1:]
	templateFile := "./templates/cloudfront_log.json"
	if len(args) > 0 {
		templateFile = args[0]
	}

	template, err := data.NewFactoryFile(templateFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when read template: %v\n", err)
		os.Exit(2)
	}

	var wg sync.WaitGroup

	ticket, err := ticket.NewGoTicket(runtime.GOMAXPROCS(0) * 3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when read template: %v", err)
		os.Exit(2)
	}

	for f := range template.Iterate() {
		wg.Add(1)

		ticket.Take()
		go func(f *data.File) {
			defer ticket.Return()
			defer wg.Done()
			f.Save()
		}(f)
	}

	wg.Wait()
}
