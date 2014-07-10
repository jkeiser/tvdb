package tvdb

type Episode struct {
	Id                    uint64   `xml:"id"`
	CombinedEpisodeNumber uint64   `xml:"Combined_episodenumber"`
	CombinedSeason        uint64   `xml:"Combined_season"`
	DVDChapter            string `xml:"DVD_chapter"`
	DVDDiscId             string `xml:"DVD_discid"`
	DVDEpisodeNumber      string `xml:"DVD_episodenumber"`
	DVDSeason             string `xml:"DVD_season"`
	Director              string
	EpImgFlag             uint64
	EpisodeName           string
	EpisodeNumber         uint64
	FirstAired            string
	GuestStars            PipeDelimitedString
	ImdbId                string `xml:"IMDB_ID"`
	Language              string
	Overview              string
	ProductionCode        string
	Rating                float32
	RatingCount           uint64
	SeasonNumber          uint64
	Writer                string
	AbsoluteNumber        uint64   `xml:"absolute_number"`
	AirsAfterSeason       uint64   `xml:"airsafter_season"`
	AirsBeforeEpisode     uint64   `xml:"airsbefore_episode"`
	AirsBeforeSeason      uint64   `xml:"airsbefore_season"`
	Filename              string `xml:"filename"`
	LastUpdated           uint64   `xml:"lastupdated"`
	SeasonId              uint64   `xml:"seasonid"`
	SeriesId              uint64   `xml:"seriesid"`
	ThumbAdded            string `xml:"thumb_added"`
	ThumbHeight           uint64   `xml:"thumb_height"`
	ThumbWidth            uint64   `xml:"thumb_width"`
}
