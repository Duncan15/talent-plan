.PHONY: all

all: test bench pprof vet

test:
	go test

bench:
	go test -bench Benchmark -run xx -count 5 -benchmem

pprof:
	go test -bench BenchmarkMergeSort -run xx -cpuprofile cpu.out -memprofile mem.out

vet:
	go test -bench BenchmarkVer -run xx -cpuprofile vercpu.out -memprofile vermem.out
