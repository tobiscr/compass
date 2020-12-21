/*
 * Copyright 2020 The Compass Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package http

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

//TODO extract headers to forward in config
var ForwardHeaders = []string{
	"Authorization",
}

type HeaderForwarder struct {
	headers []string
}

func NewHeaderForwarder(headers []string) *HeaderForwarder {
	return &HeaderForwarder{
		headers: headers,
	}
}

// HeaderForwarder stores the specified request in the context so that they can later be used and forwarded to other backends
func (hf *HeaderForwarder) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		for _, header := range hf.headers {
			if value := r.Header.Get(header); value != "" {
				ctx = SaveToContext(ctx, header, value)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(rw, r)
	})
}

type key int

const HeadersContextKey key = iota

func LoadFromContext(ctx context.Context) (map[string]string, error) {
	value := ctx.Value(HeadersContextKey)
	headers, ok := value.(map[string]string)
	if !ok {
		return nil, errors.Errorf("headers not found in context")
	}
	return headers, nil
}

func SaveToContext(ctx context.Context, key, val string) context.Context {
	headers := make(map[string]string)
	if value := ctx.Value(HeadersContextKey); value != nil {
		if currentHeaders, ok := value.(map[string]string); ok {
			headers = currentHeaders
		}
	}
	headers[key] = val
	return context.WithValue(ctx, HeadersContextKey, headers)
}
