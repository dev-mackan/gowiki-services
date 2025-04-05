

This repository contains code for a API and a frontend webserver, for a
service to store and view markdown documents.

The webserver simply works as a frontend to navigate and view markdown documents
as webpages. The API works as the backend.

The service is meant to become a wiki-like website where you can view pages of
information. Each page can be provided with a markdown file, which is shown as a
webpage when that particular page is visited. Each page keeps revisions of the
markdown file when a new file is uploaded.

## Building and running

A Makefile is provided:

```
$ make all
```

then run each service in separate terminals/TTYs:

```
$ ./bin/api
```
and
```
$ ./bin/web
```


Visit the locally hosted site `localhost:3001/`


