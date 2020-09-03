#!/bin/sh

# !NB: if you change UUID for GEN USER_ID or CHAT_ID, as well as pprof pw => script will not work!


curl -X POST localhost:9000/users/add -d '{"username": "user7"}'

curl -X POST localhost:9000/users/add -d '{"username": "user8"}'

curl -X POST localhost:9000/chats/add -d '{"name":"chat1","users":["0fdd1d54-736e-5749-90cb-987c32f5a3dc","a90bc930-77f8-5abc-91b1-cbcec106b07a"]}
'

mkdir -p results/raw

ab -n 5000 -c 100 -r -p sendMsg localhost:9000/messages/add > results/ab 2>&1 &

./pprof.sh

# waiting ab to finish
wait $!

sleep 1
curl -H "Authorization: f00b@r" localhost:9000/admin/pprof/goroutine?debug=2 -o results/goroutines.txt

echo "\nCheck 'results/' dir"