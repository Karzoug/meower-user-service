package middleware

import "net/http"

type codeWriter struct {
	http.ResponseWriter
	code int
}

func newWrapResponseWriter(w http.ResponseWriter) *codeWriter {
	return &codeWriter{w, http.StatusOK}
}

func (w *codeWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *codeWriter) Status() int {
	return w.code
}
