/*
Copyright (c) 2019 Vladimir Glafirov

MIT License

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	iex "github.com/goinvest/iexcloud"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	namespace = "iexcloud"
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last query of iexcloud successful.",
		nil, nil,
	)

	stockPrice = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "stock_price"),
		"Current stock price",
		[]string{"symbol"},
		nil,
	)
)

type promHTTPLogger struct {
	logger log.Logger
}

func (l promHTTPLogger) Println(v ...interface{}) {
	level.Error(l.logger).Log("msg", fmt.Sprint(v...))
}

type iexcloudOpts struct {
	endpoint   string
	apiToken   string
	apiVersion string
}

func (o iexcloudOpts) String() string {
	return fmt.Sprintf("Endpoint: %s\n API version: %s", o.endpoint, o.apiVersion)
}

// Describe describes all the metrics ever exported by the Consul exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- stockPrice
}

// Collect fetches the stats from configured Consul location and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ok := e.collectPriceMetric(ch)
	level.Info(e.logger).Log("msg", "Collecting Price metric", "result", ok)

	if ok {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 1.0,
		)
	} else {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0.0,
		)
	}
}

func (e *Exporter) collectPriceMetric(ch chan<- prometheus.Metric) bool {
	price, err := e.client.Price("MSFT")
	if err != nil {
		level.Error(e.logger).Log("msg", "Can't query IEX Cloud API", "err", err)
		return false
	}
	ch <- prometheus.MustNewConstMetric(
		stockPrice, prometheus.GaugeValue, float64(price), "MSFT",
	)
	return true
}

// Exporter object
type Exporter struct {
	client   *iex.Client
	kvPrefix string
	kvFilter *regexp.Regexp
	logger   log.Logger
}

// NewExporter returns an initialized Exporter.
func NewExporter(opts iexcloudOpts, kvPrefix, kvFilter string, logger log.Logger) (*Exporter, error) {
	endpoint := opts.endpoint
	if !strings.Contains(endpoint, "://") {
		endpoint = "https://" + endpoint + "/" + opts.apiVersion + "/"
	}
	e, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid iexcloud endpoint: %s", err)
	}

	level.Info(logger).Log("msg", "Initializing endpoint", "endpoint", e)

	client := iex.NewClient(opts.apiToken, e.String())

	// Init our exporter.
	return &Exporter{
		client:   client,
		kvPrefix: kvPrefix,
		kvFilter: regexp.MustCompile(kvFilter),
		logger:   logger,
	}, nil
}

func init() {
	prometheus.MustRegister(version.NewCollector("iexcloud_exporter"))
}

func main() {
	var (
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9107").String()
		metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
		kvPrefix      = kingpin.Flag("kv.prefix", "Prefix from which to expose key/value pairs.").Default("").String()
		kvFilter      = kingpin.Flag("kv.filter", "Regex that determines which keys to expose.").Default(".*").String()

		opts = iexcloudOpts{}
	)

	kingpin.Flag("iexcloud.api_token", "API Token for IEX Cloud account").Required().StringVar(&opts.apiToken)
	kingpin.Flag("iexcloud.endpoint", "IEX Cloud API endpoint").Default("sandbox.iexapis.com").StringVar(&opts.endpoint)
	kingpin.Flag("iexcloud.api_version", "IEX Cloud API version").Default("stable").StringVar(&opts.apiVersion)

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting iexcloud_exporter", "version", version.Info())
	level.Info(logger).Log("build_context", version.BuildContext())

	exporter, err := NewExporter(opts, *kvPrefix, *kvFilter, logger)
	if err != nil {
		level.Error(logger).Log("msg", "Error creating the exporter", "err", err)
		os.Exit(1)
	}
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath,
		promhttp.InstrumentMetricHandler(
			prometheus.DefaultRegisterer,
			promhttp.HandlerFor(
				prometheus.DefaultGatherer,
				promhttp.HandlerOpts{
					ErrorLog: &promHTTPLogger{
						logger: logger,
					},
				},
			),
		),
	)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Consul Exporter</title></head>
             <body>
             <h1>Consul Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             <h2>Options</h2>
             <pre>` + opts.String() + `</pre>
             </dl>
             <h2>Build</h2>
             <pre>` + version.Info() + ` ` + version.BuildContext() + `</pre>
             </body>
             </html>`))
	})

	level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}

}
