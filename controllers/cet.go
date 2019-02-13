package controllers

import (
	"net/http"
	"fmt"
	cet "../server"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"math/rand"
)


//  返回四六级查询页面
func (c *MainController) Cet() {
	c.TplName = "cet.html"
}

func (c *MainController) GetCodeUrl() {
	jar, _ := cookiejar.New(nil)
	ticket := c.GetString("ticket")
	querier := cet.NewQuerier()
	res, _ := querier.GetCodeUrl(nil, ticket)
	u := &url.URL{
		Scheme: "http",
		Host:   "cache.neea.edu.cn",
		Path:   "/Imgs.do",
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return
	}
	uq := req.URL.Query()
	uq.Set("c", "CET")
	uq.Set("ik", ticket)
	uq.Set("t", strconv.FormatFloat(rand.Float64(), 'E', -1, 64))
	req.URL.RawQuery = uq.Encode()
	req.Header.Set("Referer", "http://cet.neea.edu.cn/cet/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	req.Header.Set("Host", "neea.edu.cn")
	tr := &http.Transport{}
	client := http.Client{
		Transport: tr,
		Jar:       jar,
	}
	resp, _ := client.Do(req)
	res, _ = querier.ParseCodeUrl(resp)
	c.Ctx.Output.Cookie("BIGipServercache.neea.edu.cn_pool", "580962314.39455.0000")
	c.Ctx.Output.Body([]byte(res))
	fmt.Println(res)
}

func (c *MainController) GetCetGrade() {
	ticket := c.GetString("ticket")
	name := c.GetString("name")
	code := c.GetString("yzm")
	querier := cet.NewQuerier()
	var err1 error
	cookie := make([]*http.Cookie, 1)
	cookie[0], err1 = c.Ctx.Request.Cookie("BIGipServercache.neea.edu.cn_pool")
	cookie[0].Value = "530630666.39455.0000"
	if err1 != nil {
		return
	}
	result, err := querier.Query(cookie[0], nil, ticket, name, code)
	if err != nil {
		panic(err)
	}
	fmt.Println("姓名：", result.Name)
	fmt.Println("学校：", result.University)
	fmt.Println("考试级别：", result.Level)
	fmt.Println("笔试准考证号：", result.WrittenTicket)
	fmt.Println("总分：", result.Score)
	fmt.Println("听力：", result.Listening)
	fmt.Println("阅读：", result.Reading)
	fmt.Println("写作和翻译：", result.WritingTranslation)
	fmt.Println("口试准考证号：", result.OralTicket)
	fmt.Println("口试等级：", result.OralLevel)

}
