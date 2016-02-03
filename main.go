package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/zpatrick/go-config"
)

var (
	// Es Elastic
	Es                 ElasticService
	ElasticServerHost  = ""
	ElasticSearchIndex = ""
)

func initConfig() *config.Config {

	filename := flag.String("config", "config.ini", "Default properties/ini file")

	flag.Parse()
	log.Printf("Reading config file from: %s", *filename)
	if filename == nil {
		log.Fatal("Configuration files is mandatory before running this middleware")
		os.Exit(1)
	}

	iniFile := config.NewINIFile(*filename)
	return config.NewConfig([]config.Provider{iniFile})
}

func init() {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	Es.Client = &http.Client{Transport: transport}
}

func main() {

	cfg := initConfig()
	if err := cfg.Load(); err != nil {
		log.Println(err)
	}
	// Environment Config
	host, _ := cfg.String("elasticsearch.host")
	//port, _ := cfg.String("elasticsearch.port")
	indice, _ := cfg.String("elasticsearch.index")

	ElasticServerHost = string(fmt.Sprintf("%s", host))
	ElasticSearchIndex = string(fmt.Sprintf("%s", indice))

	log.Printf("Connecting to ElasticServerHost: %s", ElasticServerHost)

	api := rest.NewApi()
	api.Use(rest.DefaultCommonStack...)

	router, err := rest.MakeRouter(
		rest.Get("/rest/index", GetSubscription),
		rest.Get("/rest/category", GetSubscriptionByCategory),
		rest.Get("/rest/tags", GetSubscriptionByTags),
	)

	if err != nil {
		log.Fatal("Error Serving Routers\n", err)
	}

	api.SetApp(router)
	log.Println("Attemp to start server at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", api.MakeHandler()))
}
