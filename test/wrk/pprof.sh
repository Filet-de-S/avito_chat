#!/bin/sh

pw=$PPROF_PW
url=localhost:9000/admin/pprof
lim=$TLIM
[ -z "$lim" ] && lim=1

curl -H "Authorization: $pw" "$url/heap" -o results/raw/mem_out.txt
curl -H "Authorization: $pw" "$url/profile?seconds=$lim" -o results/raw/cpu_out.txt

go tool pprof -svg -alloc_objects results/raw/mem_out.txt > results/mem_ao.svg
go tool pprof -svg results/raw/cpu_out.txt > results/cpu.svg
