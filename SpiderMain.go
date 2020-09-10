package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func getXingZhengChuFaList(xingZhengChuFaListUrl string) {
	urlToNameMap :=make(map[string]string)

	resp, err := http.Get(xingZhengChuFaListUrl)
	if err != nil {
		panic(err)
	}

	pattern,err :=regexp.Compile("<a href=\"../../G00306212/.*</a>")
	if(err!=nil){
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		html, err := ioutil.ReadAll(resp.Body)
		result := pattern.FindAll(html,-1)
		if err != nil {
			panic(err)
		} else {
			fmt.Println(len(result))
			for index,item := range(result){
				line := string(item)
				index1 := strings.Index(line,"..")
				index2 := strings.Index(line,".htm")
				key := line[index1+1:index2]
				index3 := strings.Index(line,";\">")
				index4:=strings.Index(line,"</a>")
				value :=line[index3+3:index4]
				urlToNameMap[key]=value
				fmt.Printf("%d,%s\n",index,line)
			}
		}
		fmt.Println(urlToNameMap)
	}
}

func main() {
	getXingZhengChuFaList("http://www.csrc.gov.cn/pub/zjhpublic/3300/3313/./index_7401.htm")
}
