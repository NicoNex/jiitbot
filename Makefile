PROGNAME = jiitbot

$(PROGNAME):
	CGO_ENABLED=0 go build -o $(PROGNAME)

armv8:
	CGO_ENABLED=0 GOARCH=arm64 go build -o $(PROGNAME)

armv7:
	CGO_ENABLED=0 GOARM=7 GOARCH=arm go build -o $(PROGNAME)

all: $(PROGNAME)

.PHONY: all $(PROGNAME) clean

small:
	CGO_ENABLED=0 go build -o $(PROGNAME) -ldflags="-s -w"

tiny:
	CGO_ENABLED=0 go build -o $(PROGNAME) -ldflags="-s -w"
	upx --brute $(PROGNAME)

clean:
	rm $(PROGNAME)