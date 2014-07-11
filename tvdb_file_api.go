package tvdb

import "io"
import "log"
import "path"
import "reflect"
import "strconv"

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

func (api TvdbFileApi) GetSeriesById(seriesId uint64, language string) (series Series, err error) {
	if language == "" {
		language = "en"
	}
	err = api.getXml(path.Join("series", strconv.FormatUint(seriesId, 10), language), &series)
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
