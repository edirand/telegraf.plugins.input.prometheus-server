# telegraf.plugins.input.prometheus-server
A custom plugin to allow gather data directly from the prometheus API.
This plugin is a work in progress, some improvements can be made : 
- Add real unit tests
- Use concurrency to make the collection faster
- Clean the code

## How to install and run this plugin
To install this plugin, you need GO and to download source code of telegraf with GO (check this documentation : https://medium.com/@punitck05/how-to-write-sample-telegraf-plugin-4b674033df97).

Then in $HOME/go/src/github.com/influxdata/telegraf/plugins/inputs/all, update all.go and add this line :

```go
_ "github.com/influxdata/telegraf/plugins/inputs/prometheusserver"
```

Go back to $HOME/go/src/github.com/influxdata/telegraf/. Use make to build telegraf. Once it's done, you can generate a config file using this command : 

```bash
./telegraf -sample-config -input-filter prometheus_scrapper -output-filter influxdb -debug >> telegraf.conf.test
```

It will configure telegraf to use prometheusserver plugin as input and influxDB as output.
Modify this file to use the influxDB instance you want, and change the parameters Host and Metrics for prometheusserver section : 

```toml
# Collects metrics from prometheus server to export or work on it.
[[inputs.prometheusserver]]
## The list of metrics to collect (you can specify multiple metrics using semicolons)
metrics = """
some_metric{someTag=\"abc\"};
some_other_metrics:rate1m{someTag1=\"abcd\", someTag2=\"efgh\"};
"""

##The prometheus url
host = "http://prometheus"
```

Warn : for the host section, just enter the host. The plugin will add /api/v1/query as path.

Then enter this command to make it run : 
```bash
 ./telegraf -config telegraf.conf.test -debug
```

If the configuration is correct, you should see this kind of messages : 
```bash
2020-07-31T09:44:42Z I! Starting Telegraf
2020-07-31T09:44:42Z I! Loaded inputs: prometheusserver
2020-07-31T09:44:42Z I! Loaded aggregators:
2020-07-31T09:44:42Z I! Loaded processors:
2020-07-31T09:44:42Z I! Loaded outputs: influxdb
2020-07-31T09:44:42Z I! [agent] Config: Interval:10s, Quiet:false, Hostname:"EMILIEN", Flush Interval:10s
2020-07-31T09:44:42Z D! [agent] Initializing plugins
2020-07-31T09:44:42Z D! [agent] Connecting outputs
2020-07-31T09:44:42Z D! [agent] Attempting connection to [outputs.influxdb]
2020-07-31T09:44:42Z D! [agent] Successfully connected to outputs.influxdb
2020-07-31T09:44:42Z D! [agent] Starting service inputs
...
2020-07-31T09:45:12Z D! [outputs.influxdb] Wrote batch of 366 metrics in 37.3236ms
2020-07-31T09:45:12Z D! [outputs.influxdb] Buffer fullness: 0 / 10000 metrics
2020-07-31T09:45:22Z D! [outputs.influxdb] Wrote batch of 121 metrics in 13.4782ms
2020-07-31T09:45:22Z D! [outputs.influxdb] Buffer fullness: 0 / 10000 metrics
2020-07-31T09:45:32Z D! [outputs.influxdb] Wrote batch of 121 metrics in 9.4097ms
2020-07-31T09:45:32Z D! [outputs.influxdb] Buffer fullness: 0 / 10000 metrics
```

Then check in your influxDB and you should see the metrics collected.
