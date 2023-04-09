.DEFAULT_GOAL := f

# Define variables
SERVICE_CLICKHOUSE = signoz-clickhouse
SERVICE_SAMPLEAPP = sampleapp
NAMESPACE_OBSERVABILITY = observability
NAMESPACE_QUALITY = quality

.PHONY: f start stop config create run repo

f:
	kubectl port-forward service/$(SERVICE_CLICKHOUSE) 9000:9000 -n $(NAMESPACE_OBSERVABILITY) &
	kubectl port-forward service/$(SERVICE_SAMPLEAPP) 8090:8090 -n $(NAMESPACE_QUALITY) &

stop:
	lsof -i :9000 -i :8090 | awk 'NR!=1 {print $2}' | xargs kill

start:
	go run ./server

config:
	qt config --endpoint "http://localhost:8080"

create:
	qt test create -d ./cli/yaml/create.yaml

run:
	qt test run -d ./cli/yaml/run.yaml -i 1

repo:
	qt repo --name myrepo --url https://github.com/myuser/myrepo --auth-type oauth2 --token kajhkjashaksjdhasjkhaskjdhkasjhd@github.com/vijeyash1 --dsn http://localhost:9000?username=admin\&password=admin
list:
	qt repo list 
update:
	qt repo update --reponame myrepo

.PHONY: stop
