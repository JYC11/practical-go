package handlers

import (
	"net/http"
	"production-ready/config"
)

func Register(mux *http.ServeMux, conf config.AppConfig) {
	mux.Handle(
		"/healthz",
		&app{conf: conf, handler: healthCheckHandler},
	)
	mux.Handle(
		"/api",
		&app{conf: conf, handler: apiHandler},
	)
	mux.Handle(
		"/panic",
		&app{conf: conf, handler: panicHandler},
	)
	expensiveHandler := http.HandlerFunc(expensiveFeature)
	// mux.Handle(
	// 	"/expensive",
	// 	&app{
	// 		conf:    conf,
	// 		handler: expensiveFeatureHandler,
	// 	},
	// )
}
