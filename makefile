build:
	go build .

pkg:
	zip config.zip config.json config/list/* config/sets/* config/config.yaml

run: build pkg
