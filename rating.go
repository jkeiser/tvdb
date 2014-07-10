package tvdb

type Rating struct {
	SeriesId        uint64 `xml:"seriesid"`
	UserRating      uint64
	CommunityRating float32
}
