package gene

import (
	"log"
	"math/rand"
)

func GenerateRandomInt(base int, power uint) []int {
	log.Print("start to generate random slice")
	arr := make([]int, base)
	for i := 0; i < base; i++ {
		arr[i] = rand.Intn(base)
	}
	for i := 0; i < int(power); i++ {
		arr = append(arr, arr...)
	}
	log.Print("finish to generate random slice")
	return arr
}
