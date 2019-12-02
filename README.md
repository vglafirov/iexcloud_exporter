# IEX Cloud Prometheus Exporter

[![Docker Pulls](https://img.shields.io/docker/pulls/vglafirov/iexcloud_exporter.svg?maxAge=604800)](https://hub.docker.com/repository/docker/vglafirov/iexcloud_exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/vglafirov/iexcloud_exporter)](https://goreportcard.com/report/github.com/vglafirov/iexcloud_exporter)
[![GitHub Actions](https://github.com/vglafirov/iexcloud_exporter/workflows/Go/badge.svg)](https://github.com/vglafirov/iexcloud_exporter/actions)

Export Stock data, provided by IEX Cloud to Prometheus

## Supported metric groups
* Price
* Dividents
* Keystats


## Build and run locally:

```bash
make
./iexcloud_exporter [flags]
```

## Flags
|Name|Default|Description|Required|
|---|---|---|---|
|--web.listen-address|:9107|Address to listen on for web interface and telemetry|No|
|---web.telemetry-path|/metrics|Path under which to expose metrics|No|
|---kv.prefix|""|Prefix from which to expose key/value pairs|No|
|---kv.filter|.*|Regex that determines which keys to expose|No|
|---iexcloud.api_token|None|API Token for IEX Cloud account|**Yes**|
|---iexcloud.endpoint|sandbox.iexapis.com|IEX Cloud API endpoint|No|
|--iexcloud.api_version|stable|IEX Cloud API version|No|
|iexcloud.config|`$(pwd)/config.json`|Config path|**Yes**|

## Config

Config file is in `json` format, which consists of `json` array of metrics groups (according to IEX Cloud API specs) with input parameters. Example:
```json
{
  "metrics": [
    {
      "price": {
        "symbols": [
          "aapl"
        ]
      }
    },
    {
      "dividends": {
        "symbols": [
          "aapl"
        ],
        "range": [
          "1y"
        ]
      }
    },
    {
      "keystats": {
        "symbols": [
          "aapl"
        ]
      }
    }
  ]
}
```

## Current stock price

### Parameters
|Parameters|Description|
|---|---|
|symbols|List of symbols|

### Metrics
|Metric|Labels|Description|
|---|---|---|
|iexcloud_price|symbol|Current stock price|

## Dividends for the given stock symbol and the given date range

### Parameters
|Parameters|Description|Example|
|---|---|---|
|symbols|List of symbols|AAPL|
|range|Date range|1y|
*Date ranges format can be found in [API documentation](https://iexcloud.io/docs/api/#dividends-basic)*


### Metrics
|Metric|Description|
|---|---|
|iexcloud_dividends|Dividends|
### Labels
|Label|Format|Description|
|---|---|---|
|declaredDate|Date|Dividend declaration date|
|declaredDate|Date|Dividend ex-date|
|paymentDate|Date|Dividend payment date|
|paymentDate|Date|Dividend payment date|
|range|String|Date range|
|recordDate|Date|Dividend record date|
|symbol|String|Stock symbol|
