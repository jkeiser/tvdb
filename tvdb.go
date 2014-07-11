package tvdb

import "errors"
import "math/rand"

type Tvdb struct {
	DiskCache  DiskCache
	ApiKey     string
	DynamicApi TvdbDynamicApi
	FileApi    TvdbFileApi
}

const TheTvdbUrl = "https://thetvdb.com/api"

func New(url string, apiKey string, mirrorPath string) (tvdb Tvdb, err error) {
	// Set up the initial api and cache
	dynamicApi := TvdbDynamicApi{
		Getter: HttpGetter{Url: url},
		ApiKey: apiKey,
	}
	diskCache := DiskCache{
		Root:           mirrorPath,
		DefaultSource:  dynamicApi.FileApi().Getter,
		DirPermissions: 0755,
		Permissions:    0644,
	}
	fileApi := TvdbFileApi{Getter: diskCache}

	tvdb = Tvdb{
		ApiKey:     apiKey,
		DiskCache:  diskCache,
		DynamicApi: dynamicApi,
		FileApi:    fileApi,
	}
	tvdb.PickMirrors()
	return
}

func (tvdb Tvdb) PickMirrors() (err error) {
	var mirrors []Mirror
	mirrors, err = tvdb.FileApi.GetMirrors()
	if err != nil {
		return
	}
	var xmlMirror HttpGetter
	var zipMirror HttpGetter
	var bannerMirror HttpGetter
	xmlMirror, bannerMirror, zipMirror, err = PickMirrors(mirrors)
	if err != nil {
		return
	}
	tvdb.DiskCache.ExtensionSources = map[string]RelativePathGetter{
		"xml":    TvdbDynamicApi{Getter: xmlMirror, ApiKey: tvdb.ApiKey}.FileApi().Getter,
		"banner": TvdbDynamicApi{Getter: bannerMirror, ApiKey: tvdb.ApiKey}.FileApi().Getter,
		"zip":    TvdbDynamicApi{Getter: zipMirror, ApiKey: tvdb.ApiKey}.FileApi().Getter,
	}
	return
}

func PickMirrors(mirrors []Mirror) (xmlMirror HttpGetter, bannerMirror HttpGetter, zipMirror HttpGetter, err error) {
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
		err = errors.New("No XML mirrors found")
		return
	}
	if len(bannerMirrors) == 0 {
		err = errors.New("No banner mirrors found")
		return
	}
	if len(zipMirrors) == 0 {
		err = errors.New("No ZIP mirrors at found")
		return
	}

	// Pick a random mirror of each type
	xmlMirror = HttpGetter{Url: xmlMirrors[rand.Intn(len(xmlMirrors))].MirrorPath}
	bannerMirror = HttpGetter{Url: bannerMirrors[rand.Intn(len(bannerMirrors))].MirrorPath}
	zipMirror = HttpGetter{Url: zipMirrors[rand.Intn(len(zipMirrors))].MirrorPath}
	return
}
