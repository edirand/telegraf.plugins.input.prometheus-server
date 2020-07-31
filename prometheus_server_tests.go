package prometheusserver

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/influxdata/telegraf/testutil"
)

//TestScrap integration test
func TestScrap(t *testing.T) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	s := &Scrap{
		Host:    "https://127.0.0.1",
		Metrics: "kafka_broker_topic_in_messages_total:rate1m{platform=\"production\", cluster=\"cluster\", topic=\"topic\"};kafka_consumergroup_lag{platform=\"production\",topic=\"topic\"}",
		Client:  &http.Client{Transport: tr},
	}

	var acc testutil.Accumulator
	s.Gather(&acc)
}
