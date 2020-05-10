<p align="center">
<br>
<img src="https://github.com/k1LoW/tbls-meta/raw/master/img/logo.png" width="200" alt="tbls-meta">
<br>
<br>
</p>

`tbls-meta` is an external subcommand of tbls for applying metadata managed by [tbls](https://github.com/k1LoW/tbls) to the datasource.

## Usage

tbls-meta is provided as an external subcommand of [tbls](https://github.com/k1LoW/tbls).

``` console
$ tbls meta plan -c /path/to/tbls.yml
```

``` console
$ tbls meta apply -c /path/to/tbls.yml
```

## Install

**deb:**

Use [dpkg-i-from-url](https://github.com/k1LoW/dpkg-i-from-url)

``` console
$ export TBLS_META_VERSION=X.X.X
$ curl -L https://git.io/dpkg-i-from-url | bash -s -- https://github.com/k1LoW/tbls-meta/releases/download/v$TBLS_META_VERSION/tbls-meta_$TBLS_META_VERSION-1_amd64.deb
```

**RPM:**

``` console
$ export TBLS_META_VERSION=X.X.X
$ yum install https://github.com/k1LoW/tbls-meta/releases/download/v$TBLS_META_VERSION/tbls-meta_$TBLS_META_VERSION-1_amd64.rpm
```

**homebrew tap:**

```console
$ brew install k1LoW/tap/tbls-meta
```

**manually:**

Download binary from [releases page](https://github.com/k1LoW/tbls-meta/releases)

**go get:**

```console
$ go get github.com/k1LoW/tbls-meta
```

## Requirements

- [tbls](https://github.com/k1LoW/tbls) > 1.35.0

## Support Datasource

**BigQuery:**

Required permissions: `bigquery.datasets.get` `bigquery.datasets.update` `bigquery.tables.get` `bigquery.tables.update` `bigquery.tables.list`
