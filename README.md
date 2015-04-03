go-latest 
====

[![GitHub release](http://img.shields.io/github/release/tcnksm/go-latest.svg?style=flat-square)][release]
[![Wercker](http://img.shields.io/wercker/ci/551e58c16b7badb9770001288.svg?style=flat-square)][wercker]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]

[release]: https://github.com/tcnksm/go-latest/releases
[wercker]: https://app.wercker.com/project/bykey/1059e8b0cf3bde5fc220477d39a1bf0e
[license]: https://github.com/tcnksm/go-latest/blob/master/LICENSE
[godocs]: http://godoc.org/github.com/tcnksm/go-latest


`go-latest` is a pacakge to check version is latest or not from various sources.

If you're building tool in Golang, you can use this pacakge for encourage user to upgrade latest version of your tool. For source to check, currecntly we can use Tags on Github, [HTML Meta tag](doc/html_meta.md), JSON response and HTML scraping.

See more details in document at [https://godoc.org/github.com/tcnksm/go-latest](https://godoc.org/github.com/tcnksm/go-latest).

## Install

To install, use `go get`:

```bash
$ go get -d github.com/tcnksm/go-latest
```

## Usage

For sources to check, currecntly we can use Tags on Github, [HTML Meta tag](doc/html_meta.md), JSON response and HTML scraping.

### GithubTag

If you want to check [https://github.com/tcnksm/ghr](https://github.com/tcnksm/ghr) version `0.1.0` is latest or not on Github.

```golang
githubTag := &latest.GithubTag{
    Owner: "tcnksm",
    Repository: "ghr"
}

res, _ := latest.Check("0.1.0",githubTag)
if res.latest {
    fmt.Printf("0.1.0 is not latest, you should upgrade to %s", res.Current)
}
```

### HTML Meta tag

You can use a simple HTTP+HTML version check. For example, if you have a tool named `reduce-worker` and want to check `0.1.0` is latest or not. First prepare HTML page which included following meta tag,

```html
<meta name="go-latest" content="reduce-worker 0.1.1 New version include security update">
```

And create following request,

```golang
html := &latest.HTMLMeta{
    URL: "http://example.com/info",
    Name: "reduce-worker",
}

res, _ := latest.Check("0.1.0", html)
if res.latest {
    fmt.Printf("0.1.0 is not latest, %s, upgrade to %s", res.Meta.Message, res.Current)
}
```

To know about HTML Meta tag spec, see [HTML Meta tag](doc/html_meta.md).

And you can prepare your own HTML structure and its scraping fuction. See more details in document at [https://godoc.org/github.com/tcnksm/go-latest](https://godoc.org/github.com/tcnksm/go-latest).

### JSON

If you have an API which return following response, you can use it for checking. For example, API returns following response,

```json
{
    "version":"1.2.3",
    "message":"New version include security update, you should update soon",
    "url":"http://example.com/info"
}
```

To check `0.1.1` is latest, just make following request.

```golang
json := &latest.JSON{
    URL: "http://example.com/json",
}

res, _ := latest.Check("0.1.0", json)
if res.latest {
    fmt.Printf("0.1.0 is not latest, %s, upgrade to %s", res.Meta.Message, res.Current)
}
```

You can prepare your own json receiver. See more details in document at [https://godoc.org/github.com/tcnksm/go-latest](https://godoc.org/github.com/tcnksm/go-latest).

## Version comparing

To compare version, we use [hashicorp/go-version](https://github.com/hashicorp/go-version). `go-version` follows [Semantic Versoning](http://semver.org/). So to use `go-latest` you need to follow SemVer format.

For user who doesn't use SemVer format, `go-latest` has function to transform it into SemVer format.

## Contribution

1. Fork ([https://github.com/tcnksm/go-latest/fork](https://github.com/tcnksm/go-latest/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create new Pull Request

## Author

[Taichi Nakashima](https://github.com/tcnksm)
