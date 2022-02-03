SHELL := /bin/bash
.PHONY: help

help: ## Exibe essa mensagem de ajuda
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

docker: ## Executa os serviços no docker-compose
	docker-compose up -d

hello-sender: ## Executa uma fila simples de exemplo - producer
	docker-compose exec go cd/src & go mod tidy && go run hello/main.go sender

hello-receive: ## Executa uma fila simples de exemplo - consumer
	docker-compose exec go cd/src & go mod tidy && go run hello/main.go receive

fanout-sender: ## Executa um producer. Exchange fanout
	docker-compose exec go cd/src & go mod tidy && go run exchange/fanout/main.go sender

fanout-receive: ## Executa um consumer. Exchange fanout
	docker-compose exec go cd/src & go mod tidy && go run exchange/fanout/main.go receive

direct-sender: ## Executa um producer. Exchange direct. Ex: make direct-sender arg=warning
	docker-compose exec go cd/src & go mod tidy && go run exchange/direct/main.go sender $(arg)

direct-receive: ## Executa um consumer. Exchange direct. Ex: make direct-receive arg=warning
	docker-compose exec go cd/src & go mod tidy && go run exchange/direct/main.go receive $(arg)


