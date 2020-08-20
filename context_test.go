package smile

import (
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestContext(t *testing.T) {
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
	c := initContext(w, r, e)

	ip := c.GetClientIP()
	t.Log("IP:", ip)
	cookie, _ := c.GetCookie("testing")
	t.Log("cookie:", cookie)
	header := c.GetHeader()
	t.Log("header:", header)
	schema := c.GetScheme()
	t.Log("schema:", schema)
	proto := c.GetProto()
	t.Log("proto:", proto)
	uri := c.GetURL()
	t.Log("url:", uri)
	host := c.GetHost()
	t.Log("host", host)
	method := c.GetMethod()
	t.Log("method:", method)
	path := c.GetPath()
	t.Log("path", path)
	postKey := c.GetPostParam("test")
	t.Log("postKey-test:", postKey)
	postParams := c.GetMultipartFormParam("key1")
	t.Log("postParams-key1:", postParams)
	postFile := c.GetMultipartFormFile("key1")
	t.Log("post-Fil:", postFile)
	id := c.GetQueryParam("id")
	t.Log("queryParam-id:", id)
	queryS := c.GetQueryString()
	t.Log("queryString:", queryS)

	ua := c.GetUserAgent()
	t.Log("ua:", ua)
	rb := c.GetRawBody()
	t.Log("rb:", rb)

	c.WriteString("testing string")
	c.Flush()
	c.Done()

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

func TestHandlerChain(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost/test?id=1", nil)

	c := initContext(w, r, Default())
	hc := newHandlerChain()
	n := 1
	for i := 1; i < 6; i++ {
		hc.add(func(c *Context) error {
			logi := n
			t.Log(logi)
			if n == 2 {
				t.Log(c.Abort())
			}
			n++
			return nil
		})
	}
	c.handlerChain = hc

	if err := c.Next(); err != nil {
		t.Error(err.Error())
	}
	hc.reset()
	if err := hc.next(c); err != nil {
		t.Error(err.Error())
	}
}
