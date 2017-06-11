package config

import (
	"net/http"
	"regexp"
	"time"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

type Config struct {
	Proxy       Proxy
	Registry    Registry
	Listen      []Listen
	Log         Log
	Metrics     Metrics
	UI          UI
	Runtime     Runtime
	AWSConfig	AWSConfig
	ProfileMode string
	ProfilePath string
}

type AWSConfig struct {
	Region string
	Endpoint string
	DynamoTableName string
	Key string
	Secret string
}

func (cfg *AWSConfig) AWSConfig() *aws.Config {
	c := &aws.Config{}
	if len(cfg.Region) > 0 {
		c.Region = aws.String(cfg.Region)
	}
	if len(cfg.Endpoint) > 0 {
		c.Endpoint = aws.String(cfg.Endpoint)
	}
	// Only set credentials if key and secret are given, otherwise fall back to IAM role
	if len(cfg.Key) > 0 && len(cfg.Secret) > 0 {
		c.Credentials = credentials.NewStaticCredentials(cfg.Key, cfg.Secret, "")
	}
	return c
}

type CertSource struct {
	Name         string
	Type         string
	CertPath     string
	KeyPath      string
	ClientCAPath string
	CAUpgradeCN  string
	Refresh      time.Duration
	Header       http.Header
}

type Listen struct {
	Addr          string
	Proto         string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	CertSource    CertSource
	StrictMatch   bool
	TLSMinVersion uint16
	TLSMaxVersion uint16
	TLSCiphers    []uint16
}

type UI struct {
	Listen Listen
	Color  string
	Title  string
	Access string
}

type Proxy struct {
	Strategy              string
	Matcher               string
	NoRouteStatus         int
	MaxConn               int
	ShutdownWait          time.Duration
	DialTimeout           time.Duration
	ResponseHeaderTimeout time.Duration
	KeepAliveTimeout      time.Duration
	FlushInterval         time.Duration
	LocalIP               string
	ClientIPHeader        string
	TLSHeader             string
	TLSHeaderValue        string
	GZIPContentTypes      *regexp.Regexp
	RequestID             string
}

type Runtime struct {
	GOGC       int
	GOMAXPROCS int
}

type Circonus struct {
	APIKey   string
	APIApp   string
	APIURL   string
	CheckID  string
	BrokerID string
}

type Log struct {
	AccessFormat string
	AccessTarget string
	RoutesFormat string
}

type Metrics struct {
	Target       string
	Prefix       string
	Names        string
	Interval     time.Duration
	GraphiteAddr string
	StatsDAddr   string
	Circonus     Circonus
}

type Registry struct {
	Backend string
	Static  Static
	File    File
	Timeout time.Duration
	Retry   time.Duration
}

type Static struct {
	Routes string
}

type File struct {
	Path string
}

type Consul struct {
	Addr               string
	Scheme             string
	Token              string
	KVPath             string
	TagPrefix          string
	Register           bool
	ServiceAddr        string
	ServiceName        string
	ServiceTags        []string
	ServiceStatus      []string
	CheckInterval      time.Duration
	CheckTimeout       time.Duration
	CheckScheme        string
	CheckTLSSkipVerify bool
}
