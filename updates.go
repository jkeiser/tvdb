package tvdb

type Updates struct {
	Time     uint64
	Series   []uint64
	Episodes []uint64 `xml:"Episode"`
}
