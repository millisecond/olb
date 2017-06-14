package proxy

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/millisecond/olb/config"
	"github.com/millisecond/olb/logger"
	"github.com/millisecond/olb/metrics"
	"github.com/millisecond/olb/proxy/gzip"
	"github.com/millisecond/olb/uuid"
	"github.com/millisecond/olb/model"
)

// HTTPProxy is a dynamic reverse proxy for HTTP and HTTPS protocols.
type HTTPProxy struct {
	// Config is the proxy configuration as provided during startup.
	Config config.Proxy

	// Time returns the current time as the number of seconds since the epoch.
	// If Time is nil, time.Now is used.
	Time func() time.Time

	// Transport is the http connection pool configured with timeouts.
	// The proxy will panic if this value is nil.
	Transport http.RoundTripper

	// InsecureTransport is the http connection pool configured with
	// InsecureSkipVerify set. This is used for https proxies with
	// self-signed certs.
	InsecureTransport http.RoundTripper

	// Lookup returns a target host for the given request.
	// The proxy will panic if this value is nil.
	Lookup func(http.ResponseWriter, *http.Request) *model.Target

	// Requests is a timer metric which is updated for every request.
	Requests metrics.Timer

	// Noroute is a counter metric which is updated for every request
	// where Lookup() returns nil.
	Noroute metrics.Counter

	// Logger is the access logger for the requests.
	Logger logger.Logger

	// UUID returns a unique id in uuid format.
	// If UUID is nil, uuid.NewUUID() is used.
	UUID func() string
}

func (p *HTTPProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if p.Lookup == nil {
		panic("no lookup function")
	}

	target := p.Lookup(w, req)
	if target == nil {
		w.WriteHeader(p.Config.NoRouteStatus)
		return
	}

	// build the request url since req.URL will get modified
	// by the reverse proxy and contains only the RequestURI anyway
	requestURL := &url.URL{
		Scheme:   scheme(req),
		Host:     req.Host,
		Path:     req.URL.Path,
		RawQuery: req.URL.RawQuery,
	}

	// build the real target url that is passed to the proxy
	targetURL := &url.URL{
		Scheme: target.URL.Scheme,
		Host:   target.URL.Host,
		Path:   req.URL.Path,
	}
	if target.URL.RawQuery == "" || req.URL.RawQuery == "" {
		targetURL.RawQuery = target.URL.RawQuery + req.URL.RawQuery
	} else {
		targetURL.RawQuery = target.URL.RawQuery + "&" + req.URL.RawQuery
	}

	if target.Host == "dst" {
		req.Host = targetURL.Host
	}

	// TODO(fs): The HasPrefix check seems redundant since the lookup function should
	// TODO(fs): have found the target based on the prefix but there may be other
	// TODO(fs): matchers which may have different rules. I'll keep this for
	// TODO(fs): a defensive approach.
	if target.StripPath != "" && strings.HasPrefix(req.URL.Path, target.StripPath) {
		targetURL.Path = targetURL.Path[len(target.StripPath):]
	}

	if err := addHeaders(req, p.Config, target.StripPath); err != nil {
		http.Error(w, "cannot parse "+req.RemoteAddr, http.StatusInternalServerError)
		return
	}

	if p.Config.RequestID != "" {
		id := p.UUID
		if id == nil {
			id = uuid.NewUUID
		}
		req.Header.Set(p.Config.RequestID, id())
	}

	upgrade, accept := req.Header.Get("Upgrade"), req.Header.Get("Accept")

	tr := p.Transport
	if target.TLSSkipVerify {
		tr = p.InsecureTransport
	}

	var h http.Handler
	switch {
	case upgrade == "websocket" || upgrade == "Websocket":
		if targetURL.Scheme == "https" || targetURL.Scheme == "wss" {
			h = newRawProxy(targetURL, func(network, address string) (net.Conn, error) {
				return tls.Dial(network, address, tr.(*http.Transport).TLSClientConfig)
			})
		} else {
			h = newRawProxy(targetURL, net.Dial)
		}

	case accept == "text/event-stream":
		// use the flush interval for SSE (server-sent events)
		// must be > 0s to be effective
		h = newHTTPProxy(targetURL, tr, p.Config.FlushInterval)

	default:
		h = newHTTPProxy(targetURL, tr, time.Duration(0))
	}

	if p.Config.GZIPContentTypes != nil {
		h = gzip.NewGzipHandler(h, p.Config.GZIPContentTypes)
	}

	timeNow := p.Time
	if timeNow == nil {
		timeNow = time.Now
	}

	start := timeNow()
	h.ServeHTTP(w, req)
	end := timeNow()
	dur := end.Sub(start)

	if p.Requests != nil {
		p.Requests.Update(dur)
	}
	if target.Timer != nil {
		target.Timer.Update(dur)
	}

	// get response and update metrics
	rp, ok := h.(*httputil.ReverseProxy)
	if !ok {
		return
	}
	rpt, ok := rp.Transport.(*transport)
	if !ok {
		return
	}
	if rpt.resp == nil {
		return
	}
	metrics.DefaultRegistry.GetTimer(key(rpt.resp.StatusCode)).Update(dur)

	// write access log
	if p.Logger != nil {
		p.Logger.Log(&logger.Event{
			Start:           start,
			End:             end,
			Request:         req,
			Response:        rpt.resp,
			RequestURL:      requestURL,
			UpstreamAddr:    targetURL.Host,
			UpstreamService: target.Service,
			UpstreamURL:     targetURL,
		})
	}
}

func key(code int) string {
	b := []byte("http.status.")
	b = strconv.AppendInt(b, int64(code), 10)
	return string(b)
}
