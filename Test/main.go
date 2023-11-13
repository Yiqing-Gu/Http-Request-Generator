package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

var fileArg = flag.String("f", "mock-http.json", "模拟数据文件名")
var onceCyclePublished []int64 = make([]int64, 0)
var repeatCyclePublished []int64 = make([]int64, 0)

func myfunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hi")
}

func main() {

	arr := [5]int{1, 2, 3, 4}
	fmt.Println(arr)

	http.HandleFunc("/", myfunc)
	http.ListenAndServe(":8080", nil)

	flag.Parse()
	fileName := *fileArg
	jsonFile, _ := os.Open(fileName)
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	jsonStr := string(byteValue)
	valueMap := gjson.Get(jsonStr, "config")
	if valueMap.Exists() {
		addr := valueMap.Get("baseUrl").String()
		body := valueMap.Get("body").String()
		path := valueMap.Get("path").String()
		resp, err := http.Post(addr, path, strings.NewReader(body))
		//fmt.Println(body)
		if err != nil {
			fmt.Println(addr)
			return
		}
		defer resp.Body.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

		onceArray := gjson.Get(jsonStr, "data.once").Array()
		fmt.Println(onceArray)
		fmt.Println("body")

	}
}
