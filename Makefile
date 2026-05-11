all:
	go build -o build/server ./cmd/server

clean:
	rm -fv build/*
