package main

import (
	"fmt"
	"time"

	"bitbucket.org/tebeka/selenium"
)

var (
	capabilities = selenium.Capabilities{
		"browserName": "firefox",
	}
)

type Job struct {
	Id          string
	URL         string
	Interval    int
	LastScraped time.Time
	Collections map[string]Collection

	scrapedData ScrapedElements
}

type Collection struct {
	Group     *string
	Selectors []Selector
}

type Selector struct {
	Selector string
	Name     string
}

type ScrapedElements map[string][]map[string]string

func (self *Job) Scrape() error {
	wd, err := selenium.NewRemote(capabilities, "")

	if err != nil {
		return fmt.Errorf("error starting selenium: %s\n", err)
	}

	defer wd.Quit()

	if err := wd.Get(self.URL); err != nil {
		return fmt.Errorf("error fetching URL: %s\n", err)
	}

	root, err := wd.FindElements(selenium.ByCSSSelector, "html")

	if err != nil {
		return fmt.Errorf("error finding `html` element: %s\n", err)
	}

	scraped := make(ScrapedElements)

	for name, collection := range self.Collections {
		if collection.Group != nil {
			root, err = wd.FindElements(selenium.ByCSSSelector, *collection.Group)

			if err != nil {
				continue
			}
		}

		for _, parent := range root {
			for _, selector := range collection.Selectors {
				el, err := parent.FindElement(selenium.ByCSSSelector, selector.Selector)

				if err != nil {
					continue
				}

				text, err := el.GetAttribute("innerHTML")

				if err != nil {
					continue
				}

				value := make(map[string]string)
				value[selector.Name] = text

				scraped[name] = append(scraped[name], value)
			}
		}
	}

	self.scrapedData = scraped
	self.LastScraped = time.Now()

	return nil
}
