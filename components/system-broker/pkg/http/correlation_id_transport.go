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
	"github.com/kyma-incubator/compass/components/system-broker/pkg/log"
	"net/http"
)

func NewCorrelationIDTransport(roundTripper HTTPRoundTripper, uuidService UUIDService) *CorrelationIDTransport {
	return &CorrelationIDTransport{
		roundTripper: roundTripper,
		uuidService:  uuidService,
	}
}

type CorrelationIDTransport struct {
	roundTripper HTTPRoundTripper
	uuidService  UUIDService
}

func (c *CorrelationIDTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	ctx := r.Context()
	entry := log.C(ctx)
	if correlationID := log.CorrelationIDForRequest(r, c.uuidService); correlationID != "" {
		entry = entry.WithField(log.FieldCorrelationID, correlationID)
	}
	ctx = log.ContextWithLogger(ctx, entry)
	r = r.WithContext(ctx)

	return c.roundTripper.RoundTrip(r)
}