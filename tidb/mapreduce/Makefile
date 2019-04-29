.PHONY: all

#all: test_example test_homework cleanup gendata
all: cleanup test_example
test_example:
	go test -v -run=TestExampleURLTop -cpuprofile cpu.out -memprofile mem.out

test_homework:
	go test -v -run=TestURLTop

cleanup:
	go test -v -run=TestCleanData

gendata:
	go test -v -run=TestGenData
test_chan:
	go test -v -run=TestChannelSendAndReceive -cpuprofile chancpu.out -memprofile chanmem.out -count 10
