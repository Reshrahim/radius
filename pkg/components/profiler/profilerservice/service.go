/*
Copyright 2023 The Radius Authors.

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

package profilerservice

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/radius-project/radius/pkg/ucp/ucplog"
)

// Options represents the options for enabling pprof profiler.
type Options struct {
	// Enabled is a flag to enable the profiler.
	Enabled bool `yaml:"enabled,omitempty"`

	// Port is the port on which the profiler server listens.
	Port int `yaml:"port,omitempty"`
}

// Service is the profiler service.
type Service struct {
	Options *Options
}

// Name returns the name of the profiler service.
func (s *Service) Name() string {
	return "pprof profiler"
}

// Run starts the profiler server that exposes an endpoint to collect profiler from. It
// handles shutdown based on the context, and returns an error if the server fails to start.
func (s *Service) Run(ctx context.Context) error {
	logger := ucplog.FromContextOrDiscard(ctx)

	profilerPort := strconv.Itoa(s.Options.Port)
	server := &http.Server{
		Addr: ":" + profilerPort,
		BaseContext: func(ln net.Listener) context.Context {
			return ctx
		},
	}

	// Handle shutdown based on the context
	go func() {
		<-ctx.Done()
		// We don't care about shutdown errors
		_ = server.Shutdown(ctx)
	}()

	logger.Info(fmt.Sprintf("profiler Server listening on localhost port: '%s'...", profilerPort))
	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		// We expect this, safe to ignore.
		logger.Info("Server stopped...")
		return nil
	} else if err != nil {
		return err
	}

	logger.Info("Server stopped...")
	return nil
}
