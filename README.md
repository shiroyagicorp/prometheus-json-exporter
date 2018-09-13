prometheus-json-exporter
==============================

A prometheus exporter which fetches arbitrary JSON from remote API and
exports the values as gauge metrics.

Installation
--------------------

```
$ go get github.com/shiroyagicorp/prometheus-json-exporter
```

Example Usage
--------------------

```
$ curl -s "http://validate.jsontest.com/?json=%7B%22key%22:%22value%22%7D"
{
   "object_or_array": "object",
   "empty": false,
   "parse_time_nanoseconds": 24618,
   "validate": true,
   "size": 1
}

$ curl -s "http://localhost:9116/probe?target=http://validate.jsontest.com/?json=%7B%22key%22:%22value%22%7D"
# HELP empty Retrieved value
# TYPE empty gauge
empty 0
# HELP parse_time_nanoseconds Retrieved value
# TYPE parse_time_nanoseconds gauge
parse_time_nanoseconds 41626
# HELP size Retrieved value
# TYPE size gauge
size 1
# HELP validate Retrieved value
# TYPE validate gauge
validate 1
```

Note
----------

This repository is a fork of https://github.com/tolleiv/json-exporter

License
----------

MIT License
