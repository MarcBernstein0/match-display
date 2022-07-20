build_binary:
	go build -o ./build/match-display main.go

clean:
	rm ./build/match-display

docker_build:
	docker build -t match-display -f ./container/Dockerfile .

docker_run: 
	docker run --env-file .env -p 8080:8080 match-display


run_container: build docker_build docker_run