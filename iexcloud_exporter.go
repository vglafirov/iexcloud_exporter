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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/vglafirov/iexcloud_exporter/pkg/config"
	"github.com/vglafirov/iexcloud_exporter/pkg/model"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	iex "github.com/vglafirov/iexcloud"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "", "up"),
		"Was the last query of iexcloud successful.",
		nil, nil,
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
	configPath string
}

// Exporter object
type Exporter struct {
	Client   *iex.Client
	kvPrefix string
	kvFilter *regexp.Regexp
	logger   log.Logger
	config   json.RawMessage
}

func (o iexcloudOpts) String() string {
	return fmt.Sprintf("Endpoint: %s\n API version: %s", o.endpoint, o.apiVersion)
}

// Describe describes all the metrics ever exported by the Consul exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- model.PriceMetric
	ch <- model.DividendsMetric
}

// Collect fetches the stats from configured Consul location and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ok := e.collectMetrics(ch)
	level.Info(e.logger).Log("msg", "collecting metrics", "result", ok)

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

func exists(m config.Metric, s string) bool {
	_, ok := m[s]
	return ok
}

func (e *Exporter) collectMetrics(ch chan<- prometheus.Metric) bool {
	var cfg config.Config

	if err := json.Unmarshal(e.config, &cfg); err != nil {
		level.Error(e.logger).Log("msg", "cannot read JSON data", "err", err)
	}
	total := len(cfg.Metrics)
	var wg sync.WaitGroup
	wg.Add(total)

	// Concurently collecting all the metrics
	for count, m := range cfg.Metrics {
		metric := m
		go func(count int) {
			defer wg.Done()
			switch {
			case exists(metric, "price"):
				var price model.Price
				price.Client = e.Client
				price.Symbols = model.SetPriceParams(metric["price"])

				level.Info(e.logger).Log("msg", "collecting Price for", "symbols", len(price.Symbols))

				if err := price.API(ch); err != nil {
					level.Error(e.logger).Log("msg", "cannot collect Price data", "err", err)
				}
			case exists(metric, "dividends"):
				var dividends model.Dividend
				dividends.Client = e.Client
				level.Info(e.logger).Log("msg", "collecting dividents metrics")
				if err := model.SetDividendsParams(&dividends, metric["dividends"]); err != nil {
					level.Error(e.logger).Log("msg", "cannot collect Dividends data", "err", err)
				}
				if err := dividends.API(ch); err != nil {
					level.Error(e.logger).Log("msg", "cannot collect dividends data", "err", err)
				}
			default:
				level.Warn(e.logger).Log("msg", "no metrics configured")
			}
		}(count)
	}

	level.Info(e.logger).Log("msg", "waiting for metrics to be collected", "total", total)

	wg.Wait()

	return true
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

	configFile, err := os.Open(opts.configPath)
	if err != nil {
		return nil, fmt.Errorf("Error opening config file: %s", err)
	}
	level.Info(logger).Log("msg", "Reading the config file", "config", opts.configPath)

	defer configFile.Close()

	config, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, fmt.Errorf("Error reading config file: %s", err)
	}

	level.Info(logger).Log("msg", "initializing endpoint", "endpoint", e)

	client := iex.NewClient(opts.apiToken, e.String())

	// Init our exporter.
	return &Exporter{
		Client:   client,
		kvPrefix: kvPrefix,
		kvFilter: regexp.MustCompile(kvFilter),
		logger:   logger,
		config:   config,
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
	pwd, _ := os.Getwd()
	kingpin.Flag("iexcloud.config", "IEX Cloud API version").Default(pwd + "/config.json").StringVar(&opts.configPath)

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "starting iexcloud_exporter", "version", version.Info())
	level.Info(logger).Log("build_context", version.BuildContext())

	exporter, err := NewExporter(opts, *kvPrefix, *kvFilter, logger)
	if err != nil {
		level.Error(logger).Log("msg", "error creating the exporter", "err", err)
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

	level.Info(logger).Log("msg", "listening on address", "address", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		level.Error(logger).Log("msg", "error starting HTTP server", "err", err)
		os.Exit(1)
	}

}
