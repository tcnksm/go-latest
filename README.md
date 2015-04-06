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


`go-latest` is a pacakge to check a provided version is latest or not from various sources.

Once you distribute your tool and user started to use it, it's difficult to tell users that new version is released and encourage them to use new one. `go-latest` enables you to do that by just preparing simple source. For sources, currecntly you can use tags on Github, [HTML meta tag](doc/html_meta.md) (HTML scraping) and JSON response. 

See more details in document at [https://godoc.org/github.com/tcnksm/go-latest](https://godoc.org/github.com/tcnksm/go-latest).

## Install

To install, use `go get`:

```bash
$ go get -d github.com/tcnksm/go-latest
```

## Usage

For sources to check, currecntly you can use tags on Github, [HTML meta tag](doc/html_meta.md) (HTML scraping) and JSON response. 

### Github Tag

To check `0.1.0` is the latest in tags on GitHub.

```golang
githubTag := &latest.GithubTag{
    Owner: "username",
    Repository: "reponame",
}

res, _ := latest.Check(githubTag,"0.1.0")
if res.Outdated {
    fmt.Printf("0.1.0 is not latest, you should upgrade to %s", res.Current)
}
```

`go-latest` uses [Semantic Versoning](http://semver.org/) to compare versions. If tagging name strategy on GitHub may be different from it, you need to fix it by defining `FixVersionStrFunc`. For example, if you put `v` charactor on beginning of tag name (`v0.1.1`), by default `go-latest` can not interpret it. So you need to use `DeleteFrontV()`,

```golang
githubTag := &latest.GithubTag{
    Owner: "username",
    Repository: "reponame",
    FixVersionStrFunc: latest.DeleteFrontV(),
}
```

### HTML meta tag

To check version, you can use simple HTTP+HTML meta tag.

For example, if you have a tool named `reduce-worker` and want to check `0.1.0` is latest or not, prepare HTML page which includes following meta tag, 

```html
<meta name="go-latest" content="reduce-worker 0.1.1 New version include security update">
```

And make request,

```golang
html := &latest.HTMLMeta{
    URL: "http://example.com/info",
    Name: "reduce-worker",
}

res, _ := latest.Check(html, "0.1.0")
if res.Outdated {
    fmt.Printf("0.1.0 is not latest, %s, upgrade to %s", res.Meta.Message, res.Current)
}
```

To know about HTML meta tag specification, see [HTML Meta tag](doc/html_meta.md).

You can prepare your own HTML page definition and its scraping fuction. See more details in document at [https://godoc.org/github.com/tcnksm/go-latest](https://godoc.org/github.com/tcnksm/go-latest).

### JSON

You can also use JSON response.

If you want to check `0.1.0` is latest or not, prepare an API server which returns a following response,

```json
{
    "version":"1.2.3",
    "message":"New version include security update, you should update soon",
    "url":"http://example.com/info"
}
```

And make request,

```golang
json := &latest.JSON{
    URL: "http://example.com/json",
}

res, _ := latest.Check("0.1.0", json)
if res.Outdated {
    fmt.Printf("0.1.0 is not latest, %s, upgrade to %s", res.Meta.Message, res.Current)
}
```

You can use your own json schema by defining `JSONReceiver` interface. See more details in document at [https://godoc.org/github.com/tcnksm/go-latest](https://godoc.org/github.com/tcnksm/go-latest).

## Version comparing

To compare version, we use [hashicorp/go-version](https://github.com/hashicorp/go-version). `go-version` follows [Semantic Versoning](http://semver.org/). So to use `go-latest` you need to follow SemVer format.

For user who doesn't use SemVer format, `go-latest` has function to transform it into SemVer format.

## Example

Here is tools which is used `go-latest`.

- [tcnksm/ghr](https://github.com/tcnksm/ghr)


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
