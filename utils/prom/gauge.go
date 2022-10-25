package prom

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Gauge prometheus.Gauge
type Labels prometheus.Labels

// GaugeVec is a Collector that bundles a set of Gauges that all share the same
// Desc, but have different values for their variable labels. This is used if
// you want to count the same thing partitioned by various dimensions
// (e.g. number of operations queued, partitioned by user and operation
// type). Create instances with NewGaugeVec.
type GaugeVec struct {
	internal *prometheus.GaugeVec
	labels   []string
	expireMS int64

	tracker  map[uint64]*trackerValue
	trackerL sync.RWMutex
	closeCh  chan int
}

type GaugeOpts prometheus.GaugeOpts

// NewGaugeVec creates a new GaugeVec based on the provided GaugeOpts and
// partitioned by the given label names. At least one label name must be
// provided.
func NewGaugeVec(opts GaugeOpts, labelNames []string, expireMS int64) *GaugeVec {
	g := &GaugeVec{
		internal: prometheus.NewGaugeVec(prometheus.GaugeOpts(opts), labelNames),
		labels:   labelNames,
		expireMS: expireMS,
		tracker:  make(map[uint64]*trackerValue),
		closeCh:  make(chan int, 1),
	}
	go g.clearLoop()
	return g
}

// With works as GetMetricWith, but panics where GetMetricWithLabels would have
// returned an error. By not returning an error, With allows shortcuts like
//     myVec.With(Labels{"code": "404", "method": "GET"}).Add(42)
func (m *GaugeVec) With(labels Labels) Gauge {
	m.touch(labels)
	return m.internal.With(prometheus.Labels(labels))
}

func (m *GaugeVec) WithLabelValues(lvs ...string) Gauge {
	labels := make(Labels)
	for i, l := range m.labels {
		labels[l] = lvs[i]
	}
	m.touch(labels)
	return m.internal.WithLabelValues(lvs...)
}

func (m *GaugeVec) Close() {
	select {
	case m.closeCh <- 1:
	default:
	}
}

func (m *GaugeVec) Describe(ch chan<- *prometheus.Desc) {
	m.internal.Describe(ch)
}

func (m *GaugeVec) Collect(ch chan<- prometheus.Metric) {
	m.internal.Collect(ch)
}

func (m *GaugeVec) touch(labels Labels) {
	h := hash(m.labels, labels)
	now := time.Now().UnixNano() / int64(time.Millisecond)

	m.trackerL.Lock()
	defer m.trackerL.Unlock()

	val, ok := m.tracker[h]
	if !ok {
		val = &trackerValue{
			Labels: labels,
		}
		m.tracker[h] = val
	}
	val.TouchedAt = now
}

func (m *GaugeVec) clearLoop() {
	ticker := time.NewTicker(time.Duration(m.expireMS) * time.Millisecond / 3)
	for t := range ticker.C {
		select {
		case <-m.closeCh:
			ticker.Stop()
			return
		default:
		}
		var expired []Labels
		m.trackerL.Lock()
		for key, val := range m.tracker {
			if val.TouchedAt+m.expireMS < t.UnixNano()/int64(time.Millisecond) {
				delete(m.tracker, key)
				expired = append(expired, val.Labels)
			}
		}
		m.trackerL.Unlock()

		for _, labels := range expired {
			m.internal.Delete(prometheus.Labels(labels))
		}
	}
}
