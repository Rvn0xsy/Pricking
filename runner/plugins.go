package runner

import (
	"net/http"
)

type HandlerPlugins interface {
	HandlerRequest(path string,r *http.Request) bool
}

type Plugins struct {
	PluginList map[string] HandlerPlugins
}

// 初始化插件
func (p *Plugins)init(){
	p.PluginList = make(map[string]HandlerPlugins)
}

// 注册插件
func (p *Plugins)register(name string, plugin HandlerPlugins) {
	p.PluginList[name] = plugin
}
