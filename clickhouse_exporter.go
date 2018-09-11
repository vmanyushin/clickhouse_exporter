package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/kshvakov/clickhouse"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "clickhouse"
)

var (
	bind = flag.String("listen-port", ":8080", "The address to listen on for HTTP requests.")
	user = flag.String("username", "default", "Clickhouse username")
	pass = flag.String("password", "", "User password")
	host = flag.String("hostname", "localhost", "Clickhouse hostname")
	port = flag.String("port", "9000", "Clickhouse port")
)

// Exporter define a struct for you collector that contains pointers
// to prometheus descriptors for each metric you wish to expose.
// Note you can also include fields of other types if they provide utility
// but we just won't be exposing them as metrics.
type Exporter struct {
	dsn     string
	timeout time.Duration
	version *prometheus.Desc
}

func main() {
	flag.Parse()

	dsn := "tcp://" + *host + ":" + *port + "/" + "?username=" + *user + "&password=" + *pass + "&database=system"
	prometheus.MustRegister(NewExporter(dsn, 1))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<a href=\"/metrics\">/metrics</a>")
	})

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*bind, nil))
}

// NewExporter create new struct with prometheus descriptors
func NewExporter(dsn string, timeout time.Duration) *Exporter {
	return &Exporter{
		dsn:     dsn,
		timeout: timeout,
		version: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "version"),
			"The version of this clickhouse server.",
			[]string{"version"},
			nil,
		),
	}
}

// Describe write all metric descriptors to the prometheus desc channel.
func (exporter *Exporter) Describe(ch chan<- *prometheus.Desc) {

	// Update this section with the each metric you create for a given collector
	ch <- exporter.version
}

// Collect implements required collect function for all prometheus collectors
func (exporter *Exporter) Collect(ch chan<- prometheus.Metric) {

	// collect statistic from Clichouse
	connect, err := sql.Open("clickhouse", exporter.dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer connect.Close()

	var version string
	err = connect.QueryRow("SELECT value FROM build_options WHERE name = 'VERSION_DESCRIBE'").Scan(&version)

	if err != nil {
		log.Fatal(err)
	}

	// Write latest value for each metric in the prometheus metric channel.
	// Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	ch <- prometheus.MustNewConstMetric(exporter.version, prometheus.GaugeValue, 1, version)
}
