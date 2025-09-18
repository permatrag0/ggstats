matches-service:
	go run matches/cmd/main/main.go;
metadata-service:
	go run metadata/cmd/main/main.go;
players-service:
	go run players/cmd/main/main.go;

test-matches-service:

test-metadata-service:
	curl -X PUT -d id=1 -d localhost:8081/metadata
test-players-service:
	curl -v localhost:8083/players?id=1

run-consul-dev-server:
	docker run -d -p 8500:8500 -p 8600:8600/udp --name=dev-consul consul:1.15.4 agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0