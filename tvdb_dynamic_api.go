package tvdb

import "net/url"
import "strconv"
import "io"
import "log"
import "reflect"

type TvdbDynamicApi struct {
	Getter PathGetter
	ApiKey string
}

//
// API:
//
// http://thetvdb.com/wiki/index.php?title=Programmers_API
//

// http://thetvdb.com/api/GetSeries.php?seriesname=<seriesname>
// http://thetvdb.com/api/GetSeries.php?seriesname=<seriesname>&language=<language>
func (api TvdbDynamicApi) GetSeries(seriesName string, language string) (series []Series, err error) {
	params := "seriesname=" + url.QueryEscape(seriesName)
	if language != "" {
		params += "&language=" + url.QueryEscape(language)
	}
	err = api.getXmlList("GetSeries.php?"+params, "Series", &series)
	return
}

func (api TvdbDynamicApi) getSeriesByRemoteId(key string, value string) (series Series, err error) {
	var serieses []Series
	err = api.getXmlList("GetSeriesByRemoteID.php?"+key+"="+value, "Series", &serieses)
	if len(serieses) == 0 {
		return
	}
	series = serieses[0]
	return
}

func (api TvdbDynamicApi) GetSeriesByIMDB(imdbId string) (series Series, err error) {
	return api.getSeriesByRemoteId("imdb_id", imdbId)
}

func (api TvdbDynamicApi) GetSeriesByZap2ItId(zap2ItId string) (series Series, err error) {
	return api.getSeriesByRemoteId("zap2it", zap2ItId)
}

//<mirrorpath>/api/GetEpisodeByAirDate.php?apikey=<apikey>&seriesid=<seriesid>&airdate=<airdate>&[language=<language>}
func (api TvdbDynamicApi) GetEpisodeByAirDate(seriesId uint64, airDate string, language string) (episode Episode, err error) {
	params := "apikey=" + api.ApiKey + "&seriesid=" + strconv.FormatUint(seriesId, 10) + "&airDate=" + airDate
	if language != "" {
		params += "&language=" + url.QueryEscape(language)
	}

	var episodes []Episode
	err = api.getXmlList("GetEpisodeByAirDate.php?"+params, "Episode", &episodes)
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
func (api TvdbDynamicApi) GetRatingsForUser(accountId string, seriesId string) (ratings []Rating, err error) {
	params := "apikey=" + api.ApiKey + "&accountid=" + accountId
	if seriesId != "" {
		params += "&seriesid=" + url.QueryEscape(seriesId)
	}

	err = api.getXmlList("GetRatingsForUser.php?"+params, "Rating", &ratings)
	return
}

// /api/DynamicUpdates.php?type=none
func (api TvdbDynamicApi) GetServerTime() (time uint64, err error) {
	var updates DynamicUpdates
	err = api.getXml("DynamicUpdates.php?type=none", &updates)
	if err != nil {
		return
	}
	time = updates.Time
	return
}

// /api/DynamicUpdates.php?since=<time>&type=<type>
func (api TvdbDynamicApi) GetAllUpdatesSince(sinceTime uint64) (updates DynamicUpdates, err error) {
	err = api.getXml("DynamicUpdates.php?time="+strconv.FormatUint(sinceTime, 10)+"&type=all", &updates)
	return
}

func (api TvdbDynamicApi) GetSeriesUpdatesSince(sinceTime uint64) (updates DynamicUpdates, err error) {
	err = api.getXml("DynamicUpdates.php?time="+strconv.FormatUint(sinceTime, 10)+"&type=series", &updates)
	return
}

func (api TvdbDynamicApi) GetEpisodeUpdatesSince(sinceTime uint64) (updates DynamicUpdates, err error) {
	err = api.getXml("DynamicUpdates.php?time="+strconv.FormatUint(sinceTime, 10)+"&type=episode", &updates)
	return
}

// /api/User_PreferredLanguage.php?accountid=<accountidentifier>
func (api TvdbDynamicApi) GetUserPreferredLanguage(accountId string) (language Language, err error) {
	var languages []Language
	err = api.getXmlList("User_PreferredLanguage.php?accountid="+accountId, "Language", &languages)
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
func (api TvdbDynamicApi) GetUserFavorites(accountId string) (series []string, err error) {
	err = api.getXmlList("User_Favorites.php?accountid="+accountId, "Series", &series)
	return
}

// /api/User_Favorites.php?accountid=<accountidentifier>&type=add&seriesid=<seriesid>
func (api TvdbDynamicApi) AddUserFavorite(accountId string, seriesId uint64) (err error) {
	return api.noopGet("User_Favorites.php?accountid=" + accountId + "&type=add&seriesid=" + strconv.FormatUint(seriesId, 10))
}

// /api/User_Favorites.php?accountid=<accountidentifier>&type=remove&seriesid=<seriesid>
func (api TvdbDynamicApi) RemoveUserFavorite(accountId string, seriesId uint64) (err error) {
	return api.noopGet("User_Favorites.php?accountid=" + accountId + "&type=remove&seriesid=" + strconv.FormatUint(seriesId, 10))
}

func (api TvdbDynamicApi) userRating(accountId string, itemType string, itemId uint64, rating uint64) (err error) {
	params := "accountid=" + accountId
	params += "&itemtype=" + itemType
	params += "&itemid=" + strconv.FormatUint(itemId, 10)
	params += "&rating=" + strconv.FormatUint(rating, 10)
	return api.noopGet("User_Rating.php?" + params)
}

// /api/User_Rating.php?accountid=<accountidentifier>&itemtype=episode&itemid=<itemid>&rating=<rating>
func (api TvdbDynamicApi) RateEpisode(accountId string, episodeId uint64, rating uint64) (err error) {
	return api.userRating(accountId, "episode", episodeId, rating)
}

func (api TvdbDynamicApi) UnrateEpisode(accountId string, episodeId uint64) (err error) {
	return api.userRating(accountId, "episode", episodeId, 0)
}

func (api TvdbDynamicApi) RateSeries(accountId string, seriesId uint64, rating uint64) (err error) {
	return api.userRating(accountId, "series", seriesId, rating)
}

func (api TvdbDynamicApi) UnrateSeries(accountId string, seriesId uint64) (err error) {
	return api.userRating(accountId, "series", seriesId, 0)
}

func (api TvdbDynamicApi) FileApi() TvdbFileApi {
	return TvdbFileApi{Getter: RelativeGetter{RelativeTo: api.Getter, RelativePath: api.ApiKey}}
}

//
// Helpers
//
func (api TvdbDynamicApi) getXmlList(relativePath, elementName string, result interface{}) (err error) {
	var reader io.ReadCloser
	reader, err = api.Getter.Get(relativePath)
	if err != nil {
		return
	}
	defer reader.Close()
	err = GetXmlList(reader, elementName, result)
	log.Printf("GetXmlList: %s returned %d results", api.Getter.Path(relativePath), reflect.ValueOf(result).Elem().Len())
	return
}

func (api TvdbDynamicApi) getXml(relativePath string, result interface{}) (err error) {
	var reader io.ReadCloser
	reader, err = api.Getter.Get(relativePath)
	if err != nil {
		return
	}
	defer reader.Close()
	return GetXml(reader, result)
}

func (api TvdbDynamicApi) noopGet(relativePath string) (err error) {
	var reader io.ReadCloser
	reader, err = api.Getter.Get(relativePath)
	if reader != nil {
		reader.Close()
	}
	return
}
