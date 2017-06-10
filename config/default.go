package config

import (
	"os"
	"runtime"
	"time"
)

var defaultValues = struct {
	ListenerValue         string
	CertSourcesValue      string
	ReadTimeout           time.Duration
	WriteTimeout          time.Duration
	UIListenerValue       string
	GZIPContentTypesValue string
}{
	ListenerValue:   ":9999",
	UIListenerValue: ":9998",
}

var defaultConfig = &Config{
	ProfilePath: os.TempDir(),
	Log: Log{
		AccessFormat: "common",
		RoutesFormat: "delta",
	},
	Metrics: Metrics{
		Prefix:   "{{clean .Hostname}}.{{clean .Exec}}",
		Names:    "{{clean .Service}}.{{clean .Host}}.{{clean .Path}}.{{clean .TargetURL.Host}}",
		Interval: 30 * time.Second,
		Circonus: Circonus{
			APIApp: "olb",
		},
	},
	Proxy: Proxy{
		MaxConn:       10000,
		Strategy:      "rnd",
		Matcher:       "prefix",
		NoRouteStatus: 404,
		DialTimeout:   30 * time.Second,
		FlushInterval: time.Second,
		LocalIP:       LocalIPString(),
	},
	Registry: Registry{
		Backend: "dynamo",
		Timeout: 10 * time.Second,
		Retry:   500 * time.Millisecond,
	},
	Runtime: Runtime{
		GOGC:       800,
		GOMAXPROCS: runtime.NumCPU(),
	},
	UI: UI{
		Listen: Listen{
			Addr:  ":9998",
			Proto: "http",
		},
		Color:  "light-green",
		Access: "rw",
	},
}
