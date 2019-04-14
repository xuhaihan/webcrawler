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

func (c *MainController) ToGrade() {
	c.TplName = "toGrade.html"
}

//查询成绩
func (c *MainController) QueryGrade() {
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
	if len(c.Ctx.Request.Form["username"]) == 0 || len(c.Ctx.Request.Form["password"]) == 0 || len(c.Ctx.Request.Form["yzm"]) == 0 {
		c.TplName = "fault.html"
		return
	}
	//获取登陆主页
	url1 := "http://202.116.160.170/default2.aspx"
	v := url.Values{}
	encoder := mahonia.NewEncoder("gbk")
	decoder := mahonia.NewDecoder("gbk")
	but := encoder.ConvertString("学生")
	view:=getViewState("http://202.116.160.170")
	v.Add("__VIEWSTATE", view)
	v.Add("txtUserName", c.Ctx.Request.Form["username"][0])
	v.Add("TextBox2", c.Ctx.Request.Form["password"][0])
	v.Add("txtSecretCode", c.Ctx.Request.Form["yzm"][0])
	v.Add("RadioButtonList1", but)
	v.Add("Button1", "")
	v.Add("lbLanguage", "")
	v.Add("hidPdrs", "")
	v.Add("hidsc", "")
	v.Add("__EVENTVALIDATION", "/wEWDgKX/4yyDQKl1bKzCQLs0fbZDAKEs66uBwK/wuqQDgKAqenNDQLN7c0VAuaMg+INAveMotMNAoznisYGArursYYIAt+RzN8IApObsvIHArWNqOoPqeRyuQR+OEZezxvi70FKdYMjxzk=")

	//建立client发送POST请求
	body := strings.NewReader(v.Encode())
	r, _ := http.NewRequest("POST", url1, body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Referer", "http://202.116.160.170/default2.aspx")
	r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	response, err := client.Do(r)
	if err != nil {
		log.Println(err)
		c.TplName = "fault.html"
		return
	}

	//解析主页，如果有欢迎则说明获取失败
	doc := decoder.NewReader(response.Body)
	result, err := goquery.NewDocumentFromReader(doc)
	if err != nil {
		log.Println(err)
		c.TplName = "fault.html"
		return
	}
	log.Println(c.Ctx.Request.Form["username"][0], "登陆-主页获取成功", response.Status)

	username := c.Ctx.Request.Form["username"][0]
	cname := result.Find("#xhxm").Text()
	cname = strings.TrimRight(cname, "同学")
	sess, _ := globalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	sess.Set("username", username)
	sess.Set("cname", cname)
	client.Get("https://sc.ftqq.com/SCU20914Teefb444fcce3027f14828723ca1cd65e5a6c2b88500ab.send?text=" +
		url.QueryEscape(c.Ctx.Request.Form["username"][0]+" "+cname+" 登陆"))
	//获取查询成绩页
	encoder = mahonia.NewEncoder("utf-8")
	decoder = mahonia.NewDecoder("utf-8")
	resultUrl := "http://202.116.160.170/xscjcx.aspx?xh=" + c.Ctx.Request.Form["username"][0] + "&xm=" + url.QueryEscape(cname) + "&gnmkdm=N121605"
	req, _ := http.NewRequest("GET", resultUrl, nil)
	req.Header.Add("Referer", "http://202.116.160.170/xs_main.aspx?xh="+c.Ctx.Request.Form["username"][0])
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	response, err = client.Do(req)
	if err != nil {
		c.TplName = "fault.html"
		return
	}
	log.Println(username, "成绩查询页", response.Status)
	if response.StatusCode != 200 {
		c.TplName = "fault.html"
		return
	}
	doc = decoder.NewReader(response.Body)
	result, _ = goquery.NewDocumentFromReader(doc)
	view, exits := result.Find("#Form1").Find("input").Eq(2).Attr("value")
	if !exits {
		fmt.Println("__VIEWSTATE获取失败")
		c.TplName = "toEvaluate.html"
		return
	}

	event, _ := result.Find("#__EVENTVALIDATION").Attr("value")
	v = url.Values{}
	v.Add("btn_xq", encoder.ConvertString("学期成绩"))
	v.Add("__EVENTTARGET", "")
	v.Add("__EVENTARGUMENT", "")
	v.Add("__VIEWSTATE", view)
	v.Add("hidLanguage", "")
	v.Add("ddlXN", c.Ctx.Request.Form["year"][0])
	v.Add("ddlXQ", c.Ctx.Request.Form["term"][0])
	v.Add("ddl_kcxz", "")
	v.Add("__EVENTVALIDATION", event)

	//构造新请求
	body = strings.NewReader(v.Encode())
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

	log.Println(username, "成绩结果页", response.Status)

	grade := matchgrade(response)
	c.Data["Name"] =cname
	c.Data["Num"] = username
	c.Data["GradeResult"] = grade
	c.TplName = "grade.html"
}


