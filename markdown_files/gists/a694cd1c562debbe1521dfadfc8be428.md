# [Deleting go package tagged versions from pkg.go.dev] 

[View original Gist on GitHub](https://gist.github.com/Integralist/a694cd1c562debbe1521dfadfc8be428)

**Tags:** #go #golang #pkg #tag #semver

## Deleting go package tagged versions from pkg.go.dev.md

There are various issues open about this:

- https://github.com/golang/go/issues/38848#issuecomment-669945306
- https://github.com/golang/go/issues/37106#issuecomment-690419151
- https://github.com/golang/go/issues/36811#issuecomment-579404726

It seems the process is to get information about your package (e.g. go-fastly) from an endpoint like the following:

https://proxy.golang.org/github.com/fastly/go-fastly/@v/master.info

```
{"Version":"v1.15.1-0.20200611145936-5f794c8e3c37","Time":"2020-06-11T14:59:36Z"}
```

Then use the 'Version' in the following endpoint:

https://pkg.go.dev/github.com/fastly/go-fastly/master@v1.15.1-0.20200611145936-5f794c8e3c37

But that didn't work for me. I then noticed they suggest visiting:

```
pkg.go.dev/<import-path>@<semantic-version>
```

So something like:

https://pkg.go.dev/github.com/fastly/go-fastly/v2@v2.1.0

But that just shows me the existing tagged version, where as we're trying to remove deleted tagged versions from pkg.go.dev.

