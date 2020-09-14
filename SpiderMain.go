package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type xingZhengChuFa struct {
	uint64 int
	indexNum string
	sort string
	issuer string
	issueDate string
	issueNum string
	keywords string
	content string
	createdTime time.Time

}

var urlToNameMap = make(map[string]string)
var xingZhengChuFaUrl = "http://www.csrc.gov.cn/pub/zjhpublic"


func getXingZhengChuFaContent(){

}

func getXingZhengChuFaList(xingZhengChuFaListUrl string) {


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
		if err != nil || len(result)<=0{
			fmt.Println("匹配的url地址数为：%d",result)
			panic(err)
		} else {
			for index,item := range(result){
				line := string(item)
				index1 := strings.Index(line,"..")
				index2 := strings.Index(line,".htm")
				key := line[index1+5:index2+4]
				index3 := strings.Index(line,";\">")
				index4:=strings.Index(line,"</a>")
				value :=line[index3+3:index4]
				urlToNameMap[key]=strings.TrimSpace(value)
				fmt.Printf("%d,%s\n",index,line)
			}
		}
		fmt.Println(urlToNameMap)
	}
}

func main() {
	getXingZhengChuFaList("http://www.csrc.gov.cn/pub/zjhpublic/3300/3313/./index_7401.htm")
}
