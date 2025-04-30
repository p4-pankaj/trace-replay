package models

import "time"

type TraceRecord struct {
	OutgoingProtocols
	TraceID   string        `json:"trace_id" bson:"trace_id"`
	Timestamp time.Time     `json:"timestamp" bson:"timestamp"`
	Duration  time.Duration `json:"duration" bson:"duration"`
	Logs      []*LogEntry   `json:"log" bson:"log"`
	Dump      any           `json:"dump" bson:"dump"`
}

type ProtocolTrace struct {
	Duration      time.Duration `json:"duration" bson:"duration"`
	ReqTimestamp  time.Time     `json:"reqTimestamp" bson:"respTimestamp"`
	RespTimestamp time.Time     `json:"respTimestamp" bson:"respTimestamp"`
}

type OutgoingProtocols struct {
	HTTPTrace    []*HTTPTraceData    `json:"http_trace,omitempty" bson:"http_trace,omitempty"`
	GRPCTrace    []*GRPCTraceData    `json:"grpc_trace,omitempty" bson:"grpc_trace,omitempty"`
	GraphQLTrace []*GraphQLTraceData `json:"graphql_trace,omitempty" bson:"graphql_trace,omitempty"`
}

type HTTPTraceData struct {
	ProtocolTrace
	Method       string            `json:"method" bson:"method"`
	URL          string            `json:"url" bson:"url"`
	Headers      map[string]string `json:"headers" bson:"headers"`
	Body         string            `json:"body" bson:"body"`
	StatusCode   int               `json:"status_code" bson:"status_code"`
	ResponseBody string            `json:"response_body" bson:"response_body"`
}

type GRPCTraceData struct {
	ProtocolTrace
	Method   string `json:"method" bson:"method"`
	Request  string `json:"request" bson:"request"`
	Response string `json:"response" bson:"response"`
	Status   string `json:"status" bson:"status"`
}

type GraphQLTraceData struct {
	ProtocolTrace
	Query     string         `json:"query" bson:"query"`
	Variables map[string]any `json:"variables" bson:"variables"`
	Response  string         `json:"response" bson:"response"`
}
