
# Node

Setup:

```
$ npm install
```

To compile the site, run `staticgen` in the project's directory, for example:

```
$ cd _examples/node
$ staticgen
$ ls build
```

Note that `NODE_ENV=production` is used in the static.json configuration, this improves the performance with template compilation caching.