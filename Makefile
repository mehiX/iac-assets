.PHONY: install
install:
	mkdir -p dist
	go build -o dist/iac ./main.go

.PHONY: serve
serve: install
	./dist/iac serve