package prom

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	CodeRoundTripFailed = 598
)

var (
	transportMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "pili_transport_stat",
		Help: "PILI transport stat",
	}, []string{"mod", "host", "code"})
)

func init() {
	prometheus.Register(transportMetric)
}

type Transport struct {
	tr  http.RoundTripper
	mod string
}

func NewTransport(mod string, tr http.RoundTripper) *Transport {
	if tr == nil {
		tr = http.DefaultTransport
	}
	return &Transport{
		tr:  tr,
		mod: mod,
	}
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.tr.RoundTrip(req)
	var code int
	if err != nil {
		code = CodeRoundTripFailed
	} else {
		code = resp.StatusCode
	}
	host := req.Host
	if host == "" {
		host = req.URL.Host
	}
	transportMetric.With(prometheus.Labels{"mod": t.mod, "host": host, "code": strconv.Itoa(code)}).Inc()
	return
}
