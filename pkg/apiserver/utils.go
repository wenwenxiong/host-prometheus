/*
Copyright 2020 The KubeSphere Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package apiserver

import (
	"log"
	"net/http"
	"runtime"
	"strings"
)

// Avoid emitting errors that look like valid HTML. Quotes are okay.
var sanitizer = strings.NewReplacer(`&`, "&amp;", `<`, "&lt;", `>`, "&gt;")

func HandleInternalError(response http.ResponseWriter,  request *http.Request, err error) {
	_, fn, line, _ := runtime.Caller(1)
	log.Printf("%s:%d %v", fn, line, err)
	http.Error(response, sanitizer.Replace(err.Error()), http.StatusInternalServerError)
}

// HandleBadRequest writes http.StatusBadRequest and log error
func HandleBadRequest(response http.ResponseWriter, request *http.Request, err error) {
	_, fn, line, _ := runtime.Caller(1)
	log.Printf("%s:%d %v", fn, line, err)
	http.Error(response, sanitizer.Replace(err.Error()), http.StatusBadRequest)
}

func HandleNotFound(response http.ResponseWriter, request *http.Request, err error) {
	_, fn, line, _ := runtime.Caller(1)
	log.Printf("%s:%d %v", fn, line, err)
	http.Error(response, sanitizer.Replace(err.Error()), http.StatusNotFound)
}

func HandleForbidden(response http.ResponseWriter, request *http.Request, err error) {
	_, fn, line, _ := runtime.Caller(1)
	log.Printf("%s:%d %v", fn, line, err)
	http.Error(response, sanitizer.Replace(err.Error()), http.StatusForbidden)
}

func HandleConflict(response http.ResponseWriter, request *http.Request, err error) {
	_, fn, line, _ := runtime.Caller(1)
	log.Printf("%s:%d %v", fn, line, err)
	http.Error(response, sanitizer.Replace(err.Error()), http.StatusConflict)
}
