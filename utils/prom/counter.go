package prom

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Counter prometheus.Counter

// CounterVec is a Collector that bundles a set of Counters that all share the
// same Desc, but have different values for their variable labels. This is used
// if you want to count the same thing partitioned by various dimensions
// (e.g. number of HTTP requests, partitioned by response code and
// method). Create instances with NewCounterVec.
//
// CounterVec embeds MetricVec. See there for a full list of methods with
// detailed documentation.
type CounterVec struct {
	internal *prometheus.CounterVec
	labels   []string
	expireMS int64

	tracker  map[uint64]*trackerValue
	trackerL sync.RWMutex
	closeCh  chan int
}

type CounterOpts prometheus.CounterOpts

// NewCounterVec creates a new CounterVec based on the provided CounterOpts and
// partitioned by the given label names. At least one label name must be
// provided.
func NewCounterVec(opts CounterOpts, labelNames []string, expireMS int64) *CounterVec {
	g := &CounterVec{
		internal: prometheus.NewCounterVec(prometheus.CounterOpts(opts), labelNames),
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
func (m *CounterVec) With(labels Labels) Counter {
	m.touch(labels)
	return m.internal.With(prometheus.Labels(labels))
}

func (m *CounterVec) WithLabelValues(lvs ...string) Counter {
	labels := make(Labels)
	for i, l := range m.labels {
		labels[l] = lvs[i]
	}
	m.touch(labels)
	return m.internal.WithLabelValues(lvs...)
}

func (m *CounterVec) Close() {
	select {
	case m.closeCh <- 1:
	default:
	}
}

func (m *CounterVec) Describe(ch chan<- *prometheus.Desc) {
	m.internal.Describe(ch)
}

func (m *CounterVec) Collect(ch chan<- prometheus.Metric) {
	m.internal.Collect(ch)
}

func (m *CounterVec) touch(labels Labels) {
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

func (m *CounterVec) clearLoop() {
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
