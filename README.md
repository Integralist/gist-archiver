# Gist Archiver

> \[!IMPORTANT\]
> This code is pretty much ALL generated via AI.\
> This includes Go, JavaScript, CSS and HTML.\
> It was done as part of a "vibe coding" session.\
> I've made _slight_ code tweaks to make it easier to maintain, but not much.\
> The code is still very much an AI `:dumpster-fire:`,\
> But I am still surprised at how well it works,\
> and also how quickly AI got me to a working product.

## View

[https://gist-archiver.netlify.app/](https://gist-archiver.netlify.app/)

or

```shell
open index.html
```

## Build

```shell
make build
```

## CI

I have a [GitHub Actions workflow](.github/workflows/daily-build.yml) that pulls
the latest public gists daily.

## Features

The user interface includes a client-side search feature powered by [Lunr.js](https://lunrjs.com/). It supports:

- **Prefix matching**: Find gists as you type.
- **Word stemming**: Searching for `iterate` will also find gists with `iteration` or `iterating`.

For more details, see the [search guide](./docs/search-guide.md).

## Documentation

Detailed documentation on how the code works can be found in the [./docs/index.md](./docs/index.md) file.

## Gist Struture

Write and structure your gists however you like.

But if you want the "tag" feature that this UI provides then you'll need to
ensure the "Gist description..." field contains twitter-like hash tags.

For example...

```
<description> #<tag1> #<tag2>
```

So a real example might be a gist whose description looks like this:

```
structured logging #go #logs
```
