package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// KeyValue is a type used to hold the key/value pairs passed to the map and reduce functions.
type KeyValue struct {
	Key   string
	Value string
}

// ReduceF function from MIT 6.824 LAB1
type ReduceF func(key string, values []string) string

// MapF function from MIT 6.824 LAB1
type MapF func(filename string, contents string) []KeyValue

// jobPhase indicates whether a task is scheduled as a map or reduce task.
type jobPhase string

const (
	mapPhase    jobPhase = "mapPhase"
	reducePhase          = "reducePhase"
)

type task struct {
	dataDir    string
	jobName    string
	mapFile    string   // only for map, the input file
	phase      jobPhase // are we in mapPhase or reducePhase?
	taskNumber int      // this task's index in the current phase
	nMap       int      // number of map tasks
	nReduce    int      // number of reduce tasks
	mapF       MapF     // map function used in this job
	reduceF    ReduceF  // reduce function used in this job
	wg         sync.WaitGroup
}

// MRCluster represents a map-reduce cluster.
type MRCluster struct {
	nWorkers int
	wg       sync.WaitGroup
	taskCh   chan *task
	exit     chan struct{}
}

var singleton = &MRCluster{
	nWorkers: runtime.NumCPU(),
	taskCh:   make(chan *task),
	exit:     make(chan struct{}),
}

func init() {
	singleton.Start()
}

// GetMRCluster returns a reference to a MRCluster.
func GetMRCluster() *MRCluster {
	return singleton
}

// NWorkers returns how many workers there are in this cluster.
func (c *MRCluster) NWorkers() int { return c.nWorkers }

// Start starts this cluster.
func (c *MRCluster) Start() {
	for i := 0; i < c.nWorkers; i++ {
		c.wg.Add(1)
		go c.worker()
	}
}

//别看，看就是暴力优化！！！
//1.replace the json serializing with raw io writing and reading
//2.replace the sort with hash map
func (c *MRCluster) worker() {
	defer c.wg.Done()
	for {
		select {
		case t := <-c.taskCh:
			if t.phase == mapPhase {
				content, err := ioutil.ReadFile(t.mapFile)
				if err != nil {
					panic(err)
				}

				fs := make([]*os.File, t.nReduce)
				bs := make([]*bufio.Writer, t.nReduce)
				outputs := make([][]*KeyValue, 0, len(fs))
				for i := range fs {
					rpath := reduceName(t.dataDir, t.jobName, t.taskNumber, i)
					fs[i], bs[i] = CreateFileAndBuf(rpath)
					outputs = append(outputs, make([]*KeyValue, 0))
				}
				results := t.mapF(t.mapFile, string(content))
				for i := range results {
					index := ihash(results[i].Key) % t.nReduce
					outputs[index] = append(outputs[index], &results[i])
				}
				for i := range fs {
					for j := range outputs[i] {
						fmt.Fprintf(bs[i], "%s:%s\n", outputs[i][j].Key, outputs[i][j].Value)
					}
					SafeClose(fs[i], bs[i])
				}
			} else {
				// YOUR CODE HERE :)
				//t.phase == reducePhase
				kvs := make([]*KeyValue, 0)
				for i := 0; i < t.nMap; i++ {
					rpath := reduceName(t.dataDir, t.jobName, i, t.taskNumber)
					f, r := OpenFileAndBuf(rpath)
					for true {
						if str, err := r.ReadString('\n'); err == io.EOF {
							break
						} else {
							str = strings.TrimSpace(str)
							strSlice := strings.Split(str, ":")
							kvs = append(kvs, &KeyValue{strSlice[0], strSlice[1]})
						}
					}
					f.Close()
				}
				replaceMap := make(map[string][]string)
				for i := range kvs {
					if s, has := replaceMap[kvs[i].Key]; !has {
						replaceMap[kvs[i].Key] = []string{kvs[i].Value}
					} else {
						replaceMap[kvs[i].Key] = append(s, kvs[i].Value)
					}
				}
				mrPath := mergeName(t.dataDir, t.jobName, t.taskNumber)
				f, w := CreateFileAndBuf(mrPath)
				if len(replaceMap) != 0 {
					for k, v := range replaceMap {
						w.WriteString(t.reduceF(k, v))
					}
				}
				SafeClose(f, w)
			}
			t.wg.Done()
		case <-c.exit:
			return
		}
	}
}

// Shutdown shutdowns this cluster.
func (c *MRCluster) Shutdown() {
	close(c.exit)
	c.wg.Wait()
}

// Submit submits a job to this cluster.
func (c *MRCluster) Submit(jobName, dataDir string, mapF MapF, reduceF ReduceF, mapFiles []string, nReduce int) <-chan []string {
	notify := make(chan []string)
	go c.run(jobName, dataDir, mapF, reduceF, mapFiles, nReduce, notify)
	return notify
}

func (c *MRCluster) run(jobName, dataDir string, mapF MapF, reduceF ReduceF, mapFiles []string, nReduce int, notify chan<- []string) {
	// map phase
	nMap := len(mapFiles)
	tasks := make([]*task, 0, nMap)
	for i := 0; i < nMap; i++ {
		t := &task{
			dataDir:    dataDir,
			jobName:    jobName,
			mapFile:    mapFiles[i],
			phase:      mapPhase,
			taskNumber: i,
			nReduce:    nReduce,
			nMap:       nMap,
			mapF:       mapF,
		}
		t.wg.Add(1)
		tasks = append(tasks, t)
		go func() { c.taskCh <- t }()
	}
	for _, t := range tasks {
		t.wg.Wait()

	}

	// reduce phase
	// YOUR CODE HERE :D
	tasks = make([]*task, 0, nReduce)
	for i := 0; i < nReduce; i++ {
		t := &task{
			dataDir: dataDir,
			jobName: jobName,
			phase: reducePhase,
			taskNumber: i,
			nMap: nMap,
			reduceF: reduceF,
		}
		t.wg.Add(1)
		tasks = append(tasks, t)
		go func() { c.taskCh <- t }()
	}
	for _, t := range tasks {
		t.wg.Wait()
	}
	mrFiles := make([]string, 0, nReduce)
	for i := 0; i < nReduce; i++ {
		mrFiles = append(mrFiles, mergeName(dataDir, jobName, i))
	}
	notify <- mrFiles
}

func ihash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32() & 0x7fffffff)
}

func reduceName(dataDir, jobName string, mapTask int, reduceTask int) string {
	return path.Join(dataDir, "mrtmp."+jobName+"-"+strconv.Itoa(mapTask)+"-"+strconv.Itoa(reduceTask))
}

func mergeName(dataDir, jobName string, reduceTask int) string {
	return path.Join(dataDir, "mrtmp."+jobName+"-res-"+strconv.Itoa(reduceTask))
}
