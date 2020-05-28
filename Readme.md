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

This project is sponsored by [CTO.ai](https://cto.ai/), making it easy for development teams to create and share workflow automations without leaving the command line.

[![](https://apex-software.imgix.net/github/sponsors/cto.png)](https://cto.ai/)

And by my [GitHub sponsors](https://github.com/sponsors/tj):

[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/0" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/0)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/1" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/1)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/2" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/2)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/3" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/3)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/4" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/4)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/5" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/5)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/6" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/6)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/7" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/7)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/8" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/8)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/9" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/9)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/10" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/10)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/11" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/11)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/12" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/12)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/13" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/13)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/14" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/14)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/15" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/15)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/16" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/16)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/17" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/17)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/18" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/18)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/19" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/19)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/20" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/20)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/21" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/21)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/22" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/22)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/23" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/23)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/24" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/24)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/25" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/25)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/26" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/26)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/27" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/27)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/28" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/28)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/29" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/29)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/30" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/30)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/31" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/31)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/32" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/32)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/33" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/33)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/34" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/34)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/35" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/35)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/36" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/36)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/37" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/37)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/38" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/38)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/39" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/39)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/40" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/40)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/41" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/41)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/42" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/42)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/43" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/43)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/44" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/44)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/45" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/45)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/46" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/46)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/47" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/47)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/48" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/48)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/49" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/49)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/50" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/50)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/51" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/51)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/52" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/52)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/53" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/53)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/54" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/54)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/55" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/55)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/56" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/56)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/57" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/57)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/58" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/58)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/59" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/59)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/60" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/60)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/61" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/61)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/62" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/62)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/63" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/63)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/64" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/64)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/65" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/65)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/66" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/66)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/67" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/67)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/68" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/68)

