Draft version.
Only for tests.

## Basic Usage
### Получения значение свойств creationdate и getlastmodified

	client := webdavclnt.
			NewClient("connectionString").
			SetLogin("login").
			SetPassword("passw")

	propsmap, err := client.PropFind("/webdav/test.txt", "creationdate", "getlastmodified")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", propsmap)
