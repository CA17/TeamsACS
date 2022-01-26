package models

import "sync/atomic"

type Metrics struct {
	Icon  string
	Value interface{}
	Title string
}

func NewMetrics(icon string, value interface{}, title string) *Metrics {
	return &Metrics{Icon: icon, Value: value, Title: title}
}

type NameValuePair struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

func (d *NameValuePair) Incr() {
	atomic.AddInt64(&d.Value, 1)
}

