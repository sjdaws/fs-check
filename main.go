package main

import (
	"fmt"
	"log"

	"github.com/sjdaws/fs-check/find"
	"github.com/sjdaws/fs-check/notify"
)

func main() {
	config, err := parseConfig("/config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	notifier := notify.New(config.Notify.Urls)

	if len(config.Paths.Check) <= 0 {
		return
	}

	finder := find.New(config.Debug != "")
	finder.AddAllowPaths(config.Paths.Allow...)
	finder.AddAllowTerms(config.Terms.Allow...)
	finder.AddBlockPaths(config.Paths.Block...)
	finder.AddBlockTerms(config.Terms.Block...)
	finder.SetType(config.Type)

	found := make([]string, 0)
	for _, path := range config.Paths.Check {
		paths, err := finder.Check(path)
		if err != nil {
			notifier.Message(err.Error())
			log.Fatal(err)
		}

		found = append(found, paths...)
	}

	if len(found) > 0 {
		log.Println("-----------------------------------")
		log.Printf("Found %d paths that violate checks:\n", len(found))
		log.Println("-----------------------------------")

		for _, path := range found {
			log.Println(path)
		}

		notifier.Message(fmt.Sprintf("Found %d paths that violate checks.", len(found)))
	}
}
