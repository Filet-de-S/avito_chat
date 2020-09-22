DOCKER_CMP_N=docker-compose -p
CMP_BASE=-f deployments/docker-compose.yml

CMP_DEV=deployments/docker-compose.dev.yml
CMP_TEST=deployments/docker-compose.test.yml
CMP_TEST_CNTR=deployments/docker-compose.test.container.yml
CMP_PROD=deployments/docker-compose.prod.yml
DEV=dev
TEST=test
TEST_CNTR=test-contr
PROD=prod

RED=\033[0;31m
GREEN=\033[0;32m
CYAN=\033[0;36m
BLUE=\033[0;34m
NC=\033[0m

PASSED="${GREEN}Tests are PASSED${NC}"
FAILED="${RED}Tests are FAILED${NC}"
REPORT="${CYAN}Check report on: ${BLUE}test/postman/report.html${NC}"

dev-build:
	${DOCKER_CMP_N} ${DEV} ${CMP_BASE} -f ${CMP_DEV} build; exit $$?

dev-run: dev-build
	${DOCKER_CMP_N} ${DEV} ${CMP_BASE} -f ${CMP_DEV} up

dev-clean:
	${DOCKER_CMP_N} ${DEV} ${CMP_BASE} -f ${CMP_DEV} down --volume

test-run:
	@chmod 0600 secrets/.pgpass
	${DOCKER_CMP_N} ${TEST_CNTR} ${CMP_BASE} -f ${CMP_TEST}  \
		-f ${CMP_TEST_CNTR} up --build --exit-code-from test; \
		EXIT_CODE=$$?; test $$EXIT_CODE -eq 0 && echo ${PASSED} \
		|| echo ${FAILED}; \
	docker cp test-contr_api-tests:/postman/report.html test/postman/ ;\
	echo ${REPORT} ;\
	exit $$EXIT_CODE

# use, if you need db init, however tables are deleted in normal test-run
test-clean:
	${DOCKER_CMP_N} ${TEST_CNTR} ${CMP_BASE} -f ${CMP_TEST} -f ${CMP_TEST_CNTR} down --volume

test-local-run:
	@chmod 0600 secrets/.pgpassLocalhost
	${DOCKER_CMP_N} ${TEST} ${CMP_BASE} -f ${CMP_TEST} up --build -d
	SCRIPTS_FOLDER=scripts PGPASSFILE=secrets/.pgpassLocalhost POSTMAN_FOLDER=test/postman \
 		SERVICE_HOST=localhost SERVICE_PORT=9000 SERVICE_NAME=CHAT URL_PATH=status PG_SNAME=localhost \
 		sh -c scripts/test.sh; EXIT_CODE=$$?; test $$EXIT_CODE -eq 0 && echo ${PASSED} || echo ${FAILED} ;\
	${DOCKER_CMP_N} ${TEST} ${CMP_BASE} -f ${CMP_TEST} stop ;\
	echo ${REPORT} ;\
	exit $$EXIT_CODE

# use, if you need db init, however tables are deleted in normal test-local-run
test-local-clean:
	${DOCKER_CMP_N} ${TEST} ${CMP_BASE} -f ${CMP_TEST} down --volume

prod: prod-build
	${DOCKER_CMP_N} ${PROD} ${CMP_BASE} -f ${CMP_PROD} up

prod-build:
	${DOCKER_CMP_N} ${PROD} ${CMP_BASE} -f ${CMP_PROD} build; exit $$?

lint:
	golangci-lint run; exit $$?

# don't forget to: `export PPROF=ON` before service start!
PPROF_PW=$(shell cat secrets/.pprof)
ab:
	cd test/ApacheBench; eval './ab+pprof.sh $(PPROF_PW) $(tlim) $(args)'

# don't forget to: `export PPROF=ON` before service start!
wrk:
	cd test/wrk; eval './wrk+pprof.sh $(PPROF_PW) $(tlim) $(args)'