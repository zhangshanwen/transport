project = home
mod = develop

NOW = $(shell date  '+%Y-%m-%d %H:%M:%S')

### build flag
project_flag = 'main.module=$(1)'
release_flag = 'main.mod=$(2)'

### linux env
set_linux_env = CGO_ENABLED=0 GOOS=linux GOARCH=amd64

### go build
define normal_build
     go build  -ldflags "-X  $(project_flag) -X $(release_flag) "   -o bin/$(1)  apps/$(1)/main.go
endef

.PHONY: check dist build run all check

all: build

check: test all build clean fmt todo legacy


build:
	$(call normal_build,admin)
	$(call normal_build,home)
linux:
	`$(set_linux_env) $(normal_build,admin,release)`
	`$(set_linux_env) $(normal_build,home,release)`

run:build
	./bin/$(project)

clean:
	find . -name "*.DS_Store" -type f -delete
	rm -rf bin

test:
	go test -cover -race ./...


fmt:
	go fmt  ./...

todo:
	grep -rnw "TODO" internal

# Legacy code should be removed by the time of release
legacy:
	grep -rnw "\(LEGACY\|Deprecated\)" internal
