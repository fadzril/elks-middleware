package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

// Constant Properties
const (
	ElasticServer = "c4t17650.itcs.hpecorp.net"
	ElasticPort   = 9300
	ElasticIndex  = "leo"
)

var (
	// Es Elastic
	Es ElasticService
)

func init() {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	Es.Client = &http.Client{Transport: transport}
}

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultCommonStack...)

	router, err := rest.MakeRouter(
		//rest.Get("/status", Es.getElasticStatus),

		rest.Get("/log", GetSubscription),
		rest.Get("/log/category", GetSubscriptionByCategory),
		rest.Get("/log/tags", GetSubscriptionByTags),
	)

	if err != nil {
		log.Fatal("Error Serving Routers\n", err)
	}

	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8000", api.MakeHandler()))
}
