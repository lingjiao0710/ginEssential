package util

import (
	"math/rand"
	"time"
)

//RandomString 生成一个随机字符串
func RandomString(n int) (s string) {
	var letters = []byte("asdf;lkjqwerpouiasdfjfjkalsdfjASDFQWERASDXCVB")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}
