package entity

type Record struct {
	SerialNumber            int
	Country                 string
	Title                   string
	ReleaseDate             string
	DistributionCorporation string
	ProductionCompany       string
	TheaterCount            int
	WeekTickets             int
	WeekBoxOffice           int
	TotalTickets            int
	TotalBoxOffice          int
}

type Sheet struct {
	Header  []string
	Records []Record
}

type SelectItem struct {
	Record
	Output string
}

type Report struct {
	Title    string
	FileName string
}
