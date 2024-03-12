package logging

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Entry
}

type ContextKey string

const (
	LogLevel      ContextKey = "logLevel"
	RequestID     ContextKey = "requestId"
	TransactionID ContextKey = "transactionId"
	CommitInfo    ContextKey = "commitInfo"
	ContextValues ContextKey = "values"
)

type CtxValues struct {
	LogLevel      string
	RequestID     string
	TransactionID string
	CommitInfo    string
}

const (
	invalidArgument           = "Invalid Argument %s: %v"
	systemError               = "A system error has occurred"
	marshalingError           = "Marshaling Error"
	dataNotFound              = "%v Not Found"
	partialResponse           = "Partial Response"
	responseError             = "Response Error: %v"
	requestError              = "Request Error: %v"
	constructHTTPRequestError = "Failed to construct the HTTP request"
	httpRequestFailed         = "Http Request Failed"
	missingConfig             = "Unable to retrieve the config value"
	configError               = "Unable to load config :%v\n"
	fullRequest               = "Full Request Completed"
	childRequest              = "%v Request Completed"
)

// New creates a new instance of Logger with the context.
func New(ctx context.Context) *Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
	logger := logrus.NewEntry(log)

	if ctx != nil { //nolint:nestif
		if logLevel, ok := ctx.Value(LogLevel).(string); ok {
			if level, err := logrus.ParseLevel(logLevel); err == nil {
				log.SetLevel(level)
			}
		}

		if reqID, ok := ctx.Value(RequestID).(string); ok {
			logger = logger.WithField("requestId", reqID)
		}

		if tranID, ok := ctx.Value(TransactionID).(string); ok {
			logger = logger.WithField("transactionId", tranID)
		}

		if cInfo, ok := ctx.Value(CommitInfo).(string); ok {
			logger = logger.WithField("commitInfo", cInfo)
		}
	}

	return &Logger{logger}
}

// SetLoggingFormat sets the logging timestamp format.
func SetLoggingFormat(level logrus.Level) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(level)
}

// SetLogging Creates Context for new Logger.
func SetLogging(logLevel string) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, CommitInfo, os.Getenv("GIT_COMMIT"))
	level, err := logrus.ParseLevel(logLevel)

	if err != nil {
		level = logrus.ErrorLevel
		logrus.Error("unable to parse logging level from config, defaulting to error level")
	}
	ctx = context.WithValue(ctx, LogLevel, level)
}

func FieldsFromCTX(ctx context.Context) (fields logrus.Fields) {
	fields = make(map[string]interface{})

	if ctx != nil {
		if reqID, ok := ctx.Value(ContextValues).(CtxValues); ok {
			fields["transactionId"] = reqID.TransactionID
			fields["requestId"] = reqID.RequestID
			fields["commitInfo"] = reqID.CommitInfo
		}
	}

	return fields
}

func (l *Logger) ConfigError(loadErr string) {
	l.Error(configError, loadErr)
}

func (l *Logger) RequestError(source interface{}, err error, status interface{}, start, end time.Time) {
	l.WithFields(logrus.Fields{
		"source":     source,
		"httpStatus": status,
		"start":      start,
		"end":        end,
		"duration":   getDuration(start, end),
	}).Errorf(requestError, err.Error())
}

func (l *Logger) ResponseError(source, rootCause, status interface{}, start, end time.Time) {
	l.WithFields(logrus.Fields{
		"source":     source,
		"httpStatus": status,
		"start":      start,
		"end":        end,
		"duration":   getDuration(start, end),
	}).Errorf(responseError, rootCause)
}

func (l *Logger) DataNotFound(data, source interface{}, start, end time.Time) {
	l.WithFields(logrus.Fields{
		"source":     source,
		"httpStatus": http.StatusNotFound,
		"start":      start,
		"end":        end,
		"duration":   getDuration(start, end),
	}).Infof(dataNotFound, data)
}

func (l *Logger) PartialResponse(source interface{}, start, end time.Time) {
	l.WithFields(logrus.Fields{
		"source":     source,
		"httpStatus": http.StatusPartialContent,
		"start":      start,
		"end":        end,
		"duration":   getDuration(start, end),
	}).Warnf(partialResponse)
}

func (l *Logger) InvalidArg(argName, argValue interface{}) {
	l.Warnf(invalidArgument, argName, argValue)
}

func (l *Logger) SystemError() {
	l.Error(systemError)
}

func (l *Logger) MarshalError(obj interface{}) {
	l.Errorf("%v: %v", marshalingError, obj)
}

func (l *Logger) ConstructHTTPRequestError(method, url string, body interface{}) {
	l.WithFields(logrus.Fields{
		"method": method,
		"url":    url,
		"body":   body,
	}).Error(constructHTTPRequestError)
}

func (l *Logger) HTTPRequestError(method, url string, body, status interface{}, err error) {
	l.WithFields(logrus.Fields{
		"method":     method,
		"url":        url,
		"body":       body,
		"httpStatus": status,
		"error":      err.Error(),
	}).Error(httpRequestFailed)
}

func (l *Logger) MissingConfig(configName string) {
	l.WithField("configName", configName).Error(missingConfig)
}

func (l *Logger) RequestCompleted(source interface{}, start, end time.Time) {
	l.WithFields(logrus.Fields{
		"end":      end,
		"start":    start,
		"duration": getDuration(start, end),
		"source":   source,
	}).Info(fullRequest)
}

func (l *Logger) RequestCompletedV2(reqArgs, status interface{}, start, end time.Time) {
	l.WithFields(logrus.Fields{
		"end":        end,
		"start":      start,
		"httpStatus": status,
		"duration":   getDuration(start, end),
		"reqArgs":    reqArgs,
	}).Info(fullRequest)
}

func (l *Logger) ChildRequestCompleted(source interface{}, start time.Time, duration time.Duration) {
	l.WithFields(logrus.Fields{
		"start":    start,
		"duration": duration,
		"source":   source,
	}).Tracef(childRequest, source)
}

func getDuration(start, end time.Time) time.Duration {
	return end.Sub(start) / time.Millisecond
}
