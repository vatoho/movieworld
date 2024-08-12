package middleware

import (
	"log"
	"net/http"
)

func ErrorLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger, err := GetLoggerFromContext(r.Context())
		if err != nil {
			log.Printf("can not get logger from context: %s", err)
			WriteNoLoggerResponse(w)
		}
		recorder := &responseRecorder{
			ResponseWriter: w,
		}
		next.ServeHTTP(recorder, r)
		if recorder.statusCode >= 400 {
			logger.Errorf("error occurred: %s", string(recorder.respBody))
		}
	})
}

type responseRecorder struct {
	statusCode int
	respBody   []byte
	http.ResponseWriter
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.respBody = b
	return r.ResponseWriter.Write(b)
}

func (r *responseRecorder) Header() http.Header {
	return r.ResponseWriter.Header()
}
