# benchito

Like [pgbench](https://www.postgresql.org/docs/current/pgbench.html) or [sysbench](https://github.com/akopytov/sysbench) but only for testing maximum number of connections. `benchito` will start multiple threads to issue very simple queries in order to avoid CPU or memory starvation.

`benchito` supports:
* MySQL
* PostgreSQL

## Requirements

* `make`
* go 1.18

## Setup

Compile the `benchito` binary:

```
make
```

Start database instances:

```
docker-compose pull
docker-compose up -d
```

## Usage

```
./bin/benchito -help
```

## Cleanup

```
docker-compose down -v
```