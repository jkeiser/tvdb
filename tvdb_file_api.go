package tvdb

import "io"
import "log"
import "path"
import "reflect"
import "fmt"

type TvdbFileApi struct {
	Getter RelativePathGetter
}

//
// Files
//
func (api TvdbFileApi) GetMirrors() (mirrors []Mirror, err error) {
	err = api.getXmlList("mirrors.xml", "Mirror", &mirrors)
	return
}

func (api TvdbFileApi) GetLanguages() (languages []Language, err error) {
	err = api.getXmlList("languages.xml", "Language", &languages)
	return
}

func (api TvdbFileApi) GetSeries(seriesId uint64, language string) (series Series, err error) {
	if language == "" {
		language = "en"
	}
	err = api.getXml(pathJoin("series", seriesId, language+".xml"), &series)
	return
}

func (api TvdbFileApi) GetEpisode(episodeId uint64, language string) (episode Episode, err error) {
	if language == "" {
		language = "en"
	}
	err = api.getXml(pathJoin("episodes", episodeId, language+".xml"), &episode)
	return
}

func (api TvdbFileApi) GetSeriesBanners(seriesId uint64) (banners []Banner, err error) {
	err = api.getXmlList(pathJoin("series", seriesId, "banners.xml"), "Banner", &banners)
	return
}

func (api TvdbFileApi) GetSeriesActors(seriesId uint64) (actors []Actor, err error) {
	err = api.getXmlList(pathJoin("series", seriesId, "actors.xml"), "Actor", &actors)
	return
}

func (api TvdbFileApi) GetAllSeriesEpisodes(seriesId uint64, language string) (series Series, episodes []Episode, err error) {
	if language == "" {
		language = "en"
	}
	type seriesEpisodes struct {
		Series   Series
		Episodes []Episode `xml:"Episode"`
	}

	var result seriesEpisodes
	err = api.getXml(pathJoin("series", seriesId, "all", language+".xml"), &result)
	if err != nil {
		return
	}

	return result.Series, result.Episodes, nil
}

func (api TvdbFileApi) GetSeriesEpisode(seriesId uint64, seasonNumber uint64, episodeNumber uint64, language string) (episode Episode, err error) {
	if language == "" {
		language = "en"
	}
	err = api.getXml(pathJoin("series", seriesId, "default", seasonNumber, episodeNumber, language+".xml"), &episode)
	return
}

func (api TvdbFileApi) GetSeriesEpisodeByDVDOrder(seriesId uint64, seasonNumber uint64, episodeNumber uint64, language string) (episode Episode, err error) {
	if language == "" {
		language = "en"
	}
	err = api.getXml(pathJoin("series", seriesId, "dvd", seasonNumber, episodeNumber, language+".xml"), &episode)
	return
}

func (api TvdbFileApi) GetSeriesEpisodeByAbsoluteOrder(seriesId uint64, episodeNumber uint64, language string) (episode Episode, err error) {
	if language == "" {
		language = "en"
	}
	err = api.getXml(pathJoin("series", seriesId, "absolute", episodeNumber, language+".xml"), &episode)
	return
}

func (api TvdbFileApi) GetAllUpdatesInTheLastDay() (updates DynamicUpdates, err error) {
	err = api.getXml("updates/updates_day.xml", &updates)
	return
}

func (api TvdbFileApi) GetAllUpdatesInTheLastWeek() (updates DynamicUpdates, err error) {
	err = api.getXml("updates/updates_week.xml", &updates)
	return
}

func (api TvdbFileApi) GetAllUpdatesInTheLastMonth() (updates DynamicUpdates, err error) {
	err = api.getXml("updates/updates_month.xml", &updates)
	return
}

func (api TvdbFileApi) GetAllUpdatesEver() (updates DynamicUpdates, err error) {
	err = api.getXml("updates/updates_all.xml", &updates)
	return
}

//
// Helpers
//
func (api TvdbFileApi) getXmlList(relativePath, elementName string, result interface{}) (err error) {
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

func (api TvdbFileApi) getXml(relativePath string, result interface{}) (err error) {
	var reader io.ReadCloser
	reader, err = api.Getter.Get(relativePath)
	if err != nil {
		return
	}
	defer reader.Close()
	return GetXml(reader, result)
}

func pathJoinArray(paths []interface{}) string {
	result := ""
	for _, childPath := range paths {
		subPaths, ok := childPath.([]interface{})
		if ok {
			childPath = pathJoinArray(subPaths)
		}
		if result == "" {
			result = fmt.Sprint(childPath)
		} else {
			result = path.Join(result, fmt.Sprint(childPath))
		}
	}
	return result
}

func pathJoin(paths ...interface{}) string {
	return pathJoinArray(paths)
}
