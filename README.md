# Collector

Project has two components:
* **node** - simulates a remote device generating telemetry data.
* **sink** - receives telemetry from sensor nodes and writes it to a log file.

The easiest way to start both `node` and `sink` is to use provided `docker-compose.yml` file, by running:
```
docker compose up
```

> Output data are written into `out/data.log`.
> Directory is created automatically if not exists.


## Sink

Supported options:
  - *out* - Path to the output file. Destintaion directory must exist. This option is **required**.
  - *listen* - Address to listen on for data requests. Default is ":8080"
  - *buffer* - Size of output buffer in bytes. Default is 1000000.
  - *interval* - How often flush the buffer. Default is 100ms.
  - *rate* - Maximum allowed input flow rate in bytes/sec.  Default is 50000 which corresponds ~1000 messages per second.


## Node

Supported options:
  - *dst* - Address of the sink where to send data. This option is **required**.
  - *name* - Name of the sensor. Default is `node`.
  - *rate* - Number of messages per second to send. Default is 100.
  - *buffer* - Buffer size in bytes to cache generated data. Default is 100000 which should keep messages for last 20 seconds with default rate.
