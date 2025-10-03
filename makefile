HASH := $(shell git rev-parse --short HEAD)

build: git
	go build .
run: build
	./go-singbox

pkg: run
	zip Config-$(HASH).zip config.json config/list/* config/sets/* config/config.yaml 1>/dev/null

git:
	git status
	git add .


clean:
	rm Config-*.zip go-singbox cache.db config/proxy-*.yaml > /dev/null 2>&1

# command > /dev/null ; === command 1>/dev/null
# command > /dev/null 2>&1; === command 1>/dev/null 2>&1
