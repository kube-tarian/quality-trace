.PHONY: clickhouse qt http config install


endpoint := "http://localhost:8080"
s: 
	go run ./server -config ./server/config.yaml.sample
	


c: 
	kubectl port-forward service/signoz-clickhouse 9000\:9000 -n observability


create: 
	qt test create -d ./cli/yaml/create.yaml



run: 
	qt test run -d ./cli/yaml/run.yaml -i 1

config:
	qt config --endpoint $(endpoint)

list:
	qt test list

