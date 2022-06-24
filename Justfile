build:
	go get -u
	go mod tidy
	go build -o init .
push:
	git add .
	git cz	
	git push
