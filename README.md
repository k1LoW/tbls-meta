# tbls-meta

`tbls-meta` is a CI-friendly tool for applying metadata managed by [tbls](https://github.com/k1LoW/tbls) to the datasource.

## Usage

tbls-meta is provided as an external subcommand of [tbls](https://github.com/k1LoW/tbls).

``` console
$ tbls meta plan -c /path/to/tbls.yml
```

``` console
$ tbls meta apply -c /path/to/tbls.yml
```

## Requirements

- [tbls](https://github.com/k1LoW/tbls) > 1.38.2

## Support Datasource

**BigQuery:**

Required permissions: `bigquery.datasets.get` `bigquery.datasets.update` `bigquery.tables.get` `bigquery.tables.update` `bigquery.tables.list`
