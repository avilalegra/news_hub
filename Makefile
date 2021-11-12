run/local: run/dev
	@xdg-open http://localhost:8085

run/dev:
	@docker-compose up -d

test/unit: run/dev
	@docker-compose exec go bash -c 'GOFLAGS="-count=1";SKIP_DB_INIT="true" go test -short ./...'

test/int: run/dev
	@@docker-compose exec go bash -c 'GOFLAGS="-count=1" go test -run Integration -p 1 ./...'

test/all: test/clrcache test/unit test/int

test/clrcache:
	@@docker-compose exec go bash -c 'go clean -testcache'

stop:
	@docker-compose stop

destroy:
	@docker-compose down --rmi all

config/reload:
	@docker-compose exec go bash -c 'pgrep main | xargs kill -s SIGUSR1'