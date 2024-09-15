package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"github.com/xhayamix/proto-gen-golang/pkg/domain/constant"
	"github.com/xhayamix/proto-gen-golang/pkg/logging/option/gcp"
)

const logName = "app"

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	DPanic(ctx context.Context, msg string, fields ...zap.Field)
	Panic(ctx context.Context, msg string, fields ...zap.Field)
	Fatal(ctx context.Context, msg string, fields ...zap.Field)
	ErrorReport(ctx context.Context, err error, httpRequest *gcp.ErrorHTTPRequest, user string)
	EchoErrorHTTPRequest(c echo.Context) *gcp.ErrorHTTPRequest
	GRPCErrorHTTPRequest(fullMethod string, md metadata.MD, code codes.Code) *gcp.ErrorHTTPRequest
	Sync() error
}

type logger struct {
	*zap.Logger
	projectID string
	isLocal   bool
}

var appLogger Logger

func New(projectID, service, version string, isLocal bool) (Logger, error) {
	var config zap.Config
	if isLocal {
		config = zap.NewDevelopmentConfig()
	} else {
		config = gcp.NewConfig()
	}

	sctx := &gcp.ServiceContext{
		Service: service,
		Version: version,
	}
	l, err := config.Build(
		zap.WrapCore(gcp.NewCore),
		zap.Fields(zap.Object(sctx.Key(), sctx)),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		return nil, err
	}
	return &logger{Logger: l.Named(logName), projectID: projectID, isLocal: isLocal}, nil
}

func GetLogger() Logger {
	if testing.Testing() {
		return &logger{Logger: zap.NewNop(), projectID: "", isLocal: true}
	}

	if appLogger == nil {
		err := errors.New("not initial set AppLogger")
		l, _ := New("", "", "", false)
		l.ErrorReport(context.Background(), err, nil, "")
		return l
	}
	return appLogger
}

func SetLogger(projectID, service, version string, isLocal bool) error {
	l, err := New(projectID, service, version, isLocal)
	if err != nil {
		return err
	}
	appLogger = l
	return nil
}

func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, l.appendFields(ctx, fields)...)
}
func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Info(msg, l.appendFields(ctx, fields)...)
}
func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, l.appendFields(ctx, fields)...)
}
func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Error(msg, l.appendFields(ctx, fields)...)
}
func (l *logger) DPanic(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.DPanic(msg, l.appendFields(ctx, fields)...)
}
func (l *logger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Panic(msg, l.appendFields(ctx, fields)...)
}
func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, l.appendFields(ctx, fields)...)
}

func (l *logger) appendFields(ctx context.Context, fields []zap.Field) []zap.Field {
	addFields := []zap.Field{
		gcp.GetSpanField(ctx),
		gcp.GetTraceField(ctx, l.projectID),
	}
	f := make([]zap.Field, 0, len(fields)+len(addFields))
	f = append(f, fields...)
	f = append(f, addFields...)
	return f
}

func (l *logger) ErrorReport(ctx context.Context, err error, httpRequest *gcp.ErrorHTTPRequest, user string) {
	if l.isLocal {
		return
	}

	location := &gcp.ReportLocation{}
	pc, file, line, ok := runtime.Caller(1)
	if ok {
		location.FilePath = file
		location.LineNumber = line
		location.FunctionName = runtime.FuncForPC(pc).Name()
	}
	ectx := &gcp.ErrorContext{
		HTTPRequest:    httpRequest,
		User:           user,
		ReportLocation: location,
	}
	l.Error(
		ctx,
		fmt.Sprintf("%+v", err),
		zap.String("@type", gcp.LogTypeReportedErrorEvent),
		zap.Time("eventTime", time.Now()),
		zap.Object(ectx.Key(), ectx),
	)
}

func (*logger) EchoErrorHTTPRequest(c echo.Context) *gcp.ErrorHTTPRequest {
	req := c.Request()
	res := c.Response()
	return &gcp.ErrorHTTPRequest{
		Method:             req.Method,
		URL:                fmt.Sprintf("%s://%s%s", c.Scheme(), req.Host, req.RequestURI),
		UserAgent:          req.UserAgent(),
		Referrer:           req.Referer(),
		ResponseStatusCode: res.Status,
		RemoteIP:           c.RealIP(),
	}
}

func (*logger) GRPCErrorHTTPRequest(fullMethod string, md metadata.MD, code codes.Code) *gcp.ErrorHTTPRequest {
	var referer string
	var userAgent string
	var remoteIP string
	if md != nil {
		if vals := md.Get(constant.HeaderReferer); len(vals) > 0 {
			referer = vals[0]
		}
		if vals := md.Get(constant.HeaderUserAgent); len(vals) > 0 {
			userAgent = vals[0]
		}
		if vals := md.Get(constant.HeaderXForwardedFor); len(vals) > 0 {
			remoteIP = vals[0]
		}
	}
	return &gcp.ErrorHTTPRequest{
		Method:             http.MethodPost,
		URL:                fullMethod,
		UserAgent:          userAgent,
		Referrer:           referer,
		ResponseStatusCode: int(code),
		RemoteIP:           remoteIP,
	}
}

func (l *logger) Sync() error {
	if err := l.Logger.Sync(); err != nil {
		var pathErr *os.PathError
		if errors.As(err, &pathErr) {
			// ignore path error since stdout/stderr doesn't support fsync
			return nil
		}
		return err
	}
	return nil
}
