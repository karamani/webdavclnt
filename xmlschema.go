package webdavclnt

import "encoding/xml"

type PropValue struct {
	XMLName	xml.Name	 `xml:""`
	Value	string	 `xml:",chardata"`
}

type Prop struct {
	PropList []PropValue `xml:",any"`
}

type Propstat struct {
	Prop *Prop `xml:"prop"`
}

type Response struct {
	Href     string   `xml:"href"`
	Propstat *Propstat `xml:"propstat"`
}

type Multistatus struct {
	Responses []Response `xml:"response"`
}
