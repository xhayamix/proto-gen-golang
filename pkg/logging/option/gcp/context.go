package gcp

import (
	"go.uber.org/zap/zapcore"
)

// https://cloud.google.com/logging/docs/agent/configuration#special-fields
const (
	sourceLocationKey = "logging.googleapis.com/SourceLocation"
	httpRequestKey    = "httpRequest"
	serviceContextKey = "serviceContext"
	errorContextKey   = "context"
	reportLocationKey = "reportLocation"
)

// SourceLocation https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogEntrySourceLocation
type SourceLocation struct {
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
	Function string `json:"function,omitempty"`
}

func (s *SourceLocation) MarshalLogObject(e zapcore.ObjectEncoder) error {
	e.AddString("file", s.File)
	e.AddInt("line", s.Line)
	e.AddString("function", s.Function)
	return nil
}

func (*SourceLocation) Key() string {
	return sourceLocationKey
}

// HTTPRequest https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#HttpRequest
type HTTPRequest struct {
	RequestMethod string    `json:"requestMethod,omitempty"`
	RequestURL    string    `json:"requestUrl,omitempty"`
	RequestSize   string    `json:"requestSize,omitempty"`
	Status        int       `json:"status,omitempty"`
	ResponseSize  string    `json:"responseSize,omitempty"`
	UserAgent     string    `json:"userAgent,omitempty"`
	RemoteIP      string    `json:"remoteIp,omitempty"`
	ServerIP      string    `json:"serverIp,omitempty"`
	Referer       string    `json:"referer,omitempty"`
	Latency       *Duration `json:"latency,omitempty"`
	// CacheLookup                    bool   `json:"cacheLookup"`
	// CacheHit                       bool   `json:"cacheHit"`
	// CacheValidatedWithOriginServer bool   `json:"cacheValidatedWithOriginServer"`
	// CacheFillBytes                 int    `json:"cacheFillBytes"`
	Protocol string `json:"protocol,omitempty"`
}

func (h *HTTPRequest) MarshalLogObject(e zapcore.ObjectEncoder) error {
	if h.RequestMethod != "" {
		e.AddString("requestMethod", h.RequestMethod)
	}
	if h.RequestURL != "" {
		e.AddString("requestUrl", h.RequestURL)
	}
	if h.RequestSize != "" {
		e.AddString("requestSize", h.RequestSize)
	}
	if h.Status != 0 {
		e.AddInt("status", h.Status)
	}
	if h.ResponseSize != "" {
		e.AddString("responseSize", h.ResponseSize)
	}
	if h.UserAgent != "" {
		e.AddString("userAgent", h.UserAgent)
	}
	if h.RemoteIP != "" {
		e.AddString("remoteIp", h.RemoteIP)
	}
	if h.ServerIP != "" {
		e.AddString("serverIp", h.ServerIP)
	}
	if h.Referer != "" {
		e.AddString("referer", h.Referer)
	}
	if h.Latency != nil {
		// エラー無視
		_ = e.AddObject("latency", h.Latency)
	}
	// e.AddBool("cacheLookup", h.CacheLookup)
	// e.AddBool("cacheHit", h.CacheHit)
	// e.AddBool("cacheValidatedWithOriginServer", h.CacheValidatedWithOriginServer)
	// e.AddInt("cacheFillBytes", h.CacheFillBytes)
	if h.Protocol != "" {
		e.AddString("protocol", h.Protocol)
	}
	return nil
}

func (*HTTPRequest) Key() string {
	return httpRequestKey
}

// Duration https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#google.protobuf.Duration
type Duration struct {
	Seconds int64 `json:"seconds"`
	Nanos   int32 `json:"nanos"`
}

func (d *Duration) MarshalLogObject(e zapcore.ObjectEncoder) error {
	e.AddInt64("seconds", d.Seconds)
	e.AddInt32("nanos", d.Nanos)
	return nil
}

// LogTypeReportedErrorEvent https://cloud.google.com/error-reporting/docs/formatting-error-messages#json_representation
const LogTypeReportedErrorEvent = "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent"

// ServiceContext https://cloud.google.com/error-reporting/reference/rest/v1beta1/ServiceContext
type ServiceContext struct {
	Service string `json:"service"`
	Version string `json:"version"`
}

func (s *ServiceContext) MarshalLogObject(e zapcore.ObjectEncoder) error {
	e.AddString("service", s.Service)
	e.AddString("version", s.Version)
	return nil
}

func (*ServiceContext) Key() string {
	return serviceContextKey
}

// ErrorContext https://cloud.google.com/error-reporting/reference/rest/v1beta1/ErrorContext
type ErrorContext struct {
	HTTPRequest    *ErrorHTTPRequest `json:"httpRequest"`
	User           string            `json:"user"`
	ReportLocation *ReportLocation   `json:"reportLocation"`
}

func (c *ErrorContext) MarshalLogObject(e zapcore.ObjectEncoder) error {
	if c.HTTPRequest != nil {
		if err := e.AddObject(c.HTTPRequest.Key(), c.HTTPRequest); err != nil {
			return err
		}
	}
	e.AddString("user", c.User)
	if c.ReportLocation != nil {
		if err := e.AddObject(c.ReportLocation.Key(), c.ReportLocation); err != nil {
			return err
		}
	}
	return nil
}

func (*ErrorContext) Key() string {
	return errorContextKey
}

// ErrorHTTPRequest https://cloud.google.com/error-reporting/reference/rest/v1beta1/ErrorContext#HttpRequestContext
type ErrorHTTPRequest struct {
	Method             string `json:"method,omitempty"`
	URL                string `json:"url,omitempty"`
	UserAgent          string `json:"userAgent,omitempty"`
	Referrer           string `json:"referrer,omitempty"`
	ResponseStatusCode int    `json:"responseStatusCode,omitempty"`
	RemoteIP           string `json:"remoteIp,omitempty"`
}

func (r *ErrorHTTPRequest) MarshalLogObject(e zapcore.ObjectEncoder) error {
	if r.Method != "" {
		e.AddString("method", r.Method)
	}
	if r.URL != "" {
		e.AddString("url", r.URL)
	}
	if r.UserAgent != "" {
		e.AddString("userAgent", r.UserAgent)
	}
	if r.Referrer != "" {
		e.AddString("referrer", r.Referrer)
	}
	if r.ResponseStatusCode != 0 {
		e.AddInt("responseStatusCode", r.ResponseStatusCode)
	}
	if r.RemoteIP != "" {
		e.AddString("remoteIp", r.RemoteIP)
	}
	return nil
}

func (*ErrorHTTPRequest) Key() string {
	return httpRequestKey
}

// ReportLocation https://cloud.google.com/error-reporting/reference/rest/v1beta1/ErrorContext#SourceLocation
type ReportLocation struct {
	FilePath     string `json:"filePath"`
	LineNumber   int    `json:"lineNumber"`
	FunctionName string `json:"functionName"`
}

func (s *ReportLocation) MarshalLogObject(e zapcore.ObjectEncoder) error {
	e.AddString("filePath", s.FilePath)
	e.AddInt("lineNumber", s.LineNumber)
	e.AddString("functionName", s.FunctionName)
	return nil
}

func (*ReportLocation) Key() string {
	return reportLocationKey
}
