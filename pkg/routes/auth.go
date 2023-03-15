package routes

import (
	cfg "github.com/advancemg/vimb-loader/internal/config"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"net/http"
	"strings"
	"sync"
)

var cfgToken = ""
var once = sync.Once{}

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	once.Do(func() {
		err := cfg.Load()
		if err != nil {

		}
		cfgToken = cfg.Config.Token
	})
	return func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		token := r.Header.Get("token")
		if len(strings.TrimSpace(token)) == 0 || cfgToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(utils.ToJsonBytes(map[string]string{"state": "token is nil"}))
			return
		}
		if token != cfgToken {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(utils.ToJsonBytes(map[string]string{"state": "need token"}))
			return
		}
		handler.ServeHTTP(w, r)
		return
	}
}
