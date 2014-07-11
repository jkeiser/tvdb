package tvdb

//
// http://thetvdb.com/wiki/index.php?title=API:actors.xml
//
type Actor struct {
	Id        uint64 `xml:"id"`
	Image     string
	Name      string
	Role      string
	SortOrder uint64
}

//
// http://thetvdb.com/wiki/index.php?title=API:banners.xml
//
type Banner struct {
	Id            uint64 `xml:"id"`
	BannerPath    string
	BannerType    string
	BannerType2   string
	Colors        string
	Language      string
	Season        uint64
	Rating        float64
	RatingCount   uint64
	SeriesName    bool
	ThumbnailPath string
	VignettePath  string
}

//
// http://thetvdb.com/wiki/index.php?title=API:Updates
//
type DynamicUpdates struct {
	Time     uint64
	Series   []uint64
	Episodes []uint64 `xml:"Episode"`
}

//
// http://thetvdb.com/wiki/index.php?title=API:Base_Episode_Record
//
type Episode struct {
	Id                    uint64 `xml:"id"`
	CombinedEpisodeNumber uint64 `xml:"Combined_episodenumber"`
	CombinedSeason        uint64 `xml:"Combined_season"`
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
	Rating                float64
	RatingCount           uint64
	SeasonNumber          uint64
	Writer                string
	AbsoluteNumber        uint64 `xml:"absolute_number"`
	AirsAfterSeason       uint64 `xml:"airsafter_season"`
	AirsBeforeEpisode     uint64 `xml:"airsbefore_episode"`
	AirsBeforeSeason      uint64 `xml:"airsbefore_season"`
	Filename              string `xml:"filename"`
	LastUpdated           uint64 `xml:"lastupdated"`
	SeasonId              uint64 `xml:"seasonid"`
	SeriesId              uint64 `xml:"seriesid"`
	ThumbAdded            string `xml:"thumb_added"`
	ThumbHeight           uint64 `xml:"thumb_height"`
	ThumbWidth            uint64 `xml:"thumb_width"`
}

//
// http://thetvdb.com/wiki/index.php?title=API:languages.xml
//
type Language struct {
	Id           uint64 `xml:"id"`
	Abbreviation string `xml:"abbreviation"`
	Name         string `xml:"name"`
}

//
// http://thetvdb.com/wiki/index.php?title=API:mirrors.xml
//
type Mirror struct {
	Id         uint64 `xml:"id"`
	MirrorPath string
	TypeMask   uint64 `xml:"typemask"`
}

const XML_MASK = 1
const BANNER_MASK = 2
const ZIP_MASK = 4

func (mirror Mirror) HoldsXmlFiles() bool {
	return (mirror.TypeMask & XML_MASK) > 0
}

func (mirror Mirror) HoldsBannerFiles() bool {
	return (mirror.TypeMask & BANNER_MASK) > 0
}

func (mirror Mirror) HoldsZipFiles() bool {
	return (mirror.TypeMask & ZIP_MASK) > 0
}

//
// http://thetvdb.com/wiki/index.php?title=API:GetRatingsForUser
//
type Rating struct {
	SeriesId        uint64 `xml:"seriesid"`
	UserRating      uint64
	CommunityRating float64
}

//
// http://thetvdb.com/wiki/index.php?title=API:Base_Series_Record
//
type Series struct {
	SeriesId   uint64 `xml:"seriesid"`
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

//
// http://thetvdb.com/wiki/index.php?title=API:Update_Records
//
type Updates struct {
	Time    uint64 `xml:"time"`
	Updates []Update
}

type Update interface{}

type SeriesUpdate struct {
	Time uint64 `xml:"time"`
}

type EpisodeUpdate struct {
	Time uint64 `xml:"time"`
}

type BannerUpdate struct {
	Time uint64 `xml:"time"`
}
