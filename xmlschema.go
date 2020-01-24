package webdavclnt

import "encoding/xml"

// PropValue contatins properties xml-data.
type PropValue struct {
	XMLName xml.Name `xml:""`
	Value   string   `xml:",chardata"`
}

// Prop contains list of property values.
type Prop struct {
	PropList []PropValue `xml:",any"`
}

// Propstat contains properties struct.
type Propstat struct {
	Prop *Prop `xml:"prop"`
}

// Response contains part of PROPFIND response.
type Response struct {
	Href     string    `xml:"href"`
	Propstat *Propstat `xml:"propstat"`
}

// Multistatus contains PROPFIND response.
type Multistatus struct {
	Responses []Response `xml:"response"`
}
