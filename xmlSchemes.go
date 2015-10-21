package webdavclnt

import "time"

type Multistatus struct {
	//XMLName   xml.Name    `xml:"DAV: multistatus"`
	Responses []Response `xml:"DAV: response"`
}

type Response struct {
	//XMLName xml.Name `xml:"DAV: response"`
	Href     string   `xml:"href"`
	Propstat Propstat `xml:"propstat"`
}

type Propstat struct {
	Prop Prop `xml:"prop"`
}

type Prop struct {
	Getlastmodified string `xml:"getlastmodified"`
}

type ByTime []Response

func (a ByTime) Len() int      { return len(a) }
func (a ByTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool {

	layout := "Mon, 2 Jan 2006 15:04:05 MST"

	timeA, _ := time.Parse(layout, a[i].Propstat.Prop.Getlastmodified)
	timeB, _ := time.Parse(layout, a[j].Propstat.Prop.Getlastmodified)

	return timeA.Before(timeB)

}
