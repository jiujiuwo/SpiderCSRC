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
	uint64      int
	indexNum    string
	sort        string
	issuer      string
	issueDate   string
	issueNum    string
	keywords    string
	content     string
	name        string
	createdTime time.Time
}

var urlToNameMap = make(map[string]string)
var xingZhengChuFaUrl = "http://www.csrc.gov.cn/pub/zjhpublic"

func getXingZhengChuFaList(xingZhengChuFaListUrl string) {

	resp, err := http.Get(xingZhengChuFaListUrl)
	if err != nil {
		panic(err)
	}

	pattern, err := regexp.Compile("<a href=\"../../G00306212/.*</a>")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		html, err := ioutil.ReadAll(resp.Body)
		result := pattern.FindAll(html, -1)
		if err != nil || len(result) <= 0 {
			fmt.Println("匹配的url地址数为：%d", result)
			panic(err)
		} else {
			for _, item := range result {
				line := string(item)
				index1 := strings.Index(line, "..")
				index2 := strings.Index(line, ".htm")
				key := line[index1+5 : index2+4]
				index3 := strings.Index(line, ";\">")
				index4 := strings.Index(line, "</a>")
				value := line[index3+3 : index4]
				urlToNameMap[key] = strings.TrimSpace(value)
			}
		}
	}
}

/*
	根据URL与行政处罚名称的映射，访问行政处罚详情页
*/
func getXingZhengChuFaDetail(urlMaps map[string]string) {

	for key, value := range urlMaps {
		file, err := os.Create("./" + value)
		if err != nil {
			fmt.Println("创建文件出错" + err.Error())
		}
		resp, err := http.Get(xingZhengChuFaUrl + key)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode == http.StatusOK {
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

/*
	提取html字节符数组中的，行政处罚的 索引号，分类，发布机构，
*/
func filterXingZhengChuFa(html []byte) {
	item := XingZhengChuFaItem{}
	//提取索引号
	indexNumPattern, err := regexp.Compile("索 引 号:</B>.*/</td>")
	if err != nil {
		panic(err)
	}
	indexNum := string(indexNumPattern.Find(html))
	start := strings.Index(indexNum, "</B>")
	end := strings.Index(indexNum, "</td>")
	if !indexOutOfS(indexNum,start,end){
		item.indexNum = indexNum[start+4 : end]
	}


	//提取分类
	sortPattern, err := regexp.Compile("<span id=\"lSubcat\">.*&nbsp;;&nbsp;.*</span>")
	if err != nil {
		panic(err)
	}
	sort := string(sortPattern.Find(html))
	start = strings.Index(sort, "\">")
	end = strings.Index(sort, "</span>")
	if !indexOutOfS(sort,start,end){
		sort = sort[start+2 : end]
	}
	sort = strings.ReplaceAll(sort, "&nbsp;", "")
	item.sort = sort

	//提取发布机构
	issuserPattern, err := regexp.Compile("<B>发布机构:</B> \n.*<span>.*</span>")
	if err != nil {
		panic(err)
	}
	issuer := string(issuserPattern.Find(html))
	start = strings.Index(issuer, "<span>")
	end = strings.Index(issuer, "</span>")
	if !indexOutOfS(issuer,start,end){
		item.issuer = issuer[start+6 : end]
	}
	//提取其他内容
	filterXingZhengChuFa2(&item, html)
	item.createdTime = time.Now()
	fmt.Println(item)
}

/*
	提取html字节符数组中的，发文日期，名称，文号，主题词，内容
*/
func filterXingZhengChuFa2(item *XingZhengChuFaItem, html []byte) {
	//提取发文日期
	datePattern, err := regexp.Compile("<B>发文日期:</B>\n.*<span>.*</span>")
	if err != nil {
		panic(err)
	}
	datetime := string(datePattern.Find(html))
	start := strings.Index(datetime, "<span>")
	end := strings.Index(datetime, "</span>")
	if !indexOutOfS(datetime,start,end){
		item.issueDate = datetime[start+6 : end]
	}

	//提取名称
	namePattern, err := regexp.Compile("<B>名　　称:</B> \n.*<span id=\"lTitle\">.*</span>")
	if err != nil {
		panic(err)
	}
	name := string(namePattern.Find(html))
	start = strings.Index(name, "<span")
	end = strings.Index(name, "</span>")
	if !indexOutOfS(name,start,end){
		item.name = name[start+18:end]
	}

	//提取文号
	keywordsPattern, err := regexp.Compile("<B>主 题 词:</B> \n.*<span>.*</span>")
	if err != nil {
		panic(err)
	}
	keywords := string(keywordsPattern.Find(html))
	start = strings.Index(keywords, "<span>")
	end = strings.Index(keywords, "</span>")
	if !indexOutOfS(datetime,start,end){
		item.keywords = keywords[start+6:end]
	}

	//提取主题词
	issueNumPattern, err := regexp.Compile("<B>文　　号:</B> \n.*<span>.*</span>")
	if err != nil {
		panic(err)
	}
	issueNum := string(issueNumPattern.Find(html))
	start = strings.Index(issueNum, "<span>")
	end = strings.Index(issueNum, "</span>")
	if !indexOutOfS(issueNum,start,end){
		item.issueNum = issueNum[start+6:end]
	}
	//提取内容
	contentPattern, err := regexp.Compile("<P(.|d)*<SPAN(.|d)*</SPAN>(.|d)*</P>")
	if err != nil {
		panic(err)
	}
	//var content string
	contents := contentPattern.FindAll(html,-1)
	for _,per := range(contents){
		tmp  := string(per)
/*		if -1==strings.Index(tmp,"</FONT>"){
			start = strings.Index(tmp,"10.5pt\">")
			end = strings.Index(tmp,"</SPAN>")
			if !indexOutOfS(tmp,start,end){
				content+=tmp[start:end]
			}
		}else{
			start = strings.Index(tmp,"0>")
			end = strings.Index(tmp,"</FONT>")
			if !indexOutOfS(tmp,start,end){
				content+=tmp[start:end]
			}
		}*/
		//item.content = content
		fmt.Println(tmp)
	}
}
/*
	判断截取的字符串是否越界
 */
func indexOutOfS(s string,start int,end int) bool{
	if start<0||end<0{
		return true
	}
	if start>len(s)||end>len(s){
		return true
	}
	return false
}
func main() {
	getXingZhengChuFaList("http://www.csrc.gov.cn/pub/zjhpublic/3300/3313/./index_7401.htm")
	getXingZhengChuFaDetail(urlToNameMap)
}
