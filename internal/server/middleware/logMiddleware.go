package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ResponseLogger struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func newResponseLogger(w http.ResponseWriter) *ResponseLogger {
	return &ResponseLogger{ResponseWriter: w, status: http.StatusOK, body: new(bytes.Buffer)}
}

func (lrw *ResponseLogger) WriteHeader(status int) {
	lrw.status = status
	lrw.ResponseWriter.WriteHeader(status)
}

func (lrw *ResponseLogger) Write(b []byte) (int, error) {
	lrw.body.Write(b)
	return lrw.ResponseWriter.Write(b)
}

func LogginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := newResponseLogger(w)

		logrus.Info("request started: ", r.Method, r.RequestURI)

		var body bytes.Buffer
		if r.Body != nil {
			io.Copy(&body, r.Body)
			logrus.Info("body received: ", body.String())
			r.Body = io.NopCloser(bytes.NewBuffer(body.Bytes()))
		}

		next.ServeHTTP(lrw, r)

		logrus.Info("response status code: ", lrw.status)
	})
}
