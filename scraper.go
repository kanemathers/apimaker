package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/tebeka/selenium"
)

var (
	minimumInterval = 10.0

	capabilities = selenium.Capabilities{
		"browserName": "firefox",
	}
)

type Duration time.Duration

func (self *Duration) UnmarshalJSON(data []byte) error {
	str := string(data)

	// probably a better/safer way to do this
	str = strings.Trim(str, "\"")

	// check if the duration contains a unit. if not, default to seconds
	if _, err := strconv.ParseInt(str, 10, 64); err == nil {
		str = fmt.Sprintf("%ss", str)
	}

	duration, err := time.ParseDuration(str)

	if err != nil {
		return fmt.Errorf("unable to parse duration: %s\n", err)
	}

	if duration.Seconds() < minimumInterval {
		return fmt.Errorf("duration must be a minimum of %0.1f seconds\n", minimumInterval)
	}

	*self = Duration(duration)

	return nil
}

type Job struct {
	URL         string                `json:"url"`
	Interval    Duration              `json:"interval"`
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
	return time.Duration(self.Interval)
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
