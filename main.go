package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

const templateMetrics string = `Active connections: %d
server accepts handled requests
%d %d %d
Reading: %d Writing: %d Waiting: %d
`

type NginxHost struct {
	Host string
}
type StubConnections struct {
	Uri      string
	Active   int64
	Accepted int64
	Handled  int64
	Reading  int64
	Writing  int64
	Waiting  int64
	Requests int64
	Up       int8
}

func (c *NginxHost) GetStubStates() (s StubConnections) {
	apiEndpoint := "http://" + c.Host + "/stub_status"
	s.Uri = apiEndpoint
	s.Up = 0
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*1)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 1))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 1,
		},
	}
	resp, err := client.Get(apiEndpoint)
	if err != nil {
		log.Println(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Println(err)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	resp.Body.Close()
	r := bytes.NewReader(body)
	if _, err := fmt.Fscanf(r, templateMetrics,
		&s.Active,
		&s.Accepted,
		&s.Handled,
		&s.Requests,
		&s.Reading,
		&s.Writing,
		&s.Waiting); err != nil {
		fmt.Println(err)
		return
	}
	s.Up = 1
	return
}

type NginxHostCollector struct {
	NginxHost *NginxHost
}

var (
	nginxActiveConnectionsDesc = prometheus.NewDesc(
		"nginx_connections_active_total",
		"Active client connections.",
		[]string{"uri"},
		nil,
	)
	nginxAcceptedConnectionsDesc = prometheus.NewDesc(
		"nginx_connections_accepted_total",
		"Accepted client connections.",
		[]string{"uri"},
		nil,
	)
	nginxHandledConnectionsDesc = prometheus.NewDesc(
		"nginx_connections_handled_total",
		"Handled client connections.",
		[]string{"uri"},
		nil,
	)
	nginxReadingConnectionsDesc = prometheus.NewDesc(
		"nginx_connections_reading_total",
		"Connections where NGINX is reading the request header.",
		[]string{"uri"},
		nil,
	)
	nginxWritingConnectionsDesc = prometheus.NewDesc(
		"nginx_connections_writing_total",
		"Connections where NGINX is writing the response back to the client.",
		[]string{"uri"},
		nil,
	)
	nginxWaitingConnectionsDesc = prometheus.NewDesc(
		"nginx_connections_waiting_total",
		"Idle client connections.",
		[]string{"uri"},
		nil,
	)
	nginxHttpRequestsDesc = prometheus.NewDesc(
		"nginx_http_requests_total",
		"Total http requests.",
		[]string{"uri"},
		nil,
	)
	nginxUpDesc = prometheus.NewDesc(
		"nginx_up",
		"Total http requests.",
		[]string{"uri"},
		nil,
	)
	hosts []string
)

func (cc NginxHostCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cc, ch)
}

func (cc NginxHostCollector) Collect(ch chan<- prometheus.Metric) {
	nginxStubConnection := cc.NginxHost.GetStubStates()
	if nginxStubConnection.Up == 0 {
		ch <- prometheus.MustNewConstMetric(
			nginxUpDesc,
			prometheus.GaugeValue,
			float64(nginxStubConnection.Up),
			nginxStubConnection.Uri,
		)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		nginxActiveConnectionsDesc,
		prometheus.GaugeValue,
		float64(nginxStubConnection.Active),
		nginxStubConnection.Uri,
	)
	ch <- prometheus.MustNewConstMetric(
		nginxAcceptedConnectionsDesc,
		prometheus.GaugeValue,
		float64(nginxStubConnection.Accepted),
		nginxStubConnection.Uri,
	)
	ch <- prometheus.MustNewConstMetric(
		nginxHandledConnectionsDesc,
		prometheus.GaugeValue,
		float64(nginxStubConnection.Handled),
		nginxStubConnection.Uri,
	)
	ch <- prometheus.MustNewConstMetric(
		nginxReadingConnectionsDesc,
		prometheus.GaugeValue,
		float64(nginxStubConnection.Reading),
		nginxStubConnection.Uri,
	)
	ch <- prometheus.MustNewConstMetric(
		nginxWritingConnectionsDesc,
		prometheus.GaugeValue,
		float64(nginxStubConnection.Writing),
		nginxStubConnection.Uri,
	)
	ch <- prometheus.MustNewConstMetric(
		nginxWaitingConnectionsDesc,
		prometheus.GaugeValue,
		float64(nginxStubConnection.Waiting),
		nginxStubConnection.Uri,
	)
	ch <- prometheus.MustNewConstMetric(
		nginxHttpRequestsDesc,
		prometheus.GaugeValue,
		float64(nginxStubConnection.Requests),
		nginxStubConnection.Uri,
	)
	ch <- prometheus.MustNewConstMetric(
		nginxUpDesc,
		prometheus.GaugeValue,
		float64(nginxStubConnection.Up),
		nginxStubConnection.Uri,
	)
}

func NewNginxHost(Host string, reg prometheus.Registerer) *NginxHost {

	c := &NginxHost{
		Host: Host,
	}
	cc := NginxHostCollector{NginxHost: c}
	prometheus.WrapRegistererWith(prometheus.Labels{"host": Host}, reg).MustRegister(cc)
	return c
}

func main() {
	viper.SetConfigFile("config.yml")
	err :=viper.ReadInConfig()
	if err!=nil{
		log.Fatal(err)
	}
	hosts = viper.GetStringSlice("host")	
	reg := prometheus.NewPedanticRegistry()
	for _, host := range hosts {
		log.Println("Found nginx host:",host)
		NewNginxHost(host, reg)
	}
	if len(hosts) < 1 {
		log.Fatal("No nginx host found,check config.yml")
	}
	reg.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":9113", nil))
}
