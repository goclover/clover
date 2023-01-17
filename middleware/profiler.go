package middleware

import (
	"expvar"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/goclover/clover"
)

// Profiler is a convenient subrouter used for mounting net/http/pprof. ie.
//
//	func MyService() http.Handler {
//	  r := clover.NewRouter()
//	  // ..middlewares
//	  r.Mount("/debug", middleware.Profiler())
//	  // ..routes
//	  return r
//	}
func Profiler() http.Handler {
	r := clover.NewRouter()
	r.Use(NoCache)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/pprof/", http.StatusMovedPermanently)
	})
	r.HandleFunc("/pprof", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/", http.StatusMovedPermanently)
	})

	r.HandleFunc("/pprof/*", pprof.Index)
	r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/pprof/profile", pprof.Profile)
	r.HandleFunc("/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/pprof/trace", pprof.Trace)
	r.HandleFunc("/vars", expVars)

	r.HandleStd("/pprof/goroutine", pprof.Handler("goroutine"))
	r.HandleStd("/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.HandleStd("/pprof/mutex", pprof.Handler("mutex"))
	r.HandleStd("/pprof/heap", pprof.Handler("heap"))
	r.HandleStd("/pprof/block", pprof.Handler("block"))
	r.HandleStd("/pprof/allocs", pprof.Handler("allocs"))

	return r
}

// Replicated from expvar.go as not public.
func expVars(w http.ResponseWriter, r *http.Request) {
	first := true
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\n")
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}
