Draft version.
Only for tests.

## Basic Usage
### Получения значение свойства getlastmodified

	client := webdavclnt.
			NewClient("connectionString").
			SetLogin("login").
			SetPassword("passw")
	result, err := client.PropFind("", "<prop><getlastmodified/></prop>")
	if err != nil{
		//do something
	}
	
	obj := webdavclnt.Multistatus{}
	err = xml.Unmarshal(result, &obj)
	if err != nil{
		//do something
	}


	