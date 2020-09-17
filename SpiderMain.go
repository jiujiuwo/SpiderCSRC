package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type XingZhengChuFaItem struct {
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
			for _,item := range(result){
				line := string(item)
				index1 := strings.Index(line,"..")
				index2 := strings.Index(line,".htm")
				key := line[index1+5:index2+4]
				index3 := strings.Index(line,";\">")
				index4:=strings.Index(line,"</a>")
				value :=line[index3+3:index4]
				urlToNameMap[key]=strings.TrimSpace(value)
			}
		}
	}
}


func getXingZhengChuFaDetail(urlMaps map[string]string){

	for key,value :=range(urlMaps){
		file, err := os.Create("./"+value)
		if err!=nil{
			fmt.Println("创建文件出错"+err.Error())
		}
		resp, err := http.Get(xingZhengChuFaUrl+key)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode==http.StatusOK{
			html, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			filterXingZhengChuFa(html)
			file.Write(html)
			file.Close()
		}
	}
}


func filterXingZhengChuFa(html  []byte){
	item := XingZhengChuFaItem{}
	//提取索引号
	indexNumPattern,err :=regexp.Compile("索 引 号:</B>.*/</td>")
	if err!=nil{
		panic(err)
	}
	indexNum := string(indexNumPattern.Find(html))
	start :=strings.Index(indexNum,"</B>")
	end := strings.Index(indexNum,"</td>")
	item.indexNum = indexNum[start+4:end]

	//提取分类
	sortPattern,err := regexp.Compile("<span id=\"lSubcat\">.*&nbsp;;&nbsp;.*</span>")
	if err!=nil{
		panic(err)
	}
	sort := string(sortPattern.Find(html))
	start = strings.Index(sort,"\">")
	end = strings.Index(sort,"</span>")
	sort = sort[start+2:end]
	sort = strings.ReplaceAll(sort,"&nbsp;","")
	item.sort = sort

	//提取发布机构
	issuserPattern,err := regexp.Compile("<B>发布机构:</B> \n.*<span>.*</span>")
	if err !=nil{
		panic(err)
	}
	issuer := string(issuserPattern.Find(html))
	start = strings.Index(issuer,"<span>")
	end = strings.Index(issuer,"</span>")
	issuer = issuer[start+6:end]
	item.issuer = issuer
	fmt.Println(item)
}

func main() {
	getXingZhengChuFaList("http://www.csrc.gov.cn/pub/zjhpublic/3300/3313/./index_7401.htm")
	getXingZhengChuFaDetail(urlToNameMap)
}
