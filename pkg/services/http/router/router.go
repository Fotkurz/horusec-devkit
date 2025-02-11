// Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
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

package router

import (
	"compress/flate"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Fotkurz/horusec-devkit/pkg/enums/ozzovalidation"
	"github.com/Fotkurz/horusec-devkit/pkg/services/http/router/enums"
	"github.com/Fotkurz/horusec-devkit/pkg/utils/env"
	"github.com/Fotkurz/horusec-devkit/pkg/utils/logger"
)

type IRouter interface {
	ListenAndServe()
	GetPort() string
	GetMux() *chi.Mux
	Route(pattern string, fn func(router chi.Router)) chi.Router
}

type Router struct {
	port        string
	timeout     time.Duration
	corsOptions *cors.Options
	router      *chi.Mux
}

func NewHTTPRouter(corsOptions *cors.Options, defaultPort string) IRouter {
	router := &Router{
		port:        env.GetEnvOrDefault(enums.HorusecPort, defaultPort),
		timeout:     time.Duration(env.GetEnvOrDefaultInt(enums.HorusecRouterTimeout, ozzovalidation.Length10)) * time.Second,
		corsOptions: corsOptions,
		router:      chi.NewRouter(),
	}

	return router.setRouterConfig()
}

func (r *Router) GetMux() *chi.Mux {
	return r.router
}

func (r *Router) Route(pattern string, fn func(router chi.Router)) chi.Router {
	return r.router.Route(pattern, fn)
}

func (r *Router) ListenAndServe() {
	logger.LogInfo(fmt.Sprintf(enums.MessageServiceRunningOnPort, r.port))
	logger.LogPanic(enums.MessageListenAndServeError, http.ListenAndServe(fmt.Sprintf(":%s", r.port), r.router))
}

func (r *Router) GetPort() string {
	return r.port
}

func (r *Router) setRouterConfig() *Router {
	r.enableRealIP()
	r.enableLogger()
	r.enableRecover()
	r.enableTimeout()
	r.enableCompress()
	r.enableRequestID()
	r.enableCORS()
	r.routeMetrics()

	return r
}

func (r *Router) enableRealIP() {
	r.router.Use(middleware.RealIP)
}

func (r *Router) enableLogger() {
	r.router.Use(middleware.Logger)
}

func (r *Router) enableRecover() {
	r.router.Use(middleware.Recoverer)
}

func (r *Router) enableTimeout() {
	r.router.Use(middleware.Timeout(r.timeout))
}

func (r *Router) enableCompress() {
	r.router.Use(middleware.Compress(flate.BestCompression))
}

func (r *Router) enableRequestID() {
	r.router.Use(middleware.RequestID)
}

func (r *Router) enableCORS() {
	r.router.Use(r.getCorsHandler)
}

func (r *Router) routeMetrics() {
	r.router.Handle("/metrics", promhttp.Handler())
}

func (r *Router) getCorsHandler(next http.Handler) http.Handler {
	return cors.New(*r.corsOptions).Handler(next)
}
