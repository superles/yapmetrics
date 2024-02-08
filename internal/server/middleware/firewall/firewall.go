package firewall

import (
	"github.com/superles/yapmetrics/internal/utils/logger"
	"github.com/superles/yapmetrics/internal/utils/network"
	"net/http"
)

func WithTrustedSubnet(subnet string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		logFn := func(w http.ResponseWriter, r *http.Request) {
			if len(subnet) == 0 {
				h.ServeHTTP(w, r)
				return
			}
			realIP := r.Header.Get("X-Real-Ip")
			if len(realIP) == 0 {
				logger.Log.Error("address not exist")
				http.Error(w, "address not exist in header", http.StatusUnauthorized)
				return
			}

			if inNetwork, err := network.IsAddressInNetwork(realIP, subnet); err != nil {
				logger.Log.Error(err)
				http.Error(w, "wrong ip or trusted subnet", http.StatusUnauthorized)
				return
			} else if !inNetwork {
				logger.Log.Error("address not exist in trusted network")
				http.Error(w, "address not exist in trusted network", http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(logFn)
	}
}
