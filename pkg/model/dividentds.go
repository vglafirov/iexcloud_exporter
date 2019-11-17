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

package model

import (
	"fmt"
	iex "github.com/vglafirov/iexcloud"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vglafirov/iexcloud_exporter/pkg/config"
)

var (
	// DividendsMetric Prometheus metric definition for dividends
	DividendsMetric = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "", "dividends"),
		"dividends from the IEX Cloud endpoint for the given stock symbol and the given date range",
		[]string{"symbol", "range"},
		nil,
	)
)

// Dividend data
type Dividend struct {
	Client    *iex.Client
	Symbols   []string
	Range     []iex.PathRange
	Dividends []iex.Dividend
}

// API Dividend API call
func (d *Dividend) API(ch chan<- prometheus.Metric) error {
	for _, symbol := range d.Symbols {
		for _, pathRange := range d.Range {
			div, err := d.Client.Dividends(symbol, pathRange)
			if err != nil {
				return err
			}
			d.Dividends = div
			for _, dividend := range d.Dividends {
				ch <- prometheus.MustNewConstMetric(
					DividendsMetric,
					prometheus.GaugeValue,
					dividend.Amount,
					symbol,
					string(pathRange),
					fmt.Sprintf("%v", dividend.ExDate),
				)
			}
		}
	}
	return nil
}

// SetDividendsParams Converts map of unknown parameters to symbols
func SetDividendsParams(d *Dividend, p interface{}) error {
	params := p.(map[string]interface{})
	s := params["symbols"].([]interface{})
	r := params["range"].([]interface{})
	for _, symbol := range s {
		d.Symbols = append(d.Symbols, symbol.(string))
	}
	d.Range = make([]iex.PathRange, len(r))
	for i, dividendRange := range r {
		if err := d.Range[i].Set(dividendRange.(string)); err != nil {
			return err
		}
	}
	return nil
}
