build_binary:
	go build -o ./build/match-display main.go

clean:
	rm ./build/match-display

docker_build:
	docker build --no-cache -t match-display -f ./container/Dockerfile .

docker_run: 
	docker run --rm --env-file .env -p 8080:8080 match-display

run_container: docker_build docker_run
