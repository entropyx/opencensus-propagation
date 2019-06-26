package propagation

import (
	"encoding/hex"
	"strconv"

	"go.opencensus.io/trace"
)

type Format uint8

const (
	FormatBinary Format = iota
	FormatTextMap
)

const (
	HeaderTraceID      = "Trace-Id"
	HeaderTraceOptions = "Trace-Options"
	HeaderSpanID       = "Span-Id"
)

type TextMap interface {
	Set(key string, value string)
	ForeachKey(func(key string, value string) error) error
}

func Extract(format Format, carrier interface{}) (trace.SpanContext, error) {
	spanCtx := trace.SpanContext{}
	switch format {
	case FormatTextMap:
		m := carrier.(TextMap)
		err := m.ForeachKey(func(key string, value string) error {
			switch key {
			case HeaderSpanID:
				if err := decodeAndCopyString(value, spanCtx.SpanID[:]); err != nil {
					return err
				}
			case HeaderTraceID:
				if err := decodeAndCopyString(value, spanCtx.TraceID[:]); err != nil {
					return err
				}
			case HeaderTraceOptions:
				i, err := strconv.Atoi(value)
				if err != nil {
					return err
				}
				spanCtx.TraceOptions = trace.TraceOptions(i)
			}
			return nil
		})
		if err != nil {
			return spanCtx, err
		}
	}
	return spanCtx, nil
}

func Inject(spanCtx trace.SpanContext, format Format, carrier interface{}) {
	switch format {
	case FormatTextMap:
		m := carrier.(TextMap)
		spanID := hex.EncodeToString(spanCtx.SpanID[:])
		m.Set(HeaderSpanID, spanID)
		m.Set(HeaderTraceID, hex.EncodeToString(spanCtx.TraceID[:]))
		m.Set(HeaderTraceOptions, strconv.Itoa(int(spanCtx.TraceOptions)))
	}
}

func decodeAndCopyString(s string, dst []byte) error {
	buf, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	copy(dst, buf)
	return nil
}
