package webdavclnt

type Multistatus struct {
	Responses []Response `xml:"response"`
}

type Response struct {
	Href     string   `xml:"href"`
	Propstat Propstat `xml:"propstat"`
}

type Propstat struct {
	Prop Prop `xml:"prop"`
}

type Prop struct {
	Getlastmodified string `xml:"getlastmodified"`
}
