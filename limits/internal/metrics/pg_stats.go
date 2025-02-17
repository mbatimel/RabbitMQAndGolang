package metrics

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "pgx_stats"
	subsystem = "conns"
)

type StatsGetter interface {
	Stat() *pgxpool.Stat
}

type PgStatsCollector struct {
	sg                       StatsGetter
	acquireCountDesc         *prometheus.Desc
	acquireDurationDesc      *prometheus.Desc
	acquiredConnsDesc        *prometheus.Desc
	canceledAcquireCountDesc *prometheus.Desc
	constructingConnsDesc    *prometheus.Desc
	emptyAcquireCountDesc    *prometheus.Desc
	idleConnsDesc            *prometheus.Desc
	maxConnsDesc             *prometheus.Desc
	totalConnsDesc           *prometheus.Desc
}

func NewPGStatsCollector(conn, dbName string, sg StatsGetter) *PgStatsCollector {
	labels := prometheus.Labels{
		"conn":    conn,
		"db_name": dbName,
	}
	return &PgStatsCollector{
		sg: sg,
		acquireCountDesc: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				subsystem, "acquire_count",
			),
			"returns the cumulative count of successful acquires from the pool",
			nil,
			labels,
		),
		acquireDurationDesc: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				subsystem, "acquire_duration",
			),
			"returns the total duration of all successful acquires from the pool",
			nil,
			labels,
		),
		acquiredConnsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				subsystem, "acquired_conns",
			),
			"returns the number of currently acquired connections in the pool",
			nil,
			labels,
		),
		canceledAcquireCountDesc: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				subsystem, "canceled_acquire_count",
			),
			"returns the cumulative count of acquires from the pool that were canceled by a context",
			nil,
			labels,
		),
		constructingConnsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				subsystem, "constructing_conns",
			),
			"returns the number of conns with construction in progress in the pool",
			nil,
			labels,
		),
		emptyAcquireCountDesc: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				subsystem, "empty_acquire_count",
			),
			"returns the cumulative count of successful acquires from the pool that waited for a resource to be released or constructed because the pool was empty",
			nil,
			labels,
		),
		idleConnsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				subsystem, "idle_conns",
			),
			"returns the number of currently idle conns in the pool",
			nil,
			labels,
		),
		maxConnsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				subsystem, "max_conns",
			),
			"returns the maximum size of the pool",
			nil,
			labels,
		),
		totalConnsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				subsystem, "total_conns",
			),
			"returns the total number of resources currently in the pool. the value is the sum of constructing_conns, acquired_conns, and idle_conns",
			nil,
			labels,
		),
	}
}

func (c PgStatsCollector) Collect(ch chan<- prometheus.Metric) {
	var stats = c.sg.Stat()
	ch <- prometheus.MustNewConstMetric(
		c.acquireCountDesc,
		prometheus.GaugeValue,
		float64(stats.AcquireCount()),
	)
	ch <- prometheus.MustNewConstMetric(
		c.acquireDurationDesc,
		prometheus.GaugeValue,
		float64(stats.AcquireDuration()),
	)
	ch <- prometheus.MustNewConstMetric(
		c.acquiredConnsDesc,
		prometheus.GaugeValue,
		float64(stats.AcquiredConns()),
	)
	ch <- prometheus.MustNewConstMetric(
		c.canceledAcquireCountDesc,
		prometheus.GaugeValue,
		float64(stats.CanceledAcquireCount()),
	)
	ch <- prometheus.MustNewConstMetric(
		c.constructingConnsDesc,
		prometheus.GaugeValue,
		float64(stats.ConstructingConns()),
	)
	ch <- prometheus.MustNewConstMetric(
		c.emptyAcquireCountDesc,
		prometheus.GaugeValue,
		float64(stats.EmptyAcquireCount()),
	)
	ch <- prometheus.MustNewConstMetric(
		c.idleConnsDesc,
		prometheus.GaugeValue,
		float64(stats.IdleConns()),
	)
	ch <- prometheus.MustNewConstMetric(
		c.maxConnsDesc,
		prometheus.GaugeValue,
		float64(stats.MaxConns()),
	)
	ch <- prometheus.MustNewConstMetric(
		c.totalConnsDesc,
		prometheus.GaugeValue,
		float64(stats.TotalConns()),
	)
}

func (c PgStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.acquireCountDesc
	ch <- c.acquireDurationDesc
	ch <- c.acquiredConnsDesc
	ch <- c.canceledAcquireCountDesc
	ch <- c.constructingConnsDesc
	ch <- c.emptyAcquireCountDesc
	ch <- c.idleConnsDesc
	ch <- c.maxConnsDesc
	ch <- c.totalConnsDesc
}
