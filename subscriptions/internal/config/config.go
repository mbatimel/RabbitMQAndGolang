package config

import (
	"io"
	"os"
	"time"

	metricsRegression "github.com/mbatimel/RabbitMQAndGolang/subscriptions/internal/metrics"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/seniorGolang/gokit/env"
)

var metrics *metricsRegression.Metrics

const formatJSON = "json"

var configuration *Config

type Config struct {
	gitSHA      *string
	version     *string
	nodeName    *string
	buildStamp  *string
	buildNumber *string
	serviceName *string

	LogLevel  string `env:"LOGGER_LEVEL" envDefault:"debug"`
	LogFormat string `env:"LOGGER_FORMAT" envDefault:""`

	ServiceBind string `env:"BIND_ADDR" envDefault:":9000" useFromEnv:"-"`
	HealthBind  string `env:"BIND_HEALTH" envDefault:":9091" useFromEnv:"-"`
	MetricsBind string `env:"BIND_METRICS" envDefault:":9090" useFromEnv:"-"`
	MetricsPath string `env:"METRICS_PATH" envDefault:"/metrics" useFromEnv:"-"`

	MaxRequestBodySize   int   `env:"MAX_REQUEST_BODY_SIZE" envDefault:"104857600"` // 100 MB
	MaxRequestHeaderSize int   `env:"MAX_REQUEST_HEADER_SIZE" envDefault:"16384"`   // 16 KB
	ReadTimeout          int64 `env:"READ_TIMEOUT" envDefault:"120"`
	Postgres             Postgres
}

// Postgres
type Postgres struct {
	Addr            string `env:"PG_ADDR"`
	DB              string `env:"PG_DB"`
	User            string `env:"PG_USER"`
	Password        string `env:"PG_PASSWORD"`
	MaxConn         int    `env:"PG_MAX_CONN" envDefault:"10"`
	MaxIdleLifetime string `env:"PG_MAX_IDLE_LIFETIME" envDefault:"30s"`
	MaxLifetime     string `env:"PG_MAX_LIFETIME" envDefault:"3m"`
	PrepareCacheCap int    `env:"PG_PREPARE_CACHE_CAP" envDefault:"128"`
	ReplicaAddr     string `env:"PG_REPLICA_ADDR"`
	UserRO          string `env:"PG_REPLICA_USER"`
	PasswordRO      string `env:"PG_REPLICA_PASSWORD"`
}

func Metrics() *metricsRegression.Metrics {
	if metrics == nil {
		metrics = metricsRegression.CreateMetrics("mbatimel", "regression", nil)
	}

	return metrics
}
func (cfg Config) Logger() (logger zerolog.Logger) {
	level := zerolog.InfoLevel
	if newLevel, err := zerolog.ParseLevel(cfg.LogLevel); err == nil {
		level = newLevel
	}
	var out io.Writer = os.Stdout
	if cfg.LogFormat != formatJSON {
		out = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.StampMicro}
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	return zerolog.New(out).Level(level).With().Timestamp().Logger()
}
func Values() Config {
	return *internalConfig()
}
func SetBuildInfo(serviceName, gitSHA, version, buildStamp, buildNumber string) {
	nodeName, _ := os.Hostname()

	conf := internalConfig()
	setLinkedString(&conf.gitSHA, gitSHA)
	setLinkedString(&conf.version, version)
	setLinkedString(&conf.nodeName, nodeName)
	setLinkedString(&conf.buildStamp, buildStamp)
	setLinkedString(&conf.buildNumber, buildNumber)
	setLinkedString(&configuration.serviceName, serviceName)
}

func internalConfig() *Config {
	if configuration == nil {

		configuration = &Config{}

		if err := env.Parse(configuration); err != nil {
			panic(err)
		}
	}
	return configuration
}
func getLinkedString(linked *string) string {
	if linked != nil {
		return *linked
	}
	return "unset"
}
func setLinkedString(linked **string, value string) {
	if *linked == nil {
		*linked = &value
	}
}

func ServiceName() string {
	return getLinkedString(internalConfig().serviceName)
}

func NodeName() string {
	return getLinkedString(internalConfig().nodeName)
}

func Version() string {
	return getLinkedString(internalConfig().version)
}

func GitSHA() string {
	return getLinkedString(internalConfig().gitSHA)
}

func BuildStamp() string {
	return getLinkedString(internalConfig().buildStamp)
}

func BuildNumber() string {
	return getLinkedString(internalConfig().buildNumber)
}
