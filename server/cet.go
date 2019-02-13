package server

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"io/ioutil"
	"encoding/json"
	"net"
	"os"
	"bufio"
	"io"
	"path/filepath"
)

var (
	query Querier = &querier{}
)

// Query 根据准考证号、姓名查询四六级成绩
func Query(cookie *http.Cookie,ctx context.Context, ticket, name string,code string) (*Result, error) {
	return query.Query(cookie,ctx, ticket, name,code)
}

// Querier 提供四六级成绩查询的接口
type Querier interface {
	// 根据准考证号、姓名查询
	Query(cookie *http.Cookie,ctx context.Context, ticket, name string,code string) (*Result, error)
	GetCodeUrl(ctx context.Context, ticket string) (string, error)
	SaveImage(url string,filename string) string
	Parse(resp *http.Response) (result *Result, err error)
	ParseCodeUrl(resp *http.Response) (result string, err error)

}

// Result 查询结果
type Result struct {
	Name               string  // 姓名
	University         string  // 学校
	Level              string  // 考试级别
	WrittenTicket      string  // 笔试准考证号
	Score              float32 // 总分
	Listening          float32 // 听力
	Reading            float32 // 阅读
	WritingTranslation float32 // 写作和翻译
	OralTicket         string  // 口试准考证号
	OralLevel          string  // 口试等级
}

const (
	// 四六级查询的网址
	cetImg="/static/code/"
	cetHost = "cache.neea.edu.cn"
)

var (
	// ErrNotFound not found
	ErrNotFound = errors.New("not found")
)

// NewQuerier 创建Querier接口
func NewQuerier() Querier {
	return &querier{}
}

type querier struct {
}

func (q *querier) Query(cookie *http.Cookie,ctx context.Context, ticket, name string,code string) (result *Result, err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	u := &url.URL{
		Scheme: "http",
		Host:   cetHost,
		Path:   "/cet/query",
	}
	var r http.Request
	r.ParseForm()
	r.Form.Add("data","CET6_181_DANGCI,"+ticket+","+name)
	r.Form.Add("v",strings.Trim(code," "))
	bodyStr := strings.TrimSpace(r.Form.Encode())
	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(bodyStr))
	req.AddCookie(cookie)
	if err != nil {
		return
	}

	req.Header.Set("Referer", fmt.Sprintf("http://cet.neea.edu.cn/cet/query",))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
	req.Header.Set("X-Forwarded-For", q.randomIP())
	req.Header.Set("Access-Control-Allow-Origin","http://cet.neea.edu.cn")
	req.Header.Set("Access-Control-Allow-Headers","X-Requested-With")

	err = q.httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		r, err := q.Parse(resp)
		if err != nil {
			return err
		}
		result = r
		return nil
	})
	return
}

func (q *querier) GetCodeUrl(ctx context.Context, ticket string) (result string, err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	u := &url.URL{
		Scheme: "http",
		Host:   cetHost,
		Path:   "/Imgs.do",
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return
	}
	uq := req.URL.Query()
	uq.Set("c", "CET")
	uq.Set("ik", ticket)
	uq.Set("t", strconv.FormatFloat(rand.Float64(),'E',-1,64))
	req.URL.RawQuery = uq.Encode()

	req.Header.Set("Referer","http://cet.neea.edu.cn/cet/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	req.Header.Set("Host", "neea.edu.cn")
	req.Header.Set("X-Forwarded-For",q.randomIP())

	err = q.httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		r, err := q.ParseCodeUrl(resp)
		if err != nil {
			return err
		}
		result = r
		return nil
	})
	return
}

var (
	rd = rand.New(rand.NewSource(math.MaxInt64))
)

func (q *querier) randomIP() string {
	return fmt.Sprintf("%d.%d.%d.%d", rd.Intn(255), rd.Intn(255), rd.Intn(255), rd.Intn(255))
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func(q *querier)SaveImage(url string,filename string) string{
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(resp.Body, 32 * 1024)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	dir=strings.Replace(dir, "\\", "/", -1)
	fatherDir:=substr(dir, 0, strings.LastIndex(dir, "/"))
	str:=cetImg+filename
	result:=fatherDir+str
	result=strings.Replace(result,"\\","/",-1)
	fmt.Println(result)
	file, err := os.Create(result)
	if err != nil {
		panic(err)
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)
	_,err = io.Copy(writer, reader)
	return str
}

func (q *querier) httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	tr := &http.Transport{}
	client := &http.Client{
		Transport: tr,
	}
	c := make(chan error, 1)
	go func() { c <- f(client.Do(req)) }()
	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-c
		return ctx.Err()
	case err := <-c:
		return err
	}
}

func (q *querier) ParseCodeUrl(resp *http.Response) (result string, err error) {
	body, err := ioutil.ReadAll(resp.Body)
	start := strings.Index(string(body), "(")
	end := strings.Index(string(body), ")")
	result = string(body)[start+2 : end-1]
	return strings.Replace(result, "\"", "", 1), nil
}
func (q *querier) Parse(resp *http.Response) (result *Result, err error) {
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	start:=strings.Index(string(body),"{")
	end:=strings.Index(string(body),"}")
	str:=string(body)[start+1:end-1]
	myMap:=make(map[string]string)
	json.Unmarshal([]byte(str),&myMap)
	fmt.Println(myMap)
    return result,nil
}

func(q *querier) GetIp() string{
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//fmt.Println(q.GetIp())
				return ipnet.IP.String()
			}
		}
	}
	return ""
}