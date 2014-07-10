package tvdb

type Language struct {
	Id           uint64 `xml:"id"`
	Abbreviation string `xml:"abbreviation"`
	Name         string `xml:"name"`
}
