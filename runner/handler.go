package runner

import (
	"Pricking/common"
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"mime"
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
	excludeFileTypes []string
	saveFileTypes []string
	saveFileDir string
	ConfigFile string
	LoggerFile string
	config     *viper.Viper
	staticDir         string
	PrickingPrefixUrl string
	Logger            *log.Logger
	request *http.Request
	Plugins * HandlerPlugins
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
	handler.excludeFileTypes = handler.config.GetStringSlice("exclude_file_types")
	handler.staticDir = handler.config.GetString("static_dir")
	handler.PrickingPrefixUrl = handler.config.GetString("pricking_prefix_url")
	handler.ListenAddress = handler.config.GetString("listen_address")
	handler.saveFileTypes = handler.config.GetStringSlice("save_file_types")
	handler.saveFileDir = handler.config.GetString("save_file_dir")
	// handler.LoadPlugins()
}
//
// func  (handler *Handler) LoadPlugins(){
// 	handler.Plugins = new(HandlerPlugins)
// 	plist := handler.config.GetStringMap("plugins")
// 	for pluginName,plugin := range plist{
//
// 	}
// }

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
		location := w.Header.Get("Location")
		// log.Println("Location ... ",location)
		if strings.HasPrefix(location,"http") || strings.HasPrefix(location,"https"){
			locationUrl , err := url.Parse(location)
			if err != nil{
				w.Header.Set("Location","/")
			}else{
				w.Header.Set("Location",locationUrl.RequestURI())
			}
		}
	}

	// gzip 解压
	if w.Header.Get("Content-Encoding") == "gzip" {
		w.Header.Del("Content-Encoding")
		respBodyByte = common.UnGzipData(respBodyByte)
	}

	// 移除内容安全策略
	w.Header.Del("Content-Security-Policy-Report-Only")
	w.Header.Del("Content-Security-Policy")
	w.Header.Set("Access-Control-Allow-Origin","*")
	w.Header.Del("Content-Security-Policy")
	w.Header.Set("Access-Control-Allow-Credentials","true")
	contentType := strings.Trim(w.Header.Get("Content-Type")," ")

	for _,ft := range handler.filterTypes{
		if strings.HasPrefix(contentType,ft) {
			// 如果匹配到Content-Type，注入内容
			respBodyByte = bytes.Replace(respBodyByte, []byte("</body>"), []byte(handler.injectBody + "</body>") , -1)
		}
	}

	w.Body = ioutil.NopCloser(bytes.NewReader(respBodyByte))
	w.ContentLength = int64(len(respBodyByte))
	w.Header.Set("Content-Length", strconv.Itoa(len(respBodyByte)))
	return nil
}

func (handler *Handler)saveResponseFile(w *http.Response)  {
	var filename string
	// 保存特定类型文件
	fileType := w.Header.Get("Content-Type")

	for _, saveType := range handler.saveFileTypes {
		if fileType == saveType {
			ext, _ := mime.ExtensionsByType(w.Header.Get("Content-Type"))
			if len(ext)> 0 {
				filename = handler.saveFileDir + common.UUID() + ext[0]
			}else{
				filename = handler.saveFileDir + common.UUID()
			}
			respBodyByte, _ := ioutil.ReadAll(w.Body)
			common.SaveFile(filename,respBodyByte)
			w.Body = ioutil.NopCloser(bytes.NewReader(respBodyByte))
		}
	}
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
	for _, staticType := range handler.excludeFileTypes {
		staticType = strings.ToUpper(staticType)
		upperPath := strings.ToUpper(r.URL.Path)
		if strings.HasSuffix(upperPath, staticType) {
			return
		}
	}

	dump, err := httputil.DumpRequest(r, true)
	if err == nil{
		handler.Logger.Println(r.RemoteAddr + "\n" + string(dump) + "\n")
	}else{
		log.Fatalln(err)
	}

}

func (handler * Handler)modifyHost(r *http.Request,header string ,host string)  {
	r.Header.Set(header,host)
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.request = r
	remote, err := url.Parse(handler.Url)
	if err != nil {
		panic(err)
	}

	if r.TLS != nil {
		r.URL.Scheme = "https"
	}else{
		r.URL.Scheme = "http"
	}

	// 覆盖Host头
	r.URL.Host = remote.Host
	r.Host = remote.Host

	handler.modifyHost(r,"Referer",handler.Url)
	handler.modifyHost(r,"Origin",handler.Url)

	if handler.prickingResponse(w, r) {
		return
	}

	handler.loggingRequest(r)
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ModifyResponse = handler.modifyResponse

	proxy.ServeHTTP(w, r)
}

