package crawler

import (
	"fmt"
	"strings"

	"github.com/WakeupTsai/boxoffice-tw/entity"
	"github.com/gocolly/colly"
	"github.com/manifoldco/promptui"
)

func GetReportList() (string, error) {
	reports := []entity.Report{}

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: tfi.org.tw
		colly.AllowedDomains("www.tfi.org.tw"),
	)

	// On every a element which has tr attribute call callback
	c.OnHTML("tr", func(e *colly.HTMLElement) {

		ch := e.DOM.Children()

		if link, _ := ch.Eq(4).Children().Eq(0).Attr("href"); link != "" {
			report := entity.Report{}
			report.Title = ch.Eq(1).Text()

			attr, _ := ch.Eq(4).Children().Eq(0).Attr("onclick")
			report.FileName = strings.Split(attr, "'")[1]

			reports = append(reports, report)
		}
	})

	c.Visit("https://www.tfi.org.tw/BoxOfficeBulletin/weekly")

	templates := &promptui.SelectTemplates{
		Active:   "{{ \"▸\" | cyan }} {{ .Title | cyan }}",
		Inactive: "  {{ .Title }}",
		Selected: "{{ \"▸\" | red | cyan }} {{ .Title | red | cyan }}",
	}

	prompt := promptui.Select{
		Label:     "Select Week",
		Templates: templates,
		Items:     reports,
		Size:      20,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}
	return result, nil
}
