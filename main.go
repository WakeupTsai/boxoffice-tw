package main

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/WakeupTsai/boxoffice-tw/crawler"
	"github.com/WakeupTsai/boxoffice-tw/entity"
	"github.com/WakeupTsai/promptui"
	"github.com/WakeupTsai/tablewriter"
	humanize "github.com/dustin/go-humanize"
)

func main() {
	result, err := crawler.GetReportList()

	re := regexp.MustCompile(" (.*?)\\}")
	match := re.FindStringSubmatch(result)
	fileURL := "https://www.tfi.org.tw" + match[1]

	err = renderInfo(fileURL)
	if err != nil {
		panic(err)
	}
}

func drawTable(sheet entity.Sheet) {
	table := tablewriter.NewWriter(os.Stdout)

	// sort the list
	//sort.Slice(sheet.Records, func(i, j int) bool {
	//	return sheet.Records[i].WeekBoxOffice > sheet.Records[j].WeekBoxOffice
	//})
	//table.SetAutoMergeCells(true)
	// table.SetRowLine(true)

	table.SetHeader([]string{
		sheet.Header[2],
		sheet.Header[3],
		sheet.Header[8],
		sheet.Header[10],
	})

	data := [][]string{}
	for _, record := range sheet.Records {
		data = append(data, []string{
			record.Title,
			record.ReleaseDate,
			humanize.Comma(int64(record.WeekBoxOffice)),
			humanize.Comma(int64(record.TotalBoxOffice))})
	}

	table.AppendBulk(data)
	table.Render() // Send output
}

func renderInfo(url string) error {

	xlFile, err := crawler.GetXlsx(url)
	if err != nil {
		return nil
	}

	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	sheet := crawler.ReadData(xlFile)
	drawTable(sheet)

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	var header string
	var footer string

	movies := []entity.SelectItem{}
	lines := strings.Split(out, "\n")

	for _, line := range lines[:3] {
		header = header + "\n    " + line
	}

	for i, line := range lines[3 : len(lines)-2] {
		var item entity.SelectItem
		if i < len(sheet.Records) {
			item.Record = sheet.Records[i]
		}
		item.Output = line
		movies = append(movies, item)
	}

	for _, line := range lines[len(lines)-2:] {
		footer = footer + "\n    " + line
	}

	templates := &promptui.SelectTemplates{
		Header:   header,
		Footer:   footer,
		Active:   "\U0001f37f {{ .Output | cyan }}",
		Inactive: "  {{ .Output }}",
		Selected: "\U0001f37f {{ .Title | red | cyan }}",
		Details: `    {{ "` + "國別地區：　　　" + `:" | faint }}　{{ .Country }}
    {{ "` + "申請人：　　　　" + `:" | faint }}　{{ .DistributionCorporation }}
    {{ "` + "出品：　　　　　" + `:" | faint }}　{{ .ProductionCompany }}
    {{ "` + "上映院數：　　　" + `:" | faint }}　{{ .TheaterCount }}
    {{ "` + "銷售票數：　　　" + `:" | faint }}　{{ .WeekTickets }}
    {{ "` + "累計銷售票數：　" + `:" | faint }}　{{ .TotalTickets }}`,
	}

	searcher := func(input string, index int) bool {
		movie := movies[index]
		name := movie.Title

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     []string{"Select Movie"},
		Items:     movies,
		Size:      15,
		Templates: templates,
		Searcher:  searcher,
	}

	_, _, err = prompt.Run()

	return nil
}
