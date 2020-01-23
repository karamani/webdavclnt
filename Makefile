compile:
	goimports -w ./*.go
	go vet ./*.go
	golint ./*.go
	go install
