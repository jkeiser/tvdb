package tvdb

import "net/url"
import "errors"
import "math/rand"
import "strconv"
import "io"

type TVDB struct {
	Url          string
	ApiKey       string
	MirrorPath   string
	apiClient    ApiClient
	xmlClient    ApiClient
	zipClient    ApiClient
	bannerClient ApiClient
	diskMirror   DiskMirror
}

const TheTVDBUrl = "https://thetvdb.com/api"

func New(url string, apiKey string, mirrorPath string) (tvdb TVDB, err error) {
	apiClient := ApiClient{Url: url}
	diskMirror := DiskMirror{Path: mirrorPath, DirPermissions: 0755, Permissions: 0644}
	tvdb = TVDB{Url: url, ApiKey: apiKey, MirrorPath: mirrorPath, apiClient: apiClient, diskMirror: diskMirror}
	err = tvdb.PickMirrors()
	if err != nil {
		return
	}
	return
}

//
// API:
//
// http://thetvdb.com/wiki/index.php?title=Programmers_API
//

// http://thetvdb.com/api/GetSeries.php?seriesname=<seriesname>
// http://thetvdb.com/api/GetSeries.php?seriesname=<seriesname>&language=<language>
func (tvdb TVDB) GetSeries(seriesName string, language string) (series []Series, err error) {
	params := "seriesname=" + url.QueryEscape(seriesName)
	if language != "" {
		params += "&language=" + url.QueryEscape(language)
	}
	err = tvdb.apiClient.GetXmlList("GetSeries.php?"+params, "Series", &series)
	return
}

func (tvdb TVDB) getSeriesByRemoteId(key string, value string) (series Series, err error) {
	var serieses []Series
	err = tvdb.apiClient.GetXmlList("GetSeriesByRemoteID.php?"+key+"="+value, "Series", &serieses)
	if len(serieses) == 0 {
		return
	}
	series = serieses[0]
	return
}

func (tvdb TVDB) GetSeriesByIMDB(imdbId string) (series Series, err error) {
	return tvdb.getSeriesByRemoteId("imdb_id", imdbId)
}

func (tvdb TVDB) GetSeriesByZap2ItId(zap2ItId string) (series Series, err error) {
	return tvdb.getSeriesByRemoteId("zap2it", zap2ItId)
}

//<mirrorpath>/api/GetEpisodeByAirDate.php?apikey=<apikey>&seriesid=<seriesid>&airdate=<airdate>&[language=<language>}
func (tvdb TVDB) GetEpisodeByAirDate(seriesId uint64, airDate string, language string) (episode Episode, err error) {
	params := "apikey=" + tvdb.ApiKey + "&seriesid=" + strconv.FormatUint(seriesId, 10) + "&airDate=" + airDate
	if language != "" {
		params += "&language=" + url.QueryEscape(language)
	}

	var episodes []Episode
	err = tvdb.apiClient.GetXmlList("GetEpisodeByAirDate.php?"+params, "Episode", &episodes)
	if err != nil {
		return
	}
	if len(episodes) == 0 {
		return
	}
	episode = episodes[0]
	return
}

// /api/GetRatingsForUser.php?apikey=<apikey>&accountid=<accountidentifier>[&seriesid=<seriesid>]
func (tvdb TVDB) GetRatingsForUser(accountId string, seriesId string) (ratings []Rating, err error) {
	params := "apikey=" + tvdb.ApiKey + "&accountid=" + accountId
	if seriesId != "" {
		params += "&seriesid=" + url.QueryEscape(seriesId)
	}

	err = tvdb.apiClient.GetXmlList("GetRatingsForUser.php?"+params, "Rating", &ratings)
	return
}

// /api/Updates.php?type=none
func (tvdb TVDB) GetServerTime() (time uint64, err error) {
	var updates Updates
	err = tvdb.apiClient.GetXml("Updates.php?type=none", &updates)
	if err != nil {
		return
	}
	time = updates.Time
	return
}

// /api/Updates.php?since=<time>&type=<type>
func (tvdb TVDB) GetAllUpdates(sinceTime uint64) (updates Updates, err error) {
	err = tvdb.apiClient.GetXml("Updates.php?time="+strconv.FormatUint(sinceTime, 10)+"&type=all", &updates)
	return
}

func (tvdb TVDB) GetSeriesUpdates(sinceTime uint64) (updates Updates, err error) {
	err = tvdb.apiClient.GetXml("Updates.php?time="+strconv.FormatUint(sinceTime, 10)+"&type=series", &updates)
	return
}

func (tvdb TVDB) GetEpisodeUpdates(sinceTime uint64) (updates Updates, err error) {
	err = tvdb.apiClient.GetXml("Updates.php?time="+strconv.FormatUint(sinceTime, 10)+"&type=episode", &updates)
	return
}

// /api/User_PreferredLanguage.php?accountid=<accountidentifier>
func (tvdb TVDB) GetUserPreferredLanguage(accountId string) (language Language, err error) {
	var languages []Language
	err = tvdb.apiClient.GetXmlList("User_PreferredLanguage.php?accountid="+accountId, "Language", &languages)
	if err != nil {
		return
	}
	if len(languages) == 0 {
		return
	}
	language = languages[0]
	return
}

// /api/User_Favorites.php?accountid=<accountidentifier>
func (tvdb TVDB) GetUserFavorites(accountId string) (series []string, err error) {
	err = tvdb.apiClient.GetXmlList("User_Favorites.php?accountid="+accountId, "Series", &series)
	return
}

// /api/User_Favorites.php?accountid=<accountidentifier>&type=add&seriesid=<seriesid>
func (tvdb TVDB) AddUserFavorite(accountId string, seriesId uint64) (err error) {
	return tvdb.noopGet("User_Favorites.php?accountid=" + accountId + "&type=add&seriesid=" + strconv.FormatUint(seriesId, 10))
}

// /api/User_Favorites.php?accountid=<accountidentifier>&type=remove&seriesid=<seriesid>
func (tvdb TVDB) RemoveUserFavorite(accountId string, seriesId uint64) (err error) {
	return tvdb.noopGet("User_Favorites.php?accountid=" + accountId + "&type=remove&seriesid=" + strconv.FormatUint(seriesId, 10))
}

func (tvdb TVDB) userRating(accountId string, itemType string, itemId uint64, rating uint64) (err error) {
	params := "accountid=" + accountId
	params += "&itemtype=" + itemType
	params += "&itemid=" + strconv.FormatUint(itemId, 10)
	params += "&rating=" + strconv.FormatUint(rating, 10)
	return tvdb.noopGet("User_Rating.php?" + params)
}

// /api/User_Rating.php?accountid=<accountidentifier>&itemtype=episode&itemid=<itemid>&rating=<rating>
func (tvdb TVDB) RateEpisode(accountId string, episodeId uint64, rating uint64) (err error) {
	return tvdb.userRating(accountId, "episode", episodeId, rating)
}

func (tvdb TVDB) UnrateEpisode(accountId string, episodeId uint64) (err error) {
	return tvdb.userRating(accountId, "episode", episodeId, 0)
}

func (tvdb TVDB) RateSeries(accountId string, seriesId uint64, rating uint64) (err error) {
	return tvdb.userRating(accountId, "series", seriesId, rating)
}

func (tvdb TVDB) UnrateSeries(accountId string, seriesId uint64) (err error) {
	return tvdb.userRating(accountId, "series", seriesId, 0)
}

func (tvdb TVDB) GetMirrors() (mirrors []Mirror, err error) {
	err = tvdb.diskMirror.GetXmlList("mirrors.xml", ApiClient{tvdb.Url + "/" + tvdb.ApiKey}, "Mirror", &mirrors)
	return
}

func (tvdb TVDB) PickMirrors() (err error) {
	mirrors, err := tvdb.GetMirrors()
	if err != nil {
		return
	}

	// Gather up all the mirrors of each type
	var xmlMirrors []Mirror
	var bannerMirrors []Mirror
	var zipMirrors []Mirror
	for _, mirror := range mirrors {
		if mirror.HoldsXmlFiles() {
			xmlMirrors = append(xmlMirrors, mirror)
		}
		if mirror.HoldsBannerFiles() {
			bannerMirrors = append(bannerMirrors, mirror)
		}
		if mirror.HoldsZipFiles() {
			zipMirrors = append(zipMirrors, mirror)
		}
	}
	if len(xmlMirrors) == 0 {
		return errors.New("No XML mirrors at " + tvdb.Url + "/" + tvdb.ApiKey + "/mirrors.xml")
	}
	if len(bannerMirrors) == 0 {
		return errors.New("No banner mirrors at " + tvdb.Url + "/" + tvdb.ApiKey + "/mirrors.xml")
	}
	if len(zipMirrors) == 0 {
		return errors.New("No ZIP mirrors at " + tvdb.Url + "/" + tvdb.ApiKey + "/mirrors.xml")
	}

	// Pick a random mirror of each type
	xmlMirror := xmlMirrors[rand.Intn(len(xmlMirrors))]
	bannerMirror := bannerMirrors[rand.Intn(len(bannerMirrors))]
	zipMirror := zipMirrors[rand.Intn(len(zipMirrors))]

	// Return the API servers for each type
	tvdb.xmlClient = ApiClient{Url: xmlMirror.MirrorPath + "/api/" + tvdb.ApiKey}
	tvdb.bannerClient = ApiClient{Url: bannerMirror.MirrorPath + "/api/" + tvdb.ApiKey}
	tvdb.zipClient = ApiClient{Url: zipMirror.MirrorPath + "/api/" + tvdb.ApiKey}
	return
}

func (tvdb TVDB) noopGet(relative string) (err error) {
	var reader io.ReadCloser
	reader, err = tvdb.apiClient.Get(relative)
	if reader != nil {
		reader.Close()
	}
	return
}
