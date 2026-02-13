package handler

import "net/http"

type Context struct {
    Request  *http.Request
    Response http.ResponseWriter
    Params   map[string]string
}