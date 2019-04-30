package main

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// URLTop10 .
func URLTop10(nWorkers int) RoundsArgs {
	// YOUR CODE HERE :)
	// And don't forget to document your idea.
	var args RoundsArgs
	// round 1: do url count
	args = append(args, RoundArgs{
		MapFunc:    URLTop10MapCount,
		ReduceFunc: URLTop10ReduceCount,
		NReduce:    nWorkers,
	})

	args = append(args, RoundArgs{
		MapFunc: URLTop10MapMerge,
		ReduceFunc: URLTop10ReduceMerge,
		NReduce: 1,
	})
	return args
}
func URLTop10MapCount(filename string, contents string) []KeyValue  {
	tmp := make(map[string]int)
	lines := strings.Split(contents, "\n")
	for k := range lines {
		if len(lines[k]) == 0 {
			continue
		}
		if num, ok := tmp[lines[k]]; ok {
			tmp[lines[k]] = num + 1
		} else {
			tmp[lines[k]] = 1
		}
	}
	kvs := make([]KeyValue, 0, len(tmp))
	for k, v := range tmp {
		kvs = append(kvs, KeyValue{k, strconv.FormatInt(int64(v), 10)})
	}
	return kvs
}
func URLTop10ReduceCount(key string, values []string) string  {
	total := int64(0)
	for _, v := range values {
		num, _ := strconv.ParseInt(v, 10, 64)
		total += num
	}
	return fmt.Sprintf("%s %s\n", key, strconv.FormatInt(total, 10))
}
func URLTop10MapMerge(filename string, contents string) []KeyValue  {
	lines := strings.Split(contents, "\n")
	kvs := make([]KeyValue, 0, len(lines))
	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		kvs = append(kvs, KeyValue{"", l})
	}
	return kvs
}
func URLTop10ReduceMerge(key string, values []string) string  {
	kvs := make([]urlCount, 0, len(values))
	for k := range values {
		values[k] = strings.TrimSpace(values[k])
		s := strings.Split(values[k], " ")
		num, _ := strconv.Atoi(s[1])
		kvs = append(kvs, urlCount{s[0], num})
	}
	return top10(kvs)
}
func top10(kvs []urlCount) string {
	i, j := 0, len(kvs) - 1
	if j > 9 {
		for j!= 9 {
			f := i
			start, end := i, j
			for i < j {
				for kvs[f].cnt < kvs[i].cnt &&  i < end {
					i++
				}
				for (kvs[f].cnt > kvs[j].cnt ||(kvs[f].cnt == kvs[j].cnt && kvs[f].url < kvs[j].url)) && j > start {
					j--
				}
				kvs[i], kvs[j] = kvs[j], kvs[i]
			}
			kvs[i], kvs[j] = kvs[j], kvs[i]
			kvs[f], kvs[j] = kvs[j], kvs[f]
			if j > 9 {
				i = start
				j = j - 1
			} else if j < 9 {
				i = j + 1
				j = end
			}
		}
		kvs = kvs[0:10]
	}

	sort.Slice(kvs, func(i, j int) bool {
		if kvs[i].cnt > kvs[j].cnt || kvs[i].cnt == kvs[j].cnt && kvs[i].url < kvs[j].url {
			return true
		} else {
			return false
		}
	})
	buf := new(bytes.Buffer)
	for i := range kvs {
		fmt.Fprintf(buf, "%s: %d\n", kvs[i].url, kvs[i].cnt)
	}
	return buf.String()
}
