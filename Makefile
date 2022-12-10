clean:
	rm -rf build/*

build: clean
	GOOS="linux" go build -o build/light-bot .

up:
	go get