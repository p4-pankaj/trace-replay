package models

import "time"

type TraceRecord struct {
	OutgoingProtocols
	TraceID       string          `json:"trace_id" bson:"trace_id"`
	Timestamp     time.Time       `json:"timestamp" bson:"timestamp"`
	Duration      time.Duration   `json:"duration" bson:"duration"`
	Logs          []*LogEntry     `json:"log" bson:"log"`
	FunctionTrace []FunctionTrace `json:"functionTrace" bson:"functionTrace"`
	Dump          any             `json:"dump" bson:"dump"`
}

type FunctionTrace struct {
	TraceInfo
	Input  any
	Output any
}

type TraceInfo struct {
	Duration      time.Duration `json:"duration" bson:"duration"`
	ReqTimestamp  time.Time     `json:"reqTimestamp" bson:"reqTimestamp"`
	RespTimestamp time.Time     `json:"respTimestamp" bson:"respTimestamp"`
	CallID        string        `json:"call_id" bson:"call_id"`
	CallIndex     int           `json:"call_index" bson:"call_index"`
}

type OutgoingProtocols struct {
	HTTPTrace    map[string]*HttpTraceData `json:"http_trace,omitempty" bson:"http_trace,omitempty"`
	GRPCTrace    []*GRPCTraceData          `json:"grpc_trace,omitempty" bson:"grpc_trace,omitempty"`
	GraphQLTrace []*GraphQLTraceData       `json:"graphql_trace,omitempty" bson:"graphql_trace,omitempty"`
}

type HttpTraceData struct {
	TraceInfo
	HttpRequestIdentifier
	ErrString    *string           `json:"errString,omitempty" bson:"errString,omitempty"`
	StatusCode   int               `json:"status_code" bson:"status_code"`
	ResponseBody []byte            `json:"response_body" bson:"response_body"`
	RespHeaders  map[string]string `json:"respHeaders" bson:"respHeaders"`
	NilResp      bool              `json:"nilResp" bson:"nilResp"`
}

type HttpRequestIdentifier struct {
	Method     string            `json:"method" bson:"method"`
	URL        string            `json:"url" bson:"url"`
	ReqHeaders map[string]string `json:"headers" bson:"headers"`
	ReqBody    []byte            `json:"body" bson:"body"`
}

type GRPCTraceData struct {
	TraceInfo
	Method   string `json:"method" bson:"method"`
	Request  string `json:"request" bson:"request"`
	Response string `json:"response" bson:"response"`
	Status   string `json:"status" bson:"status"`
}

type GraphQLTraceData struct {
	TraceInfo
	Query     string         `json:"query" bson:"query"`
	Variables map[string]any `json:"variables" bson:"variables"`
	Response  string         `json:"response" bson:"response"`
}
