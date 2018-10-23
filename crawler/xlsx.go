package crawler

import (
	"io/ioutil"
	"net/http"

	"github.com/WakeupTsai/boxoffice-tw/entity"
	"github.com/WakeupTsai/boxoffice-tw/util"
	"github.com/tealeg/xlsx"
)

func GetXlsx(url string) (*xlsx.File, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return xlsx.OpenBinary(body)
}

func ReadData(xlFile *xlsx.File) entity.Sheet {
	data := entity.Sheet{}

	for _, sheet := range xlFile.Sheets {
		for _, item := range sheet.Rows[0].Cells {
			data.Header = append(data.Header, item.String())
		}

		for _, row := range sheet.Rows[1:] {
			data.Records = append(data.Records, entity.Record{
				SerialNumber:            util.StringToInt(row.Cells[0].String()),
				Country:                 row.Cells[1].String(),
				Title:                   row.Cells[2].String(),
				ReleaseDate:             row.Cells[3].String(),
				DistributionCorporation: row.Cells[4].String(),
				ProductionCompany:       row.Cells[5].String(),
				TheaterCount:            util.StringToInt(row.Cells[6].String()),
				WeekTickets:             util.StringToInt(row.Cells[7].String()),
				WeekBoxOffice:           util.StringToInt(row.Cells[8].String()),
				TotalTickets:            util.StringToInt(row.Cells[9].String()),
				TotalBoxOffice:          util.StringToInt(row.Cells[10].String()),
			})
		}
	}

	return data
}
