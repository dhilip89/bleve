//  Copyright (c) 2014 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"fmt"
	"net/http"
)

type ListFieldsHandler struct {
	defaultIndexName string
	IndexNameLookup  varLookupFunc
}

func NewListFieldsHandler(defaultIndexName string) *ListFieldsHandler {
	return &ListFieldsHandler{
		defaultIndexName: defaultIndexName,
	}
}

func (h *ListFieldsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// find the index to operate on
	var indexName string
	if h.IndexNameLookup != nil {
		indexName = h.IndexNameLookup(req)
	}
	if indexName == "" {
		indexName = h.defaultIndexName
	}
	index := IndexByName(indexName)
	if index == nil {
		showError(w, req, fmt.Sprintf("no such index '%s'", indexName), 404)
		return
	}

	fields, err := index.Fields()
	if err != nil {
		showError(w, req, fmt.Sprintf("error: %v", err), 500)
		return
	}

	fieldsResponse := struct {
		Fields []string `json:"fields"`
	}{
		Fields: fields,
	}

	// encode the response
	mustEncode(w, fieldsResponse)
}
