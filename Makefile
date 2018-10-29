OLD_SHA:=$(shell shasum -a 256 templates/a_templates-packr.go | cut -d' ' -f1)
NEW_SHA= $(shell shasum -a 256 templates/a_templates-packr.go | cut -d' ' -f1)

all: test install post
test: install
	rm -rf ./tmp
	go test -v -count 1 ./...
	go build -o modelgen ./cli
	docker-compose --no-ansi -f docker-compose.yml up -d --force-recreate
	sleep 5
	./modelgen -c root:@localhost:3307 -d modelgen_tests -p models -o tmp generate
	golint -set_exit_status tmp
	# rm -rf modelgen
	# rm -rf ./tmp
test-ci:
	go test -v -count 1 ./...
	go build -o modelgen ./cli
	docker-compose --no-ansi -f docker-compose.yml up -d --force-recreate
	sleep 30 # annoying, but for ci.
	./modelgen -c root:@localhost:3307 -d modelgen_tests -p models -o tmp generate
	golint -set_exit_status tmp
	rm -rf modelgen
	rm -rf ./tmp
clean:
	docker rm -f modelgen-tests
install:
	go install ./cli
post:
	@if [ "$(NEW_SHA)" != "$(OLD_SHA)" ]; then\
        echo "sha comparison failed on templates/a_templates-packr.go";\
		exit 1;\
    fi
