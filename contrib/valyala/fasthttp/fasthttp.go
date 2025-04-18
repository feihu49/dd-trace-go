// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

// Package fasthttp provides functions to trace the valyala/fasthttp package (https://github.com/valyala/fasthttp)
package fasthttp // import "github.com/DataDog/dd-trace-go/contrib/valyala/fasthttp/v2"

import (
	"fmt"
	"strconv"

	"github.com/valyala/fasthttp"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/ext"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/DataDog/dd-trace-go/v2/instrumentation"
)

var instr *instrumentation.Instrumentation

func init() {
	instr = instrumentation.Load(instrumentation.PackageValyalaFastHTTP)
}

// WrapHandler wraps a fasthttp.RequestHandler with tracing middleware
func WrapHandler(h fasthttp.RequestHandler, opts ...Option) fasthttp.RequestHandler {
	cfg := newConfig()
	for _, fn := range opts {
		fn.apply(cfg)
	}
	instr.Logger().Debug("contrib/valyala/fasthttp.v1: Configuring Middleware: cfg: %#v", cfg)
	return func(fctx *fasthttp.RequestCtx) {
		if cfg.ignoreRequest(fctx) {
			h(fctx)
			return
		}
		spanOpts := []tracer.StartSpanOption{
			tracer.ServiceName(cfg.serviceName),
		}
		spanOpts = append(spanOpts, defaultSpanOptions(fctx)...)
		fcc := &HTTPHeadersCarrier{
			ReqHeader: &fctx.Request.Header,
		}
		if sctx, err := tracer.Extract(fcc); err == nil {
			// If there are span links as a result of context extraction, add them as a StartSpanOption
			if sctx != nil && sctx.SpanLinks() != nil {
				spanOpts = append(spanOpts, tracer.WithSpanLinks(sctx.SpanLinks()))
			}
			spanOpts = append(spanOpts, tracer.ChildOf(sctx))
		}
		span := StartSpanFromContext(fctx, "http.request", spanOpts...)
		defer span.Finish()
		h(fctx)
		span.SetTag(ext.ResourceName, cfg.resourceNamer(fctx))
		status := fctx.Response.StatusCode()
		if cfg.isStatusError(status) {
			span.SetTag(ext.Error, fmt.Errorf("%d: %s", status, string(fctx.Response.Body())))
		}
		span.SetTag(ext.HTTPCode, strconv.Itoa(status))
	}
}

func defaultSpanOptions(fctx *fasthttp.RequestCtx) []tracer.StartSpanOption {
	opts := []tracer.StartSpanOption{
		tracer.Tag(ext.Component, instrumentation.PackageValyalaFastHTTP),
		tracer.Tag(ext.SpanKind, ext.SpanKindServer),
		tracer.SpanType(ext.SpanTypeWeb),
		tracer.Tag(ext.HTTPMethod, string(fctx.Method())),
		tracer.Tag(ext.HTTPURL, string(fctx.URI().FullURI())),
		tracer.Tag(ext.HTTPUserAgent, string(fctx.UserAgent())),
		tracer.Measured(),
	}
	if host := string(fctx.Host()); len(host) > 0 {
		opts = append(opts, tracer.Tag("http.host", host))
	}
	return opts
}
