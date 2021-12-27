package runner

import (
	"Pricking/common"
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

type Handler struct {
	Url           string
	ListenAddress string
	injectBody    string
	filterTypes       []string
	excludeFile []string
	ConfigFile string
	LoggerFile string
	config     *viper.Viper
	staticDir         string
	PrickingPrefixUrl string
	Logger            *log.Logger
}



func (handler *Handler) LoadConfig() {
	handler.config = viper.New()
	handler.config.SetConfigName("config")
	handler.config.SetConfigType("yaml")
	handler.config.AddConfigPath("/etc/pricking/")
	handler.config.AddConfigPath("$HOME/.pricking")
	handler.config.AddConfigPath(".")
	if handler.ConfigFile != "" {
		log.Println("Loading config ", handler.ConfigFile)
		content, err := ioutil.ReadFile(handler.ConfigFile)
		if err != nil {
			panic(err)
		}
		if err = handler.config.ReadConfig(bytes.NewReader(content)); err != nil {
			panic(err)
		}
	} else {
		err := handler.config.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}
	}
	handler.injectBody = handler.config.GetString("inject_body")
	handler.filterTypes = handler.config.GetStringSlice("filter_type")
	handler.excludeFile = handler.config.GetStringSlice("exclude_file")
	handler.staticDir = handler.config.GetString("static_dir")
	handler.PrickingPrefixUrl = handler.config.GetString("pricking_prefix_url")
	handler.ListenAddress = handler.config.GetString("listen_address")
}


func (handler *Handler) modifyResponse(w *http.Response) error {
	respBodyByte, err := ioutil.ReadAll(w.Body)
	if err != nil {
		return err
	}
	err = w.Body.Close()
	if err != nil {
		return err
	}

	// 让客户端跟随跳转
	if w.StatusCode == 302 {
		replacer := strings.NewReplacer(handler.Url,"")
		newLocation := replacer.Replace(w.Header.Get("Location"))
		w.Header.Set("Location",newLocation)
	}

	// gzip 解压
	if w.Header.Get("Content-Encoding") == "gzip" {
		w.Header.Del("Content-Encoding")
		respBodyByte = common.UnGzipData(respBodyByte)
	}

	// 移除内容安全策略
	w.Header.Del("Content-Security-Policy-Report-Only")
	w.Header.Del("Content-Security-Policy")
	respBodyByte = bytes.Replace(respBodyByte, []byte("Content-Security-Policy"), []byte( "") , -1)
	respBodyByte = bytes.Replace(respBodyByte, []byte("content-security-policy"), []byte( "") , -1)

	// 注入内容
	respBodyByte = bytes.Replace(respBodyByte, []byte("</body>"), []byte(handler.injectBody + "</body>") , -1)

	w.Body = ioutil.NopCloser(bytes.NewReader(respBodyByte))
	w.ContentLength = int64(len(respBodyByte))
	w.Header.Set("Content-Length", strconv.Itoa(len(respBodyByte)))
	return nil
}

func (handler *Handler) prickingResponse(w http.ResponseWriter, r *http.Request) bool {
	if strings.HasPrefix(r.URL.Path, handler.PrickingPrefixUrl) {
		file := handler.staticDir + r.URL.Path[len(handler.PrickingPrefixUrl):]
		http.ServeFile(w, r, file)
		return true
	}
	return false
}

func (handler *Handler) loggingRequest(r *http.Request) {
	// 对一般静态文件进行忽略处理
	for _, staticType := range handler.excludeFile {
		staticType = strings.ToUpper(staticType)
		upperPath := strings.ToUpper(r.URL.Path)
		if strings.HasSuffix(upperPath, staticType) {
			return
		}
	}

	var bodyBytes []byte
	var newLine = "\r\n"
	if r.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(r.Body)
	}
	// Copy Request Body
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	buffer := new(bytes.Buffer)
	buffer.WriteString(r.RemoteAddr + newLine)
	buffer.WriteString(r.Method + " " + r.URL.String() + newLine)
	for header, value := range r.Header {
		buffer.WriteString(header + ": " + strings.Join(value, ",") + newLine)
	}
	buffer.WriteString(newLine)
	buffer.Write(bodyBytes)
	handler.Logger.Println(buffer.String())
}

func (handler * Handler)modifyHost(r *http.Request,header string ,host string)  {
	Url, err := url.Parse(r.Header.Get(header))
	if  err != nil{
		return
	}
	Url.Host = host
	r.Header.Set(header, Url.String())
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remote, err := url.Parse(handler.Url)
	if err != nil {
		panic(err)
	}
	r.Host = remote.Host // 覆盖Host头

	handler.modifyHost(r,"Referer",remote.Host)
	handler.modifyHost(r,"Origin",remote.Host)

	if handler.prickingResponse(w, r) {
		return
	}

	handler.loggingRequest(r)
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ModifyResponse = handler.modifyResponse

	proxy.ServeHTTP(w, r)
}

