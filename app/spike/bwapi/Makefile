# Spike helper: goenv may reject go.mod's "go 1.26" before GOTOOLCHAIN runs.
# Prefer: ./run.sh  or  make build / make run

.PHONY: build run tidy

build:
	./run.sh build -o /tmp/bwapi .

run:
	./run.sh run .

tidy:
	./run.sh mod tidy
