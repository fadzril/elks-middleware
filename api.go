package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/mattbaird/elastigo/lib"
)

var (
	lock = sync.RWMutex{}
)

// Subscription struct
type Subscription struct {
	Uri              string `json:"uri"`
	Selection        string `json:"selection"`
	Type             string `json:"type"`
	App_Id           string `json:"app_id"`
	Component        string `json:"component"`
	Hostname         string `json:"hostname"`
	Transaction_Type string `json:"transaction_type"`
	Subscription_Id  string `json:"subscription_id"`
	Event_Id         string `json:"event_id"`
	Category         string `json:"category"`
	Source_Filename  string `json:"source_filename"`
	Target_Filename  string `json:"target_filename"`
	Messages         string `json:"messages"`
	Tags             string `json:"tags",omitempty`
	Received_At      string `json:"received_at"`
	Status           string
}

// Subscriptions List struct
type Subscriptions struct {
	Total   int
	Message string
	Result  []Subscription
}

// GetVersion Public
func GetVersion() string {
	return "1.0.0.1"
}

func formatDate(s string, t string) string {
	parseTime, error := time.Parse(time.RFC3339, s)
	if error != nil {
		tm := time.Now()
		if t == "from" {
			ts := tm.Add(-24 * time.Hour)
			return fmt.Sprintf("%v", ts.Format(time.RFC3339))
		}
		return fmt.Sprintf("%v", tm.Format(time.RFC3339))
	}
	return fmt.Sprintf("%v", parseTime.Format(time.RFC3339))
}

func formatResponse(out *elastigo.Hits) (res *Subscriptions, err error) {
	if len(out.Hits) >= 1 {
		res := &Subscriptions{
			Total: len(out.Hits),
		}
		for _, item := range out.Hits {
			t, e := json.Marshal(item.Source)
			if e != nil {
				log.Println("Error marshalling", e)
			}

			var sub Subscription
			ers := json.Unmarshal(t, &sub)
			if ers != nil {
				log.Println("Error Unmarshalling", ers)
			}

			sub.Uri = item.Id
			res.Result = append(res.Result, sub)
		}

		res.Message = "OK"
		res.Total = len(res.Result)
		return res, nil
	} else {
		res := &Subscriptions{
			Total:   0,
			Message: "Failed",
			Result:  []Subscription{},
		}
		return res, err
	}
}

/*
 * GetSubscription
 * Params:
 *		Query: size
 *		Query: from
 *		Query: to
 * Example:
 * http://localhost:8000/rest/index?size=500&q=keywords&from=2015-12-12T23:00:00Z&to=
 */
func GetSubscription(w rest.ResponseWriter, r *rest.Request) {
	lock.Lock()

	client := elastigo.NewConn()
	client.Domain = ElasticServerHost
	log.Printf("Dialing ElasticServerHost:ElasticSearchIndex \t%s:%s", ElasticServerHost, ElasticSearchIndex)

	// DSL Query Using Range
	size := r.URL.Query().Get("size")
	query_string := r.URL.Query().Get("q")

	from := formatDate(
		r.URL.Query().Get("from"),
		"from",
	)
	to := formatDate(
		r.URL.Query().Get("to"),
		"to",
	)

	query := `{
	  "size": ` + size + `,
	  "sort": [
		{
		  "@timestamp": {
			"order": "asc",
			"unmapped_type": "boolean"
		  }
		}
	  ],
	  "query": {
		"filtered": {
		  "query": {
			"query_string": {
			  "analyze_wildcard": true,
			  "query": "` + query_string + `"
			}
		  },
		  "filter": {
			"bool": {
			  "must": [
				{
				  "range": {
					"@timestamp": {
					  "gte": "` + from + `",
					  "lte": "` + to + `"
					}
				  }
				}
			  ],
			  "must_not": []
			}
		  }
		}
	  },
	  "fields": [
		"*",
		"_source"
	  ],
	  "fielddata_fields": [
		"@timestamp",
		"arrived_at"
	  ]
	}`

	defer client.Close()
	out, err := client.Search(ElasticSearchIndex, "", nil, query)

	if err != nil {
		log.Println("Error Getting Faceted Search \n", err)
	}

	result, err := formatResponse(&out.Hits)
	if err != nil {
		log.Println(err)
		w.WriteJson([]byte(`{"Result": "Not Found"}`))
	}

	lock.Unlock()
	w.WriteJson(&result)
}

/*
 * GetSubscriptionByCategory
 * Params:
 * 		Query: size
 * 		Query: category
 * Example:
 * localhost:8000/log/category?query=info,error&size=200
 */
func GetSubscriptionByCategory(w rest.ResponseWriter, r *rest.Request) {
	lock.Lock()

	client := elastigo.NewConn()
	client.Domain = ElasticServerHost

	size := r.URL.Query().Get("size")
	categories := r.URL.Query().Get("query")

	query := `{
		"size": ` + size + `,
		"sort": [{
			"@timestamp": {
				"order": "desc",
				"unmapped_type": "boolean"
			}
		}],
		"query": {
			"filtered": {
				"query": {
					"query_string": {
						"analyze_wildcard": true,
						"query": "NOT type:gui AND category:` + categories + `"
					}
				},
				"filter": {
					"bool": {
						"must": [{
							"range": {
								"@timestamp": {
									"gte": "now-90d/d",
									"lte": "now/d"
								}
							}
						}],
						"must_not": []
					}
				}
			}
		},
		"fields": ["*", "_source"],
		"fielddata_fields": ["@timestamp","ts"]
	}`

	out, err := client.Search(ElasticSearchIndex, "", nil, query)
	if err != nil {
		log.Println("Error Getting Result:\n", err)
	}

	result, err := formatResponse(&out.Hits)

	if err != nil {
		w.WriteJson([]byte(`{"Result": "Not Found"}`))
		log.Println("Response format failed:\n", err)
	}

	lock.Unlock()
	w.WriteJson(&result)
}

/*
 * GetSubscriptionByTags
 * This search only target range 1 month: now-30d/d
 * Params:
 *		Query: query
 */
func GetSubscriptionByTags(w rest.ResponseWriter, r *rest.Request) {
	client := elastigo.NewConn()
	client.Domain = ElasticServerHost

	lock.Lock()
	tags := r.URL.Query().Get("query")

	log.Printf("Getting By Tags. Params is : %v", tags)
	query := `{
		"size": 1000,
		"sort": [{
			"@timestamp": {
				"order": "desc",
				"unmapped_type": "boolean"
			}
		}],
		"query": {
			"filtered": {
				"query": {
					"query_string": {
						"analyze_wildcard": true,
						"query": "messages:[\"\" TO *] AND subscription_id:[\"\" TO *] AND tags:` + tags + `"
					}
				},
				"filter": {
					"bool": {
						"must": [{
							"range": {
								"@timestamp": {
									"gte": "now-30d/d",
									"lte": "now/d"
								}
							}
						}],
						"must_not": []
					}
				}
			}
		},
		"fields": ["*", "_source"],
		"fielddata_fields": ["@timestamp","ts"]
	}`

	out, err := client.Search(ElasticSearchIndex, "", nil, query)
	if err != nil {
		log.Println("Error triggering GetSubscriptionByTags \n", err)
	}

	result, err := formatResponse(&out.Hits)
	if err != nil {
		log.Println("Error while formatting Response from GetSubscriptionByTags", err)
	}

	lock.Unlock()
	w.WriteJson(&result)
}
