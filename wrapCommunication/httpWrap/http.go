package httpWrap

import (
	"context"
	"errors"
	"time"

	"github.com/p4-pankaj/trace-replay/models"
	"github.com/p4-pankaj/trace-replay/tracer"
	"github.com/p4-pankaj/trace-replay/utility"
)

type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
}

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

// TODO: Add more http features here
// for now it's a proof of concept
// will make it more sophisticated in upcoming version
// user will be able to decide how they want to make http request
// and will be able to pass an implementation of following interface
// bydefault we will use net/http implemented out of the box.
type Client interface {
	Do(ctx context.Context, req *Request) (
		*Response, error)
}

func Do__Trace(r *tracer.TraceRecorder,
	c Client, ctx context.Context, req *Request) (
	resp *Response, err error) {
	callId := utility.GetHashForObject(req)
	if r.Env == "DEBUG" {
		trace, exist := r.TraceRecord.HTTPTrace[callId]
		if !exist {
			panic("this is unexpected, server didn't recieve this request on prod env.please review TRACE_ID or manually or try again")
		}
		if !trace.NilResp {
			resp = &Response{}
			resp.Body = trace.ResponseBody
			resp.Headers = trace.RespHeaders
			resp.StatusCode = trace.StatusCode
		}
		if trace.ErrString != nil {
			err = errors.New(*trace.ErrString)
		}

		time.Sleep(trace.Duration)
	} else {
		if r.TraceRecord.HTTPTrace == nil {
			r.TraceRecord.HTTPTrace = map[string]*models.HttpTraceData{}
		}
		httpTrace := models.HttpTraceData{}
		httpTrace.ReqTimestamp = time.Now()
		httpTrace.CallID = callId
		httpTrace.ReqHeaders = req.Headers
		httpTrace.ReqBody = req.Body
		httpTrace.URL = req.URL
		httpTrace.Method = req.Method

		resp, err = c.Do(ctx, req)

		if resp != nil {
			httpTrace.ResponseBody = resp.Body
			httpTrace.StatusCode = resp.StatusCode
			httpTrace.RespHeaders = resp.Headers
		} else {
			httpTrace.NilResp = true
		}

		if err != nil {
			s := err.Error()
			httpTrace.ErrString = &s
		}

		httpTrace.RespTimestamp = time.Now()
		httpTrace.Duration = httpTrace.RespTimestamp.
			Sub(httpTrace.ReqTimestamp)

		r.TraceRecord.HTTPTrace[callId] = &httpTrace
		r.Updater <- func(t *models.TraceRecord) {
			t.FunctionTrace = r.TraceRecord.FunctionTrace
		}

	}

	return

}
