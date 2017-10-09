# Switchman - A simple proxy pass for developers

Switchman is currently an easy to use proxy pass for development environments.
It forwards requests from the `listen` interface to the destination server and allow to `rewrite` the routes.

## Installation

```sh
go get -u github.com/mdouchement/switchman
```


## Requirements

- Golang 1.7 and above


## Usage

```sh
switchman -c my_switchman.yml
```


## Configuration file

```yaml
listen: localhost:4242
rules:
  "/files/*":
    name: Minio
    type: proxy
    url: http://localhost:9000
  "/api/*":
    name: MyAPI
    type: proxy
    url: http://localhost:3000
    rewrite:
      from: "/api(?P<keep>.*)"
      to: "<keep>"
  "/*":
    name: MyFrontendAssets
    type: proxy
    url: http://localhost:8000
```

The key of a rule is the matching pattern of all requests path.

The rewrite `from` field use the Golang Regexp named groups functionality that you can combine as you want in the `to` field.
You can write more complicated example as:

```yaml
rewrite:
  from: "/api(?P<keep>[^\?]*)\?(?P<queries>.*)"
  to: "<keep>/new?<queries>"

#  in: /api/datasets?toto=trololo
# out: /datasets/new?toto=trololo
```


## License

**MIT**


## Contributing

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
5. Push to the branch (git push origin my-new-feature)
6. Create new Pull Request
