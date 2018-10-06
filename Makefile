.PHONY: help

# self-documented makefile, thanks to the Marmelab team
# see http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

DOCKER_PREFIX ?= chrome_proxy
DOCKER_HTTP_PORT ?= 9222

docker-build-chrome: ## build the chrome headless container using docker
	docker build --rm --tag $(DOCKER_PREFIX) chrome

docker-run-chrome: ## start running the chrome headless container
	docker run --rm --detach --publish $(DOCKER_HTTP_PORT):9222 --name $(DOCKER_PREFIX)_$(DOCKER_HTTP_PORT) $(DOCKER_PREFIX)

docker-stop-chrome: ## stop the chrome headless container
	docker stop $(DOCKER_PREFIX)_$(DOCKER_HTTP_PORT)

docker-logs-chrome: ## display the logs from the chrome headless container
	docker logs $(DOCKER_PREFIX)_$(DOCKER_HTTP_PORT)
