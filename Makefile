build:
	GOOS=linux GOARCH=amd64 go build -o=./bin/pprof-upload-linuxamd64
	GOOS=darwin GOARCH=amd64 go build -o=./bin/pprof-upload-darwinamd64
	GOOS=windows GOARCH=amd64 go build -o=./bin/pprof-upload-windowsamd64

push:
	gsutil cp bin/* gs://jbd-releases