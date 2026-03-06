package ayopentracing

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

	"github.com/AyakuraYuki/go-aybox/log"
	ayStack "github.com/AyakuraYuki/go-aybox/stack"
)

func jaegerHost() string {
	if s := os.Getenv("JAEGER_HOST"); s != "" {
		return s
	}
	return "127.0.0.1:6831"
}

// NewJaegerTracer udp trace
func NewJaegerTracer(serviceName, env string) (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		ServiceName: fmt.Sprintf("%s-%s", env, serviceName),
		Tags:        []opentracing.Tag{{Key: "env", Value: env}},
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  jaegerHost(),
		},
	}

	tracer, closer, err := cfg.NewTracer()

	if err != nil {
		log.Error("NewTracer").Msgf("%v", err)
	}
	if tracer != nil {
		opentracing.SetGlobalTracer(tracer)
	}

	return tracer, closer, err
}

// // NewJaegerTracerHttp http trace
// func NewJaegerTracerHttp(serviceName, env string) (opentracing.Tracer, io.Closer, error) {
//	cfg := jaegercfg.Configuration{
//		ServiceName: fmt.Sprintf("%s-%s", env, serviceName),
//		Tags:        []opentracing.Tag{{Key: "env", Value: env}},
//		Sampler: &jaegercfg.SamplerConfig{
//			Type:  jaeger.SamplerTypeConst,
//			Param: 1,
//		},
//		Reporter: &jaegercfg.ReporterConfig{
//			LogSpans:            true,
//			BufferFlushInterval: 1 * time.Second,
//		},
//	}
//	var (
//		_httpUrl = _httpUrlInner
//	)
//
//	sender := transport.NewHTTPTransport(_httpUrl)
//	reporter := jaeger.NewRemoteReporter(sender)
//	tracer, closer, err := cfg.NewTracer(
//		jaegercfg.Reporter(reporter),
//	)
//	if err != nil {
//		log.Error("NewTracer").Msgf("%v", err)
//	}
//	if tracer != nil {
//		opentracing.SetGlobalTracer(tracer)
//	}
//
//	return tracer, closer, err
// }

// StartSpanFromContext 新建span
func StartSpanFromContext(ctx context.Context, opName ...string) (opentracing.Span, context.Context) {
	operationName := "unknow"
	if len(opName) > 0 && len(opName[0]) > 0 {
		operationName = opName[0]
	} else {
		pc, _, _, _ := runtime.Caller(1)
		if f := runtime.FuncForPC(pc); f != nil {
			if arr := strings.Split(f.Name(), "/"); len(arr) > 0 {
				operationName = arr[len(arr)-1]
			}
		}
	}
	return opentracing.StartSpanFromContext(ctx, operationName)
}

// StartSpanFromContextWithSt 新建span
func StartSpanFromContextWithSt(ctx context.Context, opName string, startTime time.Time) (opentracing.Span, context.Context) {
	ops := opentracing.StartTime(startTime)

	operationName := "unknow"
	if len(opName) > 0 {
		operationName = opName
	} else {
		pc, _, _, _ := runtime.Caller(1)
		if f := runtime.FuncForPC(pc); f != nil {
			if arr := strings.Split(f.Name(), "/"); len(arr) > 0 {
				operationName = arr[len(arr)-1]
			}
		}
	}

	return opentracing.StartSpanFromContext(ctx, operationName, ops)
}

func SpanError(span opentracing.Span, err error) {
	if err == nil {
		return
	}

	ext.Error.Set(span, true)
	span.LogKV("err", err)

	var stacktrace *ayStack.Stacktrace
	if errors.As(err, &stacktrace) {
		span.LogKV("__stacktrace__", stacktrace.Sources())
	}
}

func SpanLogKV(span opentracing.Span, keyValues ...any) {
	if len(keyValues) == 0 || len(keyValues)%2 != 0 {
		return
	}
	for i := 0; i*2 < len(keyValues); i++ {
		key, ok := keyValues[i*2].(string)
		if !ok {
			continue
		}
		switch typedVal := keyValues[i*2+1].(type) {
		case bool, string, int, int8, int16, int32, int64, uint, uint64, uint8, uint16, uint32, float32, float64:
			span.LogKV(key, typedVal)
		case error:
			if typedVal == nil || (reflect.ValueOf(typedVal).Kind() == reflect.Ptr && reflect.ValueOf(typedVal).IsNil()) {
				span.LogKV(key, "nil")
			} else {
				ext.Error.Set(span, true)
				span.LogKV(key, typedVal)

				var stacktrace *ayStack.Stacktrace
				if errors.As(typedVal, &stacktrace) {
					span.LogKV("__stacktrace__", stacktrace.Sources())
				}
			}
		default:
			if typedVal == nil || (reflect.ValueOf(typedVal).Kind() == reflect.Ptr && reflect.ValueOf(typedVal).IsNil()) {
				span.LogKV(key, "nil")
			} else {
				raw, _ := sonic.MarshalString(keyValues[i*2+1])
				span.LogKV(key, raw)
			}
		}
	}
}
