# latest

`latest` is a command to check a provided version is latest or not in GitHub. 

## Usage

To check cloned repository is latest or not, just run with owner name and repository name which you want to check. If it is not latest version, it returns non-zero exit code.

```go
$ latest -owner=tcnksm -repo=go-latest 2.4.1
$ echo $?
0
```

You can check version is new, it means version is not exist on GitHub and greater than others, and more outputs can be enabled with `-debug` flag, 

```go
$ latest -debug -new -owner=tcnksm repo=go-latest 2.4.1
2.2.1 is new
```

See more usage with `-help` options.


