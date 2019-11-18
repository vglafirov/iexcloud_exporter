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
	"github.com/prometheus/client_golang/prometheus"
	// TODO: Replace with github.com/goinvest/iexcloud once https://github.com/goinvest/iexcloud/issues/41 is closed
	iex "github.com/vglafirov/iexcloud"
	"github.com/vglafirov/iexcloud_exporter/pkg/config"
)

var (
	// PriceMetric Prometheus metric definition
	PriceMetric = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "", "price"),
		"Current stock price",
		[]string{"symbol"},
		nil,
	)
)

// Price data
type Price struct {
	Client  *iex.Client
	Symbols []string
	Price   float64
}

// API Price API call
func (p *Price) API(ch chan<- prometheus.Metric) error {
	for _, symbol := range p.Symbols {
		var err error
		p.Price, err = p.Client.Price(symbol)
		if err != nil {
			return err

		}
		ch <- prometheus.MustNewConstMetric(
			PriceMetric, prometheus.GaugeValue, float64(p.Price), symbol,
		)
	}
	return nil
}

// SetPriceParams Converts map of unknown parameters to symbols
func SetPriceParams(p interface{}) []string {
	params := p.(map[string]interface{})
	s := params["symbols"].([]interface{})
	var symbols []string
	for _, symbol := range s {
		symbols = append(symbols, symbol.(string))
	}
	return symbols
}
