package controllers

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"
)

var globalSessions *session.Manager

//验证码url
var CheckCodeUrl = "http://202.116.160.170/CheckCode.aspx"

func init() {
	sessionConfig := &session.ManagerConfig{
		CookieName:      "gosessionid",
		EnableSetCookie: true,
		Gclifetime:      3600,
		Maxlifetime:     3600,
		Secure:          false,
		CookieLifeTime:  3600,
		ProviderConfig:  "./tmp",
	}
	globalSessions, _ = session.NewManager("memory", sessionConfig)
	go globalSessions.GC()
}

type MainController struct {
	beego.Controller
}

func GetSession() *session.Manager {
	return globalSessions
}

// 返回界面
func (c *MainController) Login() {
	c.TplName = "login.html"
}

func (c *MainController) ToCredit() {
	c.TplName = "toCredit.html"
}

//获取验证码
func (c *MainController) CheckCode() {
	jar, _ := cookiejar.New(nil)
	checkCodeUrl, _ := url.Parse(CheckCodeUrl)
	client := http.Client{
		Jar: jar,
	}
	req, err := client.Get(CheckCodeUrl)
	if err != nil {
		log.Println(err)
		c.TplName = "fault.html"
		return
	}
	cook := client.Jar.Cookies(checkCodeUrl)
	c.Ctx.Output.Cookie(cook[0].Name, cook[0].Value)
	imageCode, _ := ioutil.ReadAll(req.Body)
	c.Ctx.Output.Body(imageCode)
}

// Craw 登陆函数
func (c *MainController) Craw() {
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

	if len(c.Ctx.Request.Form["flag"]) == 0 {
		c.TplName = "fault.html"
		return
	}
	flag:=c.Ctx.Request.Form["flag"][0]
	if len(c.Ctx.Request.Form["username"]) == 0 || len(c.Ctx.Request.Form["password"]) == 0 || len(c.Ctx.Request.Form["yzm"]) == 0 {
		c.TplName = "fault.html"
		return
	}

	//获取主页
	url1 := "http://202.116.160.170/default2.aspx"
	v := url.Values{}
	encoder := mahonia.NewEncoder("gbk")
	decoder := mahonia.NewDecoder("gbk")
	but := encoder.ConvertString("学生")
	v.Add("__VIEWSTATE", "dDwxNTMxMDk5Mzc0Ozs+SG77OyOzybTqnYFhiAU0smy8ot4=")
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

	//解析主页
	doc := decoder.NewReader(response.Body)
	result, err := goquery.NewDocumentFromReader(doc)
	if err != nil {
		log.Println(err)
		c.TplName = "fault.html"
		return
	}

	username := c.Ctx.Request.Form["username"][0]
	log.Println(username, "登陆-主页获取成功", response.Status)
	cname := result.Find("#xhxm").Text()
	sess, _ := globalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	defer sess.SessionRelease(c.Ctx.ResponseWriter)
	sess.Set("username", username)
	cname = strings.TrimRight(cname, "同学")
	encoder = mahonia.NewEncoder("gbk")
	cname = encoder.ConvertString(cname)
	sess.Set("cname", cname)
	client.Get("https://sc.ftqq.com/SCU20914Teefb444fcce3027f14828723ca1cd65e5a6c2b88500ab.send?text=" +
		url.QueryEscape(username+" "+cname+" 登陆"))
	c.Data["Name"] = cname
	c.Data["Num"] = username
	if flag =="1"{
		c.Ctx.Redirect(302, "/school/toCredit")
		return
	}else if flag =="2"{
		c.Ctx.Redirect(302, "/school/toEvaluate")
		return
	}
	c.TplName = "fault.html"
}
