.PHONY: all

#all: test_example test_homework cleanup gendata
all: cleanup test_example cleanup test_homework
test_example:
	go test -v -run=TestExampleURLTop -cpuprofile cpu.out -memprofile mem.out

test_homework:
	go test -v -run=TestURLTop -cpuprofile cpu.out -memprofile mem.out

cleanup:
	go test -v -run=TestCleanData

gendata:
	go test -v -run=TestGenData
