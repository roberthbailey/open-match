// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"context"
	"net/http"
	"sync"

	"net"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"open-match.dev/open-match/internal/monitoring"
	"open-match.dev/open-match/internal/util/netlistener"
)

var (
	insecureLogger = logrus.WithFields(logrus.Fields{
		"app":       "openmatch",
		"component": "insecure_server",
	})
)

type insecureServer struct {
	grpcLh       *netlistener.ListenerHolder
	grpcListener net.Listener
	grpcServer   *grpc.Server

	httpLh       *netlistener.ListenerHolder
	httpListener net.Listener
	httpMux      *http.ServeMux
	proxyMux     *runtime.ServeMux
	httpServer   *http.Server
}

func (s *insecureServer) start(params *ServerParams) (func(), error) {
	var serverStartWaiter sync.WaitGroup

	s.httpMux = params.ServeMux
	s.proxyMux = runtime.NewServeMux()

	// Configure the gRPC server.
	grpcListener, err := s.grpcLh.Obtain()
	if err != nil {
		return func() {}, errors.WithStack(err)
	}
	s.grpcListener = grpcListener
	s.grpcServer = grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))
	// Bind gRPC handlers
	for _, handlerFunc := range params.handlersForGrpc {
		handlerFunc(s.grpcServer)
	}

	serverStartWaiter.Add(1)
	go func() {
		serverStartWaiter.Done()
		gErr := s.grpcServer.Serve(s.grpcListener)
		if gErr != nil {
			return
		}
	}()

	// Configure the HTTP proxy server.
	httpListener, err := s.httpLh.Obtain()
	if err != nil {
		return func() {}, errors.WithStack(err)
	}
	s.httpListener = httpListener

	// Bind gRPC handlers
	ctx, cancel := context.WithCancel(context.Background())

	for _, handlerFunc := range params.handlersForGrpcProxy {
		if err = handlerFunc(ctx, s.proxyMux, grpcListener.Addr().String(), []grpc.DialOption{grpc.WithInsecure()}); err != nil {
			cancel()
			return func() {}, errors.WithStack(err)
		}
	}

	s.httpMux.HandleFunc("/healthz", monitoring.NewHealthProbe(params.handlersForHealthCheck))
	s.httpMux.Handle("/", s.proxyMux)
	s.httpServer = &http.Server{
		Addr:    s.httpListener.Addr().String(),
		Handler: s.httpMux,
	}
	serverStartWaiter.Add(1)
	go func() {
		serverStartWaiter.Done()
		hErr := s.httpServer.Serve(s.httpListener)
		defer cancel()
		if hErr != nil {
			insecureLogger.Debugf("error closing gRPC server: %s", hErr)
		}
	}()

	return serverStartWaiter.Wait, nil
}

func (s *insecureServer) stop() {
	s.grpcServer.Stop()
	if err := s.grpcListener.Close(); err != nil {
		insecureLogger.Debugf("error closing gRPC listener: %s", err)
	}

	if err := s.httpServer.Close(); err != nil {
		insecureLogger.Debugf("error closing HTTP server: %s", err)
	}

	if err := s.httpListener.Close(); err != nil {
		insecureLogger.Debugf("error closing HTTP listener: %s", err)
	}
}

func newInsecureServer(grpcLh *netlistener.ListenerHolder, httpLh *netlistener.ListenerHolder) *insecureServer {
	return &insecureServer{
		grpcLh: grpcLh,
		httpLh: httpLh,
	}
}
