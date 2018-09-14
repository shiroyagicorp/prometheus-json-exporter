package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ReceiverFunc func(key string, value float64)

func (receiver ReceiverFunc) Receive(key string, value float64) {
	receiver(key, value)
}

type Receiver interface {
	Receive(key string, value float64)
}

func WalkJSON(path string, jsonData interface{}, receiver Receiver) {
	switch v := jsonData.(type) {
	case int:
		receiver.Receive(path, float64(v))
	case float64:
		receiver.Receive(path, v)
	case bool:
		n := 0.0
		if v {
			n = 1.0
		}
		receiver.Receive(path, n)
	case string:
		// ignore
	case nil:
		// ignore
	case []interface{}:
		prefix := path + "__"
		for i, x := range v {
			WalkJSON(fmt.Sprintf("%s%d", prefix, i), x, receiver)
		}
	case map[string]interface{}:
		var prefix string
		if path != "" {
			prefix = path + "."
		} else {
			prefix = ""
		}
		for k, x := range v {
			WalkJSON(fmt.Sprintf("%s%s", prefix, k), x, receiver)
		}
	default:
		log.Printf("unkown type: %#v", v)
	}
}

func doProbe(target string) (interface{}, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Get(target)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsonData interface{}
	err = json.Unmarshal([]byte(bytes), &jsonData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func probeHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	target := params.Get("target")
	if target == "" {
		http.Error(w, "Target parameter is missing", http.StatusBadRequest)
		return
	}

	jsonData, err := doProbe(target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// log.Printf("Retrieved value %v", jsonData)

	registry := prometheus.NewRegistry()

	WalkJSON("", jsonData, ReceiverFunc(func(key string, value float64) {
		g := prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: key,
				Help: "Retrieved value",
			},
		)
		registry.MustRegister(g)
		g.Set(value)
	}))

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

var indexHTML = []byte(`<html>
<head><title>Json Exporter</title></head>
<body>
<h1>Json Exporter</h1>
<p><a href="/probe">Run a probe</a></p>
<p><a href="/metrics">Metrics</a></p>
</body>
</html>`)

func main() {
	addr := flag.String("listen-address", ":9116", "The address to listen on for HTTP requests.")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(indexHTML)
	})
	flag.Parse()
	http.HandleFunc("/probe", probeHandler)
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("listenning on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
