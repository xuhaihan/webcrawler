package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
)

//查询学分
func (c *MainController) QueryCredit() {
	//初始化client
	jar, _ := cookiejar.New(nil)
	checkCodeUrl, _ := url.Parse(CheckCodeUrl)
	client := http.Client{
		Jar:     jar,
		Timeout: time.Second * 10,
	}
	var err error
	cookie := make([]*http.Cookie, 1)
	cookie[0], err = c.Ctx.Request.Cookie("ASP.NET_SessionId")
	if err != nil {
		c.TplName = "fault.html"
		return
	}
	client.Jar.SetCookies(checkCodeUrl, cookie)
	c.Ctx.Request.ParseForm()

	//获取查询页
	encoder := mahonia.NewEncoder("gbk")
	decoder := mahonia.NewDecoder("gbk")
	sess, _ := globalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	defer sess.SessionRelease(c.Ctx.ResponseWriter)
	username := sess.Get("username").(string)
	cname := sess.Get("cname").(string)
	resultUrl := "http://202.116.160.170/xscjcx.aspx?xh=" + username + "&xm=" + url.QueryEscape(cname) + "&gnmkdm=N121605"
	req, _ := http.NewRequest("GET", resultUrl, nil)
	req.Header.Add("Referer", "http://202.116.160.170/xs_main.aspx?xh="+username)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	response, err := client.Do(req)
	if err != nil {
		c.TplName = "fault.html"
		return
	}

	log.Println(username, "学分查询页", response.Status)
	if response.StatusCode != 200 {
		c.TplName = "fault.html"
		return
	}

	//获取view，event
	doc := decoder.NewReader(response.Body)
	result, _ := goquery.NewDocumentFromReader(doc)
	view, exits := result.Find("#Form1").Find("input").Eq(2).Attr("value")
	if !exits {
		fmt.Println("__VIEWSTATE获取失败")
		c.TplName = "toevaluate.html"
		return

	}
	v := url.Values{}
	v.Add("Button1", encoder.ConvertString("成绩统计"))
	v.Add("__EVENTTARGET", "")
	v.Add("__EVENTARGUMENT", "")
	v.Add("__VIEWSTATE", view)
	v.Add("hidLanguage", "")
	v.Add("ddlXN", "")
	v.Add("ddlXQ", "")
	v.Add("ddl_kcxz", "")

	//构建新请求
	body := strings.NewReader(v.Encode())
	req, err = http.NewRequest("POST", resultUrl, body)
	if err != nil {
		log.Println(err)
		c.TplName = "fault.html"
		return
	}
	req.Header.Add("Referer", resultUrl)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	response, err = client.Do(req)
	if err != nil {
		log.Println(err)
		c.TplName = "fault.html"
		return
	}
	log.Println(username, "学分结果页", response.Status)
	ma, xf, jd := matchcredit(response)
	c.Data["Name"] = decoder.ConvertString(cname)
	c.Data["Num"] = username
	c.Data["Xf"] = xf
	c.Data["Jd"] = jd
	c.Data["Res"] = ma
	c.TplName = "credit.html"
}
