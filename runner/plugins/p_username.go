package plugins

import (
	"log"
	"net/http"
)



type PluginUsername struct {
	Enable bool
	Name string
}

func (handler * PluginUsername)HandlerRequest(path string,r *http.Request, w *http.ResponseWriter)  {
	log.Println("[+]HandlerRequest ", r.RemoteAddr)
}
