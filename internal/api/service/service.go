package service

import (
	chats "avito-chat_service/internal/api/contexts/chats/usecase"
	msg "avito-chat_service/internal/api/contexts/messages/usecase"
	users "avito-chat_service/internal/api/contexts/users/usecase"
	"avito-chat_service/internal/api/db"
	"avito-chat_service/internal/api/db/psql"
	_router "avito-chat_service/internal/api/service/router"
	_server "avito-chat_service/internal/api/service/server"
	"avito-chat_service/internal/api/uuidgen"
	"avito-chat_service/internal/api/uuidgen/gofrsuuid"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgpassfile"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Service ...
type Service struct {
	server *_server.Server
}

type config struct {
	mode    string
	port    int
	timeout time.Duration
	// servAliveChR time.Duration
	pgPassFile     io.Reader
	uuidUser       string
	uuidChat       string
	uuidMsg        string
	pgPoolMaxConns int32
}

type _envs struct {
	pwMngFile      string
	pwMngHost      string
	pwMngPort      string
	serviceMode    string
	serviceTimeout string
	// servAliveChR   string
	pgPassFile     string
	pgPoolMaxConns string
	uuidsFile      string
}

type memory struct {
	dataStore db.Service
}

type services struct {
	mem     memory
	uuidGen uuidgen.Service
}

var (
	cfg  = config{}
	envs = _envs{}
)

// New ...
func New() (*Service, error) {
	err := initConfig()
	if err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}

	services, err := initServices()
	if err != nil {
		return nil, fmt.Errorf("init services: %w", err)
	}

	useCases, err := initUseCases(services)
	if err != nil {
		return nil, fmt.Errorf("useCases: %w", err)
	}

	router, err := _router.New(useCases)
	if err != nil {
		return nil, fmt.Errorf("routing: %w", err)
	}

	server, err := _server.New(router.Handler(), cfg.timeout, cfg.port)
	if err != nil {
		return nil, fmt.Errorf("server: %w", err)
	}

	return &Service{server: server}, nil
}

// Run ...
func (s *Service) Run() error {
	log.Println("Starting service in", cfg.mode, "mode on", cfg.port, "port")
	return s.server.Start()
}

func initConfig() error {
	// INIT ENVs
	err := initENVs()
	if err != nil {
		return errors.New("envs: " + err.Error())
	}

	// INIT SECRETS
	err = initSecretsFromPWManager()
	if err != nil {
		return errors.New("pw manager: " + err.Error())
	}

	// SERVICE_MODE
	cfg.mode = envs.serviceMode
	gin.SetMode(cfg.mode)

	// EnableDecoderDisallowUnknownFields
	binding.EnableDecoderDisallowUnknownFields = true

	// PORT
	cfg.port = 8080

	// SERVICE TIMEOUT
	stAtoi, err := strconv.Atoi(envs.serviceTimeout)
	if err != nil {
		return errors.New("please set `SERVICE_TIMEOUT` in valid format: ms, just int")
	}
	cfg.timeout = time.Duration(stAtoi) * time.Millisecond

	// PGPOOL_MAXCONNS
	maxConns, err := strconv.Atoi(envs.pgPoolMaxConns)
	if err != nil {
		return errors.New("please set `PGPOOL_MAXCONNS` in valid format: just int")
	}
	cfg.pgPoolMaxConns = int32(maxConns)

	//// SERVICE_ALIVE_CHECKRATE
	//srvAl, err := strconv.Atoi(envs.servAliveChR)
	//if err != nil {
	//	return errors.New("please set `SERVICE_ALIVE_CHECKRATE` in valid format: ms, just int")
	//}
	//cfg.servAliveChR = time.Duration(srvAl) * time.Millisecond
	//log.Printf("Timeout: %v, server alive check rate: %v\n", cfg.timeout, cfg.servAliveChR)
	return nil
}

func initSecretsFromPWManager() error {
	pw, err := ioutil.ReadFile("/run/secrets/" + envs.pwMngFile)
	if err != nil {
		return errors.New("can't read 'PW_MNG_FILE' content: " + err.Error())
	}
	userAndPw := strings.Split(string(pw), ":")

	jsn := fetchSecrets(envs.pwMngHost, envs.pwMngPort, userAndPw[0], userAndPw[1])
	if err != nil {
		return errors.New("fetching fail: " + err.Error())
	}

	secrets := map[string]interface{}{}
	err = json.Unmarshal(jsn, &secrets)
	if err != nil {
		return err
	}

	// PGPASSFILE
	val, ok := secrets[envs.pgPassFile]
	if !ok {
		return errors.New("can't find 'PGPASSFILE' in secrets JSON")
	}
	cfg.pgPassFile = strings.NewReader(val.(string))

	// UUIDS_FILE
	val, ok = secrets[envs.uuidsFile]
	if !ok {
		return errors.New("can't find 'UUIDS_FILE' in secrets JSON")
	}
	parseUUIDsFile(val.(string))

	return nil
}

func parseUUIDsFile(s string) {
	pairs := strings.Split(s, "\n")
	for i := range pairs {
		vals := strings.Split(pairs[i], ":")
		switch vals[0] {
		case "USER":
			cfg.uuidUser = vals[1]
		case "CHAT":
			cfg.uuidChat = vals[1]
		case "MSG":
			cfg.uuidMsg = vals[1]
		}
	}
}

func fetchSecrets(host, port, user, pw string) []byte {
	url := "http://" + host + ":" + port + "/secrets"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Printf("bad new request: %+v\nExiting", err)
		os.Exit(1)
	}

	req.SetBasicAuth(user, pw)

	client := &http.Client{}
	client.Timeout = 100 * time.Millisecond

	i := 0
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	log.Print("Connecting to pwmanager . . .")
	for {
		select {
		case <-time.After(time.Second * 6):
			log.Println("seems pwmanager is down, exiting")
			os.Exit(1)
		case <-ticker.C:
			bodyInBytes, err := doFetch(client, req)
			if err != nil {
				if i%10 == 0 {
					log.Printf("can't reach pw manager . . . : %+v", err)
					i++
				}

				break
			}

			return bodyInBytes
		}
	}
}

func doFetch(client *http.Client, req *http.Request) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Do: " + err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("status code: " + strconv.Itoa(resp.StatusCode))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("can't read body: " + err.Error())
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	log.Println("Got secrets")
	return respBody, nil
}

func initServices() (*services, error) {
	memory, err := dbInit()
	if err != nil {
		return nil, errors.New("db init: " + err.Error())
	}

	uuidGen, err := uuidGenInit()
	if err != nil {
		return nil, errors.New("uuidGen init: " + err.Error())
	}

	return &services{
		mem:     memory,
		uuidGen: uuidGen,
	}, nil
}

func initUseCases(s *services) (*_router.UseCases, error) {
	users, err := users.New(s.mem.dataStore, s.uuidGen)
	if err != nil {
		return nil, errors.New("users init: " + err.Error())
	}

	chats, err := chats.New(s.mem.dataStore, s.uuidGen)
	if err != nil {
		return nil, errors.New("chats init: " + err.Error())
	}

	msg, err := msg.New(s.mem.dataStore, s.uuidGen)
	if err != nil {
		return nil, errors.New("msg init: " + err.Error())
	}

	return &_router.UseCases{
		Users: users,
		Chats: chats,
		MSG:   msg,
	}, nil
}

func uuidGenInit() (uuidgen.Service, error) {
	var err error
	uuids := gofrsuuid.UUIDs{}

	uuids.User, err = uuid.FromString(cfg.uuidUser)
	if err != nil {
		return nil, err
	}
	uuids.Chat, err = uuid.FromString(cfg.uuidChat)
	if err != nil {
		return nil, err
	}
	uuids.Msg, err = uuid.FromString(cfg.uuidMsg)
	if err != nil {
		return nil, err
	}

	uidService, err := gofrsuuid.New(uuids)
	if err != nil {
		return nil, err
	}

	return uidService, nil
}

func dbInit() (memory, error) {
	config, err := getDBConfig()
	if err != nil {
		return memory{}, err
	}

	var pool *pgxpool.Pool
	timer := time.After(15 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	i := 0

CONNECT:
	for {
		if i == 0 || i%10 == 0 {
			log.Println("Trying to db connect")
			i++
		}

		select {
		case <-timer:
			log.Println("too long connect to db: ", err)
			os.Exit(1)
		case <-ticker.C:
			pool, err = pgxpool.ConnectConfig(context.Background(), config)
			if err == nil {
				break CONNECT
			}
		}
	}

	log.Println("Connection with db established")

	store, err := psql.New(pool)
	if err != nil {
		return memory{}, err
	}

	return memory{dataStore: store}, nil
}

func getDBConfig() (*pgxpool.Config, error) {
	pf, err := pgpassfile.ParsePassfile(cfg.pgPassFile)
	if err != nil {
		return nil, err
	}

	connConfig, err := pgx.ParseConfig("")
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(pf.Entries[0].Port)
	if err != nil {
		return nil, err
	}

	connConfig.RuntimeParams = map[string]string{
		"standard_conforming_strings": "on",
		"backslash_quote":             "off",
	}
	connConfig.PreferSimpleProtocol = true
	connConfig.Port = uint16(port)
	connConfig.Host = pf.Entries[0].Hostname
	connConfig.Database = pf.Entries[0].Database
	connConfig.User = pf.Entries[0].Username
	connConfig.Password = pf.Entries[0].Password

	poolConf, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, err
	}

	poolConf.ConnConfig = connConfig
	poolConf.MaxConns = cfg.pgPoolMaxConns

	return poolConf, nil
}

func initENVs() error {
	// PW_MNG_FILE
	v, ok := os.LookupEnv("PW_MNG_FILE")
	if !ok {
		return errors.New("set 'PW_MNG_FILE' env")
	}
	envs.pwMngFile = v

	// PW_MNG_HOST
	envs.pwMngHost, ok = os.LookupEnv("PW_MNG_HOST")
	if !ok {
		return errors.New("set 'PW_MNG_HOST' env")
	}

	// PW_MNG_PORT
	envs.pwMngPort, ok = os.LookupEnv("PW_MNG_PORT")
	if !ok {
		return errors.New("set 'PW_MNG_PORT' env")
	}

	// SERVICE_MODE
	envs.serviceMode, ok = os.LookupEnv("SERVICE_MODE")
	if !ok {
		return errors.New("set 'SERVICE_MODE' env")
	}

	// SERVICE_TIMEOUT
	envs.serviceTimeout, ok = os.LookupEnv("SERVICE_TIMEOUT")
	if !ok {
		return errors.New("please set `SERVICE_TIMEOUT` in ms")
	}

	//// SERVICE_ALIVE_CHECKRATE
	//envs.servAliveChR, ok = os.LookupEnv("SERVICE_ALIVE_CHECKRATE")
	//if !ok {
	//	return errors.New("please set `SERVICE_TIMEOUT` in ms")
	//}

	// PGPASSFILE
	envs.pgPassFile, ok = os.LookupEnv("PGPASSFILE")
	if !ok {
		return errors.New("set 'PGPASSFILE'")
	}

	// PGPOOL_MAXCONNS
	envs.pgPoolMaxConns, ok = os.LookupEnv("PGPOOL_MAXCONNS")
	if !ok {
		return errors.New("set 'PGPOOL_MAXCONNS'")
	}

	// UUIDS_FILE
	envs.uuidsFile, ok = os.LookupEnv("UUIDS_FILE")
	if !ok {
		return errors.New("set 'UUIDS_FILE'")
	}

	return nil
}
