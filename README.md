# tbls-meta

`tbls-meta` is an external subcommand of tbls for applying metadata managed by [tbls](https://github.com/k1LoW/tbls) to the datasource.

## Usage

tbls-meta is provided as an external subcommand of [tbls](https://github.com/k1LoW/tbls).

``` console
$ tbls meta plan -c /path/to/tbls.yml
```

``` console
$ tbls meta apply -c /path/to/tbls.yml
```

## Requirements

- [tbls](https://github.com/k1LoW/tbls) > 1.35.0

## Support Datasource

**BigQuery:**

Required permissions: `bigquery.datasets.get` `bigquery.datasets.update` `bigquery.tables.get` `bigquery.tables.update` `bigquery.tables.list`
