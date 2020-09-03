curl -H "Authorization: f00b@r" localhost:9000/admin/pprof/heap -o results/raw/mem_out.txt
curl -H "Authorization: f00b@r" localhost:9000/admin/pprof/profile?seconds=1 -o results/raw/cpu_out.txt

go tool pprof -svg -alloc_objects results/raw/mem_out.txt > results/mem_ao.svg
go tool pprof -svg results/raw/cpu_out.txt > results/cpu.svg
