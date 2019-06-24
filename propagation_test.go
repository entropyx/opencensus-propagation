package propagation

import (
	"context"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.opencensus.io/trace"
)

func TestPropagation(t *testing.T) {
	Convey("Given a span", t, func() {

		_, span := trace.StartSpan(context.Background(), "testing")

		Convey("When the span is injected into a http header", func() {
			request, _ := http.NewRequest("POST", "example.gom", nil)
			header := request.Header
			Inject(span.SpanContext(), FormatTextMap, HTTPHeader(header))

			Convey("The header should include the span", func() {
				So(header.Get(HeaderSpanID), ShouldNotBeBlank)
				So(header.Get(HeaderTraceID), ShouldNotBeBlank)
				So(header.Get(HeaderTraceOptions), ShouldNotBeBlank)
			})

			Convey("When the span context is extracted from the header", func() {
				extractedCtx, err := Extract(FormatTextMap, HTTPHeader(header))

				Convey("err should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The extracted should be equal to the original span context", func() {
					So(extractedCtx.SpanID, ShouldEqual, span.SpanContext().SpanID)
					So(extractedCtx.TraceID, ShouldEqual, span.SpanContext().TraceID)
					So(extractedCtx.TraceOptions, ShouldEqual, span.SpanContext().TraceOptions)
				})
			})
		})
	})
}
