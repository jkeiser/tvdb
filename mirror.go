package tvdb

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
