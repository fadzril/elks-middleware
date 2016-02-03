package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

// ElasticService : Struct
type ElasticService struct {
	Client *http.Client
}

// "Status: ElasticService Status"
func (e ElasticService) getElasticStatus(w rest.ResponseWriter, r *rest.Request) {
	uri := fmt.Sprintf("https://%s/%s/_status", ElasticServerHost, ElasticSearchIndex)
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

func (e ElasticService) getRelatedLog(w rest.ResponseWriter, r *rest.Request) {
	p := r.PathParam("type")
	uri := fmt.Sprintf("https://%s/%s/%s/_search?pretty=true", ElasticServerHost, ElasticSearchIndex, p)

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
