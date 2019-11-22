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
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	// TODO: Replace with github.com/goinvest/iexcloud once https://github.com/goinvest/iexcloud/issues/41 is closed
	iex "github.com/vglafirov/iexcloud"
	"github.com/vglafirov/iexcloud_exporter/pkg/config"
)

var (
	// MarketcapStatsMetric Prometheus metric definition for Market cap.
	MarketcapStatsMetric = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "marketcap"),
		"Market cap of the security calculated as shares outstanding * previous day close.",
		[]string{
			"symbol",
		},
		nil,
	)

	// Week52High 52 weeks high
	Week52High = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "week52high"),
		"52 weeks high",
		[]string{
			"symbol",
		},
		nil,
	)

	// Week52Low 52 weeks high
	Week52Low = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "week52low"),
		"52 weeks low",
		[]string{
			"symbol",
		},
		nil,
	)

	// Week52Change Percentage change
	Week52Change = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "week52change"),
		"Percentage change",
		[]string{
			"symbol",
		},
		nil,
	)

	// SharesOutstanding Number of shares outstanding as the difference between issued shares and treasury shares
	SharesOutstanding = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "sharesOutstanding"),
		"Number of shares outstanding as the difference between issued shares and treasury shares",
		[]string{
			"symbol",
		},
		nil,
	)

	// Avg30Volume Average 30 day volume
	Avg30Volume = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "avg30Volume"),
		"Average 30 day volume",
		[]string{
			"symbol",
		},
		nil,
	)

	// Avg10Volume Average 10 day volume
	Avg10Volume = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "avg10Volume"),
		"Average 10 day volume",
		[]string{
			"symbol",
		},
		nil,
	)

	// Float Returns the annual shares outstanding minus closely held shares.
	Float = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "float"),
		"Returns the annual shares outstanding minus closely held shares.",
		[]string{
			"symbol",
		},
		nil,
	)

	// Employees number of employees
	Employees = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "employees"),
		"Returns the annual shares outstanding minus closely held shares.",
		[]string{
			"symbol",
		},
		nil,
	)

	// TTMEPS Trailing twelve month earnings per share
	TTMEPS = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "ttmEPS"),
		"Trailing twelve month earnings per share",
		[]string{
			"symbol",
		},
		nil,
	)

	// TTMDividendRate Trailing twelve month dividend rate per share
	TTMDividendRate = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "ttmDividendRate"),
		"Trailing twelve month dividend rate per share",
		[]string{
			"symbol",
		},
		nil,
	)

	// DividendYield The ratio of trailing twelve month dividend compared to the previous day close price.
	DividendYield = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "dividendYield"),
		"The ratio of trailing twelve month dividend compared to the previous day close price",
		[]string{
			"symbol",
		},
		nil,
	)

	// PERatio Price to earnings ratio calculated as (previous day close price) / (ttmEPS)
	PERatio = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "peRatio"),
		"Price to earnings ratio calculated as (previous day close price) / (ttmEPS)",
		[]string{
			"symbol",
		},
		nil,
	)

	// Beta Beta is a measure used in fundamental analysis to determine the volatility of an asset or portfolio in relation to the overall market
	Beta = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "beta"),
		"Beta is a measure used in fundamental analysis to determine the volatility of an asset or portfolio in relation to the overall market",
		[]string{
			"symbol",
		},
		nil,
	)

	// Day200MovingAvg 200 days moving average
	Day200MovingAvg = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "day200MovingAvg"),
		"200 days moving average",
		[]string{
			"symbol",
		},
		nil,
	)

	// Day50MovingAvg 50 days moving average
	Day50MovingAvg = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "day50MovingAvg"),
		"50 days moving average",
		[]string{
			"symbol",
		},
		nil,
	)

	// MaxChangePercent Percent change MAX
	MaxChangePercent = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "maxChangePercent"),
		"Percent change MAX",
		[]string{
			"symbol",
		},
		nil,
	)

	// Year5ChangePercent Percent change 5 years
	Year5ChangePercent = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "year5ChangePercent"),
		"Percent change 5 years",
		[]string{
			"symbol",
		},
		nil,
	)

	// Year2ChangePercent Percent change 2 years
	Year2ChangePercent = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "year2ChangePercent"),
		"Percent change 2 years",
		[]string{
			"symbol",
		},
		nil,
	)

	// Year1ChangePercent Percent change 1 year
	Year1ChangePercent = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "year1ChangePercent"),
		"Percent change 1 year",
		[]string{
			"symbol",
		},
		nil,
	)

	// YTDChangePercent Percent change year to date
	YTDChangePercent = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "ytdChangePercent"),
		"Percent change YTD",
		[]string{
			"symbol",
		},
		nil,
	)

	// Month6ChangePercent Percent change 6 months
	Month6ChangePercent = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "month6ChangePercent"),
		"Percent change 6 months",
		[]string{
			"symbol",
		},
		nil,
	)

	// Month3ChangePercent Percent change 3 months
	Month3ChangePercent = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "month3ChangePercent"),
		"Percent change 3 months",
		[]string{
			"symbol",
		},
		nil,
	)

	// Month1ChangePercent Percent change 1 month
	Month1ChangePercent = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "month1ChangePercent"),
		"Percent change 1 month",
		[]string{
			"symbol",
		},
		nil,
	)

	// Day30ChangePercent Percent change 30 days
	Day30ChangePercent = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "day30ChangePercent"),
		"Percent change 30 days",
		[]string{
			"symbol",
		},
		nil,
	)

	// Day5ChangePercent Percent change 5 days
	Day5ChangePercent = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "day5ChangePercent"),
		"Percent change 5 days",
		[]string{
			"symbol",
		},
		nil,
	)

	// KeyStatDates Expected ex date of the next dividend, Ex date of the last dividend, Expected next earnings report date
	KeyStatDates = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "keystats", "dates"),
		"Expected ex date of the next dividend, Ex date of the last dividend, Expected next earnings report date",
		[]string{
			"symbol",
			"nextDividendDate",
			"exDividendDate",
			"nextEarningsDate",
		},
		nil,
	)
)

// KeyStats data
type KeyStats struct {
	Client   *iex.Client
	Symbols  []string
	KeyStats iex.KeyStats
}

func toUnknown(s string) string {
	if s == "1929-10-24" {
		return "unknown"
	}
	return s
}

// API Dividend API call
func (s *KeyStats) API(ch chan<- prometheus.Metric) error {
	for _, symbol := range s.Symbols {
		var err error
		s.KeyStats, err = s.Client.KeyStats(symbol)
		if err != nil {
			return err

		}
		ch <- prometheus.MustNewConstMetric(
			MarketcapStatsMetric, prometheus.GaugeValue, s.KeyStats.MarketCap, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			Week52High, prometheus.GaugeValue, s.KeyStats.Week52High, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			Week52Low, prometheus.GaugeValue, s.KeyStats.Week52Low, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			Week52Change, prometheus.GaugeValue, s.KeyStats.Week52Change, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			SharesOutstanding, prometheus.GaugeValue, s.KeyStats.SharesOutstanding, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			Avg30Volume, prometheus.GaugeValue, s.KeyStats.Avg30Volume, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			Avg10Volume, prometheus.GaugeValue, s.KeyStats.Avg10Volume, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			Float, prometheus.GaugeValue, s.KeyStats.Float, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			Employees, prometheus.GaugeValue, float64(s.KeyStats.Employees), symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			TTMEPS, prometheus.GaugeValue, s.KeyStats.TTMEPS, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			TTMDividendRate, prometheus.GaugeValue, s.KeyStats.TTMDividendRate, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			DividendYield, prometheus.GaugeValue, s.KeyStats.DividendYield, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			PERatio, prometheus.GaugeValue, s.KeyStats.PERatio, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			Beta, prometheus.GaugeValue, s.KeyStats.Beta, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			Day200MovingAvg, prometheus.GaugeValue, s.KeyStats.Day200MovingAvg, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			Day50MovingAvg, prometheus.GaugeValue, s.KeyStats.Day50MovingAvg, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			MaxChangePercent, prometheus.GaugeValue, s.KeyStats.MaxChangePercent, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			Year5ChangePercent, prometheus.GaugeValue, s.KeyStats.Year5ChangePercent, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			Year2ChangePercent, prometheus.GaugeValue, s.KeyStats.Year2ChangePercent, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			Year1ChangePercent, prometheus.GaugeValue, s.KeyStats.Year1ChangePercent, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			YTDChangePercent, prometheus.GaugeValue, s.KeyStats.YTDChangePercent, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			Month6ChangePercent, prometheus.GaugeValue, s.KeyStats.Month6ChangePercent, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			Month3ChangePercent, prometheus.GaugeValue, s.KeyStats.Month3ChangePercent, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			Month1ChangePercent, prometheus.GaugeValue, s.KeyStats.Month1ChangePercent, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			Day30ChangePercent, prometheus.GaugeValue, s.KeyStats.Day30ChangePercent, symbol,
		)
		ch <- prometheus.MustNewConstMetric(
			Day5ChangePercent, prometheus.GaugeValue, s.KeyStats.Day5ChangePercent, symbol,
		)
		nextDividendDate, err := s.KeyStats.NextDividendDate.MarshalJSON()
		if err != nil {
			return err
		}
		exDividendDate, err := s.KeyStats.ExDividendDate.MarshalJSON()
		if err != nil {
			return err
		}
		nextEarningsDate, err := s.KeyStats.NextEarningsDate.MarshalJSON()
		if err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(
			KeyStatDates, prometheus.GaugeValue, 1, symbol,
			toUnknown(strings.Trim(string(nextDividendDate), `"`)),
			toUnknown(strings.Trim(string(exDividendDate), `"`)),
			toUnknown(strings.Trim(string(nextEarningsDate), `"`)),
		)

	}
	return nil
}

// SetKeyStatsParams Converts map of unknown parameters to symbols
func SetKeyStatsParams(stats *KeyStats, p interface{}) error {
	params := p.(map[string]interface{})
	s := params["symbols"].([]interface{})
	for _, symbol := range s {
		stats.Symbols = append(stats.Symbols, symbol.(string))
	}
	return nil
}
