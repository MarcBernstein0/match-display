help: # Show this help.
	@fgrep -h "#" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/#//'

build_binary: # build go binary
	go build -o ./build/match-display main.go

clean: # cleanup build
	rm ./build/match-display

docker_build: # build docker image
	docker build --no-cache -t match-display -f ./container/Dockerfile .

docker_run: # run docker image
	docker run --rm --env-file .env -p 8080:8080 match-display

build_n_run: docker_build docker_run # build container image and run 
	 
