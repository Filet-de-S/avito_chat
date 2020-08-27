package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	listenerGlobal     net.Listener
	servHost, servPort string
	checkRate          time.Duration
	done               = make(chan bool, 2)
	timeout            = make(chan bool)
	isFirstConnect     = true
)

func main() {
	defer log.Println("exiting")

	user, pw := initENVs()
	json := genSecretsInJSON()
	wg := &sync.WaitGroup{}
	wg.Add(1)

	handler := fancyHandler(json, user, pw)
	go infServ(handler)
	go serverClose()
	wg.Wait()
}

func fancyHandler(json []byte, user, pw string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("New req: ", r)

		if r.Method != "POST" || r.URL.Path != "/secrets" {
			panic("who is this?")
		}

		u, pwBA, ok := r.BasicAuth()
		if !ok {
			log.Fatal("need pw, user in basic auth")
		}

		if pwBA != pw || u != user {
			log.Fatal("wrong pw")
		}

		_, err := w.Write(json)
		if err != nil {
			log.Fatal("couldn't write data:", err)
		}
		log.Println("Secrets send")
		done <- true
	})
}

func infServ(handler http.Handler) {
	for {
		listen, err := net.Listen("tcp", "pwmanager:1337")
		if err != nil {
			log.Fatal(err)
		}

		time.AfterFunc(5*time.Second, func() {
			timeout <- true
		})
		log.Println("Listening!")

		listenerGlobal = listen
		err = http.Serve(listenerGlobal, handler)
		switch err {
		case http.ErrServerClosed:
			checkServiceAlive()

		default:
			if strings.Contains(err.Error(), "use of closed network connection") {
				checkServiceAlive()
			} else {
				log.Println("something wrong with me. RESTART serv. Err:", err.Error())
			}
		}
	} // END FOR
}

func serverClose() {
	for {
		select {
		case <-done:
			<-time.After(time.Millisecond * 200)

			err := listenerGlobal.Close()
			log.Println("SERVER is closing with err: ", err)
		case <-timeout:
			err := listenerGlobal.Close()

			switch {
			case err != nil && strings.Contains(err.Error(),
				"use of closed network connection"):
			case err != nil:
				log.Println("SERVER is closing by timeout with err: ", err)
			default:
				log.Println("SERVER is closing by timeout")
			}
		}
	}
}

func checkServiceAlive() {
	url := "http://" + servHost + ":" + servPort + "/status"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("bad new request: %+v\nExiting", err)
		os.Exit(1)
	}

	client := &http.Client{}
	client.Timeout = 100 * time.Millisecond

	if isFirstConnect {
		time.Sleep(15 * time.Second)
		isFirstConnect = false
	}

	ticker := time.NewTicker(checkRate * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		resp, err := client.Do(req)
		if err != nil {
			log.Println("something wrong with api-service, RESTARTING serv. Err: ", err)
			return
		}

		if resp.StatusCode != http.StatusOK {
			log.Println("something wrong with api-service, RESTARTING serv. StatusCode:",
				resp.StatusCode, "body.close.err: ", resp.Body.Close())
			return
		}

		if err := resp.Body.Close(); err != nil {
			log.Println("can't close body: ", err)
		}
	}
}

func initENVs() (string, string) {
	// PW_MNG_FILE
	pwFile, ok := os.LookupEnv("PW_MNG_FILE")
	if !ok {
		log.Fatal("please set 'PW_MNG_FILE'")
	}

	pwB, err := ioutil.ReadFile("/run/secrets/" + pwFile)
	if err != nil {
		panic(err)
	}
	userAndPW := strings.Split(string(pwB), ":")

	// MAIN_SERVICE
	server, ok := os.LookupEnv("MAIN_SERVICE")
	if !ok {
		log.Fatal("please set 'MAIN_SERVICE'")
	}
	serviceSet := strings.Split(server, ":")
	servHost = serviceSet[0]
	servPort = serviceSet[1]

	// SERVICE_ALIVE_CHECKRATE
	scr, ok := os.LookupEnv("SERVICE_ALIVE_CHECKRATE")
	if !ok {
		log.Fatal("please set 'SERVICE_ALIVE_CHECKRATE'")
	}
	scrInt, err := strconv.Atoi(scr)
	if err != nil ||
		time.Duration(scrInt)*time.Millisecond > 5*time.Second {
		log.Fatal("please set 'SERVICE_ALIVE_CHECKRATE' in ms or <= 5000")
	}
	checkRate = time.Duration(scrInt)

	return userAndPW[0], userAndPW[1]
}

func genSecretsInJSON() []byte {
	jsn := map[string]string{}

	fls := strings.Split(os.Getenv("FILES"), ":")
	for i := range fls {
		cnt, err := ioutil.ReadFile("/run/secrets/" + fls[i])
		if err != nil {
			log.Fatal("error reading file: ", fls[i], err)
		}
		jsn[fls[i]] = string(cnt)
	}

	jsnBytes, err := json.Marshal(jsn)
	if err != nil {
		log.Fatal("error marshaling", err)
	}
	return jsnBytes
}
