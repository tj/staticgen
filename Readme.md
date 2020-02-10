# Staticgen

A static website generator that lets you use HTTP servers and frameworks you already know. Just tell Staticgen how to start your server, then watch it crawl your site and generate a static version with all of the pages and assets required.

## About

If you're unfamiliar, you can actually use the decades-old [wget command](https://apex.sh/blog/post/pre-render-wget/) to output a static website from a dynamic one, this project is purpose-built for the same idea, letting your team to use whatever HTTP servers and frameworks you're already familiar with, in any language.

I haven't done any scientific benchmarks or comparisons yet, but here are some results on my 2014 8-core MBP:

- Compiles 3,296 pages of the [Signal v. Noise](https://m.signalvnoise.com/) blog in 1 second
- Compiles my [Apex Software](https://apex.sh/) site in 150ms

## Installation

Via [gobinaries.com](https://gobinaries.com):

```sh
$ curl -sf https://gobinaries.com/tj/staticgen/cmd/staticgen | sh
```

## Configuration

Configuration is stored within a `./static.json` file in your project's root directory. The following options are available:

- __command__ — The server command executed before crawling.
- __url__ — The target website to crawl. Defaults to `"http://127.0.0.1:3000"`.
- __dir__ —  The static website output directory. Defaults to `"build"`.
- __pages__ —  A list of paths added to crawl, typically including unlinked pages such as landing pages. Defaults to `[]`.
- __concurrency__ — The number of concurrent pages to crawl. Defaults to `30`.

## Guide

First create the `./static.json` configuration file, for example here's the config for Go server, the only required property is `command`:

```json
{
  "command": "go run main.go",
  "concurrency": 50,
  "dir": "dist"
}
```

Below is an example of a Node.js server, note that `NODE_ENV` is assigned to production so that optimizations such as Express template caches are used to improve serving performance.

```json
{
  "command": "NODE_ENV=production node server.js"
}
```

Run the `staticgen` command to start the pre-rendering process:

```
$ staticgen
```

Staticgen executes the `command` you provided, waits for the server to become available on the `url` configured. The pages and assets are copied to the `dir` configured and then your server is shut down.

By default the timeout for the generation process is 15 minutes, depending on your situation you may want to increase or decrease this with the `-t, --timeout` flag, here are some examples:

```
$ staticgen -t 30s
$ staticgen -t 15m
$ staticgen -t 1h
```

When launching the `command`, Staticgen sets the `STATICGEN` environment variable to `1`, allowing you to alter behaviour if necessary.

To view the pre-rendered site run the following command to start a static file server and open the browser:

```
$ staticgen serve
```

See the [examples](./_examples) directory for full examples.

## Notes

Staticgen does not pre-render using a headless browser, this makes it faster, however it means that you cannot rely on client-side JavaScript manipulating the page.


---

[![GoDoc](https://godoc.org/github.com/tj/staticgen?status.svg)](https://godoc.org/github.com/tj/staticgen)
![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-stable-green.svg)

## Sponsors

This project is sponsored by my [GitHub sponsors](https://github.com/sponsors/tj) and by [CTO.ai](https://cto.ai/), making it easy for development teams to create and share workflow automations without leaving the command line.

[![](https://apex-software.imgix.net/github/sponsors/cto.png)](https://cto.ai/)
