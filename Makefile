DOCKER_COMPOSE_APP=docker-compose -f deployments/docker-compose.yml
DOCKER_DEV=deployments/docker-compose.dev.yml
DOCKER_TEST=deployments/docker-compose.test.yml
DOCKER_TEST_CONTAINER=deployments/docker-compose.test.container.yml
DOCKER_PROD=deployments/docker-compose.prod.yml

RED=\033[0;31m
GREEN=\033[0;32m
CYAN=\033[0;36m
BLUE=\033[0;34m
NC=\033[0m
PASSED="${GREEN}Tests are PASSED${NC}"
FAILED="${RED}Tests are FAILED${NC}"
REPORT="${CYAN}Check report on: ${BLUE}test/postman/report.html${NC}"

dev-build:
	${DOCKER_COMPOSE_APP} -f ${DOCKER_DEV} build; exit $$?

dev-run: dev-build
	${DOCKER_COMPOSE_APP} -f ${DOCKER_DEV} up

dev-clean:
	${DOCKER_COMPOSE_APP} -f ${DOCKER_DEV} down --volume

test-run:
	${DOCKER_COMPOSE_APP} -f ${DOCKER_TEST} -f ${DOCKER_TEST_CONTAINER} up --build --exit-code-from test; \
		EXIT_CODE=$$?; test $$EXIT_CODE -eq 0 && echo ${PASSED} || echo ${FAILED}; \
	docker cp chat-api-test:/postman/report.html test/postman/ ;\
	echo ${REPORT} ;\
	exit $$EXIT_CODE

test-local-run:
	@chmod 0600 secrets/.pgpassForLocalTest
	${DOCKER_COMPOSE_APP} -f ${DOCKER_TEST} up --build -d
	SCRIPTS_FOLDER=scripts PGPASSFILE=secrets/.pgpassForLocalTest POSTMAN_FOLDER=test/postman \
 		SERVICE_HOST=localhost SERVICE_PORT=9000 SERVICE_NAME=CHAT URL_PATH=status PG_SNAME=localhost \
 		sh -c scripts/test.sh; EXIT_CODE=$$?; test $$EXIT_CODE -eq 0 && echo ${PASSED} || echo ${FAILED} ;\
	echo ${REPORT} ;\
	exit $$EXIT_CODE

prod: prod-build
	${DOCKER_COMPOSE_APP} -f ${DOCKER_PROD} up

prod-build:
	${DOCKER_COMPOSE_APP} -f ${DOCKER_PROD} build; exit $$?

lint:
	golangci-lint run; exit $$?