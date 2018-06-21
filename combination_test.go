package smile

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCombination(t *testing.T) {
	w := httptest.NewRecorder()
	gz := gzip.NewWriter(w)
	defer gz.Close()

	post := url.Values{}
	post.Add("key1", "string1")
	post.Add("key1", "string2")
	post.Add("test", "test1")

	r := httptest.NewRequest("POST", "http://localhost/test?id=1", strings.NewReader(post.Encode()))
	r.Header["Content-Type"] = []string{"application/form-data;charset=utf-8"}
	r.PostForm = post
	r.Header["User-Agent"] = []string{"QQBrowser-1.0"}
	r.AddCookie(&http.Cookie{
		Name:     "testing",
		Value:    "hello world",
		Path:     "/",
		Domain:   "localhost",
		HttpOnly: false,
		Secure:   false,
	})

	e := Default()
	c := InitCombination(w, r, e)

	ip := c.GetClientIP()
	fmt.Println("IP:", ip)
	cookie, _ := c.GetCookie("testing")
	fmt.Println("cookie:", cookie)
	header := c.GetHeader()
	fmt.Println("header:", header)
	schema := c.GetScheme()
	fmt.Println("schema:", schema)
	proto := c.GetProto()
	fmt.Println("proto:", proto)
	uri := c.GetURL()
	fmt.Println("url:", uri)
	host := c.GetHost()
	fmt.Println("host", host)
	method := c.GetMethod()
	fmt.Println("method:", method)
	path := c.GetPath()
	fmt.Println("path", path)
	postKey := c.GetPostParam("test")
	fmt.Println("postKey-test:", postKey)
	postParams := c.GetMultipartFormParam("key1")
	fmt.Println("postParams-key1:", postParams)
	postFile := c.GetMultipartFormFile("key1")
	fmt.Println("post-Fil:", postFile)
	id := c.GetQueryParam("id")
	fmt.Println("queryParam-id:", id)
	queryS := c.GetQueryString()
	fmt.Println("queryString:", queryS)

	ua := c.GetUserAgent()
	fmt.Println("ua:", ua)
	rb := c.GetRawBody()
	fmt.Println("rb:", rb)

	c.SetHeader("resp", "testSetHeader")
	c.SetCookie(&http.Cookie{
		Name:     "SetHeader",
		Value:    "hello world",
		Path:     "/",
		Domain:   "localhost",
		HttpOnly: false,
		Secure:   false,
	})

}
