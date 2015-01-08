package main

import (
	"fmt"
	"log"
	"time"

	"bitbucket.org/tebeka/selenium"
)

var (
	capabilities = selenium.Capabilities{
		"browserName": "firefox",
	}
)

type Job struct {
	URL         string                `json:"url"`
	Interval    int                   `json:"interval"`
	LastScraped time.Time             `json:"last_scraped"`
	Collections map[string]Collection `json:"collections"`
	ScrapedData ScrapedElements       `json:"scraped_data"`
}

type Collection struct {
	Group     *string    `json:"group"`
	Selectors []Selector `json:"selectors"`
}

type Selector struct {
	Selector string `json:"selector"`
	Name     string `json:"name"`
}

type ScrapedElements map[string][][]map[string]string

func (self *Job) GetInterval() time.Duration {
	// TODO: Move this parsing into main.go when received job is first
	// unmarshaled
	duration, err := time.ParseDuration(fmt.Sprintf("%ds", self.Interval))

	if err != nil {
		log.Printf("error parsing duration\n")
	}

	return duration
}

func (self *Job) Run() error {
	return self.Scrape()
}

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
			values := make([]map[string]string, 0)

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

				values = append(values, value)
			}

			scraped[name] = append(scraped[name], values)
		}
	}

	self.ScrapedData = scraped
	self.LastScraped = time.Now()

	return nil
}
