package controllers

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"fmt"
	"log"
)

// Result 学分结果
type Result struct {
	Num  string
	Name string
	Xf   string
	Res  map[string]string
	Jd   string
}

// Info 学号，姓名
type Info struct {
	Name string
	Num  string
}

// Grade 成绩
type Grade struct {
	Num         string
	Name        string
	Graderesult [][]string
}

func matchgrade(response1 *http.Response) [][]string {
	dec := mahonia.NewDecoder("gbk")
	doc := dec.NewReader(response1.Body)
	result, _ := goquery.NewDocumentFromReader(doc)
	graderesult := make([][]string, 0)
	result.Find(".datelist").Find("tr").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			row := make([]string, 6)
			row[0] = s.Find("td").Eq(3).Text()
			row[1] = s.Find("td").Eq(4).Text()
			row[2] = s.Find("td").Eq(6).Text()
			row[3] = s.Find("td").Eq(7).Text()
			row[4] = s.Find("td").Eq(8).Text()
			row[5] = s.Find("td").Eq(12).Text()
			graderesult = append(graderesult, row)
		}
	})

	return graderesult
}

func matchcredit(response1 *http.Response) (map[string]string, string, string) {
	ma := make(map[string]string)
	dec := mahonia.NewDecoder("gbk")
	doc := dec.NewReader(response1.Body)
	result, _ := goquery.NewDocumentFromReader(doc)
	result.Find(".datelist").Eq(0).Find("tr").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			ma[s.Find("td").Eq(0).Text()] = s.Find("td").Eq(2).Text()
		}
	})
	jd := result.Find("#pjxfjd").Text()
	xf := result.Find("#xftj").Text()
	return ma, xf, jd
}

func getViewState(url string) string{
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	//增加header选项
	request.Header.Add("Referer", "http://202.116.160.170/default2.aspx")
	request.Header.Add("Accept", "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	if err != nil {
		return ""
	}
	response, _ := client.Do(request)
	response,_=http.Get(url)
	defer response.Body.Close()
	if err !=nil{
		fmt.Println("页面解析失败")
		return ""
	}
	result,err:= goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println("请求登陆页失败")
		return ""
	}
	viewState,exits:=result.Find("#form1").Find("input").Eq(0).Attr("value")
	if !exits {
		fmt.Println("__VIEWSTATE获取失败")
		return ""
	}
	log.Println(viewState)
	return viewState
}
