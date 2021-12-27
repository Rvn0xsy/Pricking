package main

import (
	"Pricking/runner"
	"flag"
	"log"
	"net/http"
	"os"
)
var handler runner.Handler

func init() {
	flag.StringVar(&handler.Url, "url", "https://payloads.online", "pricking url")
	flag.StringVar(&handler.PrickingPrefixUrl, "route", "/pricking_static_files", "pricking js route path")
	flag.StringVar(&handler.ConfigFile, "config", "./config.yaml", "pricking js route path")
	flag.StringVar(&handler.LoggerFile, "log", "./access.txt", "pricking logging file")
	flag.Parse()
}

func main() {
	handler.LoadConfig()
	logFile, err := os.OpenFile(handler.LoggerFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	// multiWriter := io.MultiWriter(os.Stdout, logFile)
	handler.Logger = log.New(logFile, "", 1)
	err = http.ListenAndServe(handler.ListenAddress, &handler)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}
