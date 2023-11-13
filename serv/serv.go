package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/tidwall/gjson"
)

var fileArg = flag.String("f", "serverConfig.json", "服务器文件名")

func serverPrint(w http.ResponseWriter, r *http.Request) {
	bodyinfo, _ := ioutil.ReadAll(r.Body) // get body content from server

	// printing from server
	fmt.Println("baseUrl:" + r.Host)
	fmt.Println("path:" + r.URL.Path)
	fmt.Println("mode:" + (r.Method))
	for ki, va := range r.Header {
		fmt.Println("header: key:" + ki + " value:" + va[0])
	}

	q := r.URL.Query()

	for ki, va := range q {
		fmt.Println("query: key:" + ki + " value:" + va[0])
	}
	fmt.Println("body:" + string(bodyinfo))
}

func main() {

	flag.Parse()
	fileName := *fileArg
	jsonFile, _ := os.Open(fileName)
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	jsonStr := string(byteValue)
	valueMap := gjson.Get(jsonStr, "config")

	if valueMap.Exists() {

		port := ":" + valueMap.Get("port").String()
		pathCollection := valueMap.Get("path").Array()
		routerList := make(map[string]int)

		for _, path := range pathCollection {
			if routerList[path.String()] == 0 {
				http.HandleFunc(path.String(), serverPrint)
				routerList[path.String()] = 1
			} else {
				fmt.Println("router does exist.")
			}
		}
		http.ListenAndServe(port, nil)
	}

}
