# olivine
A simple command to fetch hadoop mapreduce job histories.

## install

`make` or `go get "github.com/plainbanana/olivine"`

## usage

```
$ olivine --hostfile jobhistorynodes > olivine.csv
```

## help

```
$ olivine --help
A command to fetch hadoop job histories.

Usage:
  olivine [flags]

Flags:
  -h, --help              help for olivine
      --hostfile string   Specify target hosts from a hostfile. default target is localhost.
  -p, --port string       Specify the port where target hadoop job history server running on hosts. (default "19888")
```