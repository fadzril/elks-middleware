package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/mattbaird/elastigo/lib"
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
		//rest.Get("/log/:type", Es.getRelatedLog),
		//rest.Get("/log/subscription/:id", Es.getElasticRelatedLog),
		rest.Get("/log", GetSubscription),
		rest.Get("/log/category", GetSubscriptionByCategory),
	)

	if err != nil {
		log.Fatal("Error Serving Routers\n", err)
	}

	api.SetApp(router)

	/*
		api.SetApp(rest.AppSimple(func(w rest.ResponseWriter, r *rest.Request) {
			w.WriteJson(map[string]string{"result": "Hello Worlds"})
		}))
	*/

	log.Fatal(http.ListenAndServe(":8000", api.MakeHandler()))
}

// ElasticService : Struct
type ElasticService struct {
	Client *http.Client
}

// "Status: ElasticService Status"
func (e ElasticService) getElasticStatus(w rest.ResponseWriter, rr *rest.Request) {
	uri := fmt.Sprintf("https://%s/%s/_status", ElasticServer, ElasticIndex)
	response, error := e.Client.Get(uri)

	log.Println("GET: ", uri)
	if error != nil {
		log.Println("Error: ", error)
	}

	defer response.Body.Close()

	responseByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("fail to read response data")
		return
	}
	w.WriteJson(string(responseByte))
}

func (e ElasticService) getRelatedLog(w rest.ResponseWriter, rr *rest.Request) {
	p := rr.PathParam("type")
	uri := fmt.Sprintf("https://%s/%s/%s/_search?pretty=true", ElasticServer, ElasticIndex, p)

	log.Println("GET: ", uri)
	response, error := e.Client.Get(uri)

	if error != nil {
		log.Println("Error: ", error)
	}

	defer response.Body.Close()

	responseByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("fail to read")
		return
	}
	w.WriteJson(string(responseByte))
}

func (e ElasticService) getElasticRelatedLog(w rest.ResponseWriter, rr *rest.Request) {
	api := elastigo.NewConn()
	api.Domain = "c4t03459.itcs.hpecorp.net"

	uid := rr.PathParam("id")
	params := fmt.Sprintf("{\"query\":{\"term\":{\"subscription_id\":\"%s\"}}}", uid)

	out, err := api.Search("leo", "tibco", nil, params)
	if err != nil {
		log.Fatal(err)
	}

	if len(out.Hits.Hits) >= 1 {
		res := &Subscriptions{
			Total: len(out.Hits.Hits),
		}
		for _, item := range out.Hits.Hits {
			t, e := json.Marshal(item.Source)
			if e != nil {
				log.Println("Error marshalling", e)
			}
			var sub Subscription
			ers := json.Unmarshal(t, &sub)
			if ers != nil {
				log.Println("Error Unmarshalling", ers)
			}
			res.Result = append(res.Result, sub)
		}
		w.WriteJson(res)
	} else {
		log.Println("OUT.HITS\n", out.Hits)
	}
}
