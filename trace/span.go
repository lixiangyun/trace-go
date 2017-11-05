package trace

import (
	"time"
)

const (
	CLIENT = 0
	SERVER = 1
)

var globalCollector *Collector

type SPAN_TYPE int

type Context struct {
	TraceID   string `json:"traceid"`
	TraceName string `json:"name"`
	SpanID    string `json:"id"`
	ParentID  string `json:"parentid"`
}

type Endpoint struct {
	SrvName string `json:"serviceName"`
	IP      string `json:"ip"`
	Port    int16  `json:"port"`
}

//value : cs,cr,ss,sr
type Stage struct {
	Timestamp int64    `json:"timestamp"`
	Value     string   `json:"value"`
	Host      Endpoint `json:"endpoint"`
}

type KeyValue struct {
	Key   string   `json:"key"`
	Value string   `json:"value"`
	Host  Endpoint `json:"endpoint"`
}

type Span struct {
	sptype SPAN_TYPE
	ctx    Context
	step   []*Stage
	kv     []*KeyValue
}

type SpanRecord struct {
	TraceID   string      `json:"traceId"`
	SpanID    string      `json:"id"`
	TraceName string      `json:"name"`
	ParentID  string      `json:"parentId"`
	Timestamp int64       `json:"timestamp"`
	Duration  int64       `json:"duration"`
	StageList []*Stage    `json:"annotations"`
	Kvlist    []*KeyValue `json:"binary_annotations"`
}

func NewEndPoint(srvname, ip string, port int16) *Endpoint {
	return &Endpoint{IP: ip, Port: port, SrvName: srvname}
}

func RecvSpan(p Context) *Span {

	s := &Span{sptype: SERVER}
	s.ctx.TraceID = p.TraceID
	s.ctx.ParentID = p.ParentID
	s.ctx.TraceName = p.TraceName
	s.ctx.SpanID = p.SpanID

	s.step = make([]*Stage, 0)
	s.kv = make([]*KeyValue, 0)

	return s
}

func NewSpan(p Context) *Span {

	s := &Span{sptype: CLIENT}
	s.ctx.TraceID = p.TraceID
	s.ctx.ParentID = p.SpanID
	s.ctx.TraceName = p.TraceName
	s.ctx.SpanID = getSpanID()

	s.step = make([]*Stage, 0)
	s.kv = make([]*KeyValue, 0)

	return s
}

func (s *Span) GetContext() Context {
	return s.ctx
}

func (s *Span) Begin(host *Endpoint) {
	stage := new(Stage)
	stage.Host = *host
	stage.Timestamp = int64(time.Now().Nanosecond())
	if s.sptype == CLIENT {
		stage.Value = "cs"
	} else {
		stage.Value = "ss"
	}
	s.step = append(s.step, stage)
}

func (s *Span) AddKV(key, value string, host *Endpoint) {
	kv := &KeyValue{Key: key, Value: value, Host: *host}
	s.kv = append(s.kv, kv)
}

func (s *Span) End(host *Endpoint) {
	stage := new(Stage)
	stage.Host = *host
	stage.Timestamp = int64(time.Now().Nanosecond())
	if s.sptype == CLIENT {
		stage.Value = "cr"
	} else {
		stage.Value = "sr"
	}
	s.step = append(s.step, stage)

	span := new(SpanRecord)
	span.TraceName = s.ctx.TraceName
	span.TraceID = s.ctx.TraceID
	span.SpanID = s.ctx.SpanID
	span.ParentID = s.ctx.ParentID
	span.StageList = s.step
	span.Kvlist = s.kv
	span.Timestamp = s.step[0].Timestamp
	span.Duration = s.step[1].Timestamp - s.step[0].Timestamp

	globalCollector.Record(span)
}