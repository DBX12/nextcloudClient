SOURCES=./

test:
	cd $(SOURCES) && \
	go test

cover:
	cd $(SOURCES) && \
	go test -coverprofile=cover.out && \
	go tool cover -html=cover.out && \
	rm cover.out

.PHONY: test cover
