run: bin/arion-shots-backend
	@PATH="$(PWD)/bin:$(PATH)" heroku local

bin/arion-shots-backend: main.go
	go build -o bin/arion-shots-backend main.go

clean:
	rm -rf bin
