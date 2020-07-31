package prometheusserver

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

//Scrap is prometheus scrapper configuration
type Scrap struct {
	Metrics string
	Host    string
	Client  *http.Client
}

//PrometheusResponse is the prometheus API response
type PrometheusResponse struct {
	Status    string
	ErrorType string
	Error     string
	Data      PrometheusData
}

//PrometheusData is the struct containing data in the prometheus response
type PrometheusData struct {
	ResultType string
	Result     []map[string]interface{}
}

const scrapConfig = `
## The list of metrics to collect (you can specify multiple metrics using semicolons)
metrics = """
some_metric{someTag=\"abc\"};
some_other_metrics:rate1m{someTag1=\"abcd\", someTag2=\"efgh\"};
"""

##The prometheus url
host = "http://prometheus"
`

//SampleConfig generates a sample configuration
func (s *Scrap) SampleConfig() string {
	return scrapConfig
}

//Description gets the description of the plugin
func (s *Scrap) Description() string {
	return "Collects metrics from prometheus server to export or work on it."
}

//Gather data
func (s *Scrap) Gather(acc telegraf.Accumulator) error {

	var metrics = strings.Split(s.Metrics, ";")
	for index := range metrics {
		var prom = callAPI(s.Host, metrics[index], s.Client)
		if prom.Status == "success" {
			for _, result := range prom.Data.Result {
				var tags = parseTags(result["metric"].(map[string]interface{}))
				var metricName = tags["__name__"]
				delete(tags, "__name__")
				for key, field := range result {
					if key != "metric" {
						fields := make(map[string]interface{})
						value, err := strconv.ParseFloat(field.([]interface{})[1].(string), 64)
						if err != nil {
							fields[key] = 0
						} else {
							fields[key] = value
						}
						acc.AddFields(metricName, fields, tags)
					}
				}
			}
		}
	}
	return nil
}

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	inputs.Add("prometheus_scrapper", func() telegraf.Input { return &Scrap{Client: &http.Client{Transport: tr}} })
}

func callAPI(host string, metrics string, client *http.Client) PrometheusResponse {
	var prom PrometheusResponse
	req, _ := http.NewRequest("GET", host+"/api/v1/query", nil)
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("User-Agent", "PostmanRuntime/7.26.2")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	query := req.URL.Query()
	query.Add("query", metrics)
	req.URL.RawQuery = query.Encode()

	response, err := client.Do(req)

	if err == nil {
		defer response.Body.Close()
		json.NewDecoder(response.Body).Decode(&prom)
	} else {
		log.Panic("An error occured during call to prometheus server : ", err)
	}
	return prom
}

func parseTags(aMap map[string]interface{}) map[string]string {
	tags := make(map[string]string)
	for key, val := range aMap {
		tags[key] = val.(string)
	}
	return tags
}
