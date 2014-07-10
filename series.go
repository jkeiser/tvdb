package tvdb

type Series struct {
	SeriesId   uint64   `xml:"seriesid"`
	Language   string `xml:"language"`
	SeriesName string
	AliasNames PipeDelimitedString
	Banner     string `xml:"banner"`
	Overview   string
	FirstAired string
	IMDBId     string `xml:"IMDB_ID"`
	Zap2ItId   string `xml:"zap2it_id"`
	Network    string
}
