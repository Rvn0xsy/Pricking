package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type Handler struct {
	url string
	listenAddress string
	injectBody string
	filterTypes []string
	excludeFile []string
	configFile string
	loggerFile string
	config     * viper.Viper
	staticDir string
	prickingPrefixUrl string
	logger *log.Logger
}

var handler  Handler


func (this * Handler)LoadConfig()  {
	this.config = viper.New()
	this.config.SetConfigName("config")
	this.config.SetConfigType("yaml")
	this.config.AddConfigPath("/etc/pricking/")
	this.config.AddConfigPath("$HOME/.pricking")
	this.config.AddConfigPath(".")
	if this.configFile != "" {
		log.Println("Loading config ", this.configFile)
		content, err := ioutil.ReadFile(this.configFile)
		if err != nil{
			panic(err)
		}
		if err = this.config.ReadConfig(bytes.NewReader(content)); err != nil{
			panic(err)
		}
	}else{
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

func (this *Handler)modifyResponse(w *http.Response) error{
	contents, _ := ioutil.ReadAll(w.Body)
	var isModify = false
	for _,t := range this.filterTypes{
		if strings.Index(w.Header.Get("Content-Type"),t) == 0{
			log.Println("Modify Body ...")
			isModify = true
		}
	}
	var newBody io.ReadCloser
	if isModify{
		log.Println("Inject Body : ",this.injectBody)
		w.ContentLength = w.ContentLength + int64(len(this.injectBody))
		newBody = ioutil.NopCloser(strings.NewReader(string(contents)+ this.injectBody))
	}else{
		newBody = ioutil.NopCloser(strings.NewReader(string(contents)))
	}
	w.Body = newBody
	return nil
}

func (this *Handler)prickingResponse(w http.ResponseWriter,r *http.Request){
	if strings.HasPrefix(r.URL.Path, this.prickingPrefixUrl) {
		file := this.staticDir + r.URL.Path[len(this.prickingPrefixUrl):]
		http.ServeFile(w, r, file)
		return
	}
}

func (this *Handler) loggingRequest(r *http.Request){
	// 对一般静态文件进行忽略处理
	for _,staticType := range this.excludeFile{
		staticType = strings.ToUpper(staticType)
		upperPath := strings.ToUpper(r.URL.Path)
		if strings.HasSuffix(upperPath,staticType) {
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
	buffer.WriteString(r.Method +" " + r.URL.String() +newLine)
	for header,value := range r.Header{
		buffer.WriteString(header + ": " + strings.Join(value,",") + newLine)
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

	this.prickingResponse(w,r)

	r.Host = remote.Host  // 覆盖Host头

	this.loggingRequest(r)

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ModifyResponse = this.modifyResponse
	proxy.ServeHTTP(w, r)
}

func init() {
	flag.StringVar(&handler.url, "url", "https://payloads.online","pricking url")
	flag.StringVar(&handler.prickingPrefixUrl, "route", "/pricking_static_files","pricking js route path")
	flag.StringVar(&handler.configFile, "config", "./config.yaml","pricking js route path")
	flag.StringVar(&handler.loggerFile, "log", "./access.txt","pricking logging file")
	flag.Parse()
}

func main() {
	handler.LoadConfig()
	logFile, err := os.OpenFile(handler.loggerFile, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	multiWriter := io.MultiWriter(os.Stdout,logFile)
	handler.logger = log.New(multiWriter,"",1)
	err = http.ListenAndServe(handler.listenAddress, &handler)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}
