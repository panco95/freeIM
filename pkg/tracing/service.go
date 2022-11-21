package tracing

import (
	"io"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

type jaegerLoggerAdapter struct {
	logger *zap.SugaredLogger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Infof(strings.Trim(msg, "\n\r "), args...)
}

func NewLogger(log *zap.SugaredLogger) jaeger.Logger {
	return &jaegerLoggerAdapter{log}
}

type TracingService struct {
	log    *zap.SugaredLogger
	closer io.Closer

	enableLogger bool
}

func NewTracingService(enableLogger bool) *TracingService {
	return &TracingService{
		enableLogger: enableLogger,
		log:          zap.S().With("module", "tracing"),
	}
}

func (s *TracingService) InitGlobal(cfg *viper.Viper) error {
	if cfg == nil {
		return nil
	}

	if s.closer != nil {
		s.closer.Close()
		s.closer = nil
	}

	jcfg := jaegercfg.Configuration{}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "yaml",
		Squash:           true,
		Metadata:         nil,
		Result:           &jcfg,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	})
	if err != nil {
		return err
	}
	err = decoder.Decode(cfg.AllSettings())
	if err != nil {
		return err
	}

	options := []jaegercfg.Option{}
	if s.enableLogger {
		options = append(options, jaegercfg.Logger(NewLogger(zap.S().With("module", "tracing.jaeger"))))
	}
	tracer, closer, err := jcfg.NewTracer(options...)
	if err != nil {
		return err
	}
	s.closer = closer

	opentracing.SetGlobalTracer(tracer)

	return nil
}
