build:
	GOOS=linux go build main.go
	zip build.zip main
	rm main