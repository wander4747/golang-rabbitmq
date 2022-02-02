SHELL := /bin/bash
.PHONY: help

help: ## Exibe essa mensagem de ajuda
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

docker: ## Executa os serviços no docker-compose
	docker-compose up -d

hello-sender:
	docker-compose exec go cd/src & go mod tidy && go run hello/main.go sender

hello-receive:
	docker-compose exec go cd/src & go mod tidy && go run hello/main.go receive


fanout-sender:
	docker-compose exec go cd/src & go mod tidy && go run exchange/fanout/main.go sender

fanout-receive:
	docker-compose exec go cd/src & go mod tidy && go run exchange/fanout/main.go receive

