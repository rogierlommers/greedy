Greedy
===========
Greedy allows you to run your own collection of `urls to read`. For example: you are reading a webpage and you want to mark it as `read later`, you can use this service to quickly save the page and read it later. From there, the service generates an RSS feed containing all these urls. You can import this RSS feed into your own, favorite RSS client (f.e. [TinyTinyRSS](https://tt-rss.org "TinyTinyRSS")).

Storage
============
Greedy uses a local database file (sqlite3) as it's storage. You can specify the location of the database file by setting the `databasefile` environment variable. See section below for more information.

Installation instructions
=========================
If you don't have a working Go environment, then you can simply download one of the pre-built binaries. For download links, see below. After putting the binary in your path, start the service by running the binary: `./greedy`. By default, it binds to 0.0.0.0:8080, but you can change it's configuration by setting some environment variables which are described below. After starting, the output should be like this:

    INFO[10-18|15:48:03] environment vars                         host=0.0.0.0 port=8080 databasefile=articles.bolt
    INFO[10-18|15:48:03] greedy meta info                         builddate=18-October-2015/15:47:55 commithash=9cb5fe13067f2dead95233e36b3b7b9fd1dd2b73
    INFO[10-18|15:48:03] deamon listening                         host=0.0.0.0 port=8080

Now open a browser and point to the `host:port` you have configurated. The greedy homepage should appear. Drag the button to your favorites/bookmarks bar. It is a bookmarklet which redirects to the service and stores the current page to your reading list. Next step is to add the /rss endpoint to your RSS aggregator.

For more information, please don't hesitate to contact me [@rogierlommers](https://twitter.com/rogierlommers).

If you have Go installed, simply `go get` it:

    go get github.com/rogierlommers/greedy

this will download the sources to your `$gopath`, build a binary and puts it in your Go binary directory. You can leave it there or you can put it in a more convenient place.

| environment var     | description               | default           |
| --------------------|:-------------------------:| ------------------|
| GREEDY_HOST         | host it binds to          | 0.0.0.0           |
| GREEDY_PORT         | http port                 | 8080              |
| GREEDY_DATABASEFILE | location of database file | ./articles.sqlite |

Running in Docker container
===========================
Explained here

Releases
=========================
| version           | download                                                                                                                         |
| ------------------|----------------------------------------------------------------------------------------------------------------------------------|
| 1.0-linux-amd64   | [greedy-1.0-linux-amd64.tar.bz2](https://github.com/rogierlommers/greedy/releases/download/1.0/greedy-1.0-linux-amd64.tar.bz2)   |
| 1.0-linux-368     | [greedy-1.0-linux-368.tar.bz2](https://github.com/rogierlommers/greedy/releases/download/1.0/greedy-1.0-linux-386.tar.bz2)       |
| 1.0-darwin-amd64  | [greedy-1.0-darwin-amd64.tar.bz2](https://github.com/rogierlommers/greedy/releases/download/1.0/greedy-1.0-darwin-amd64.tar.bz2) |
| 1.0-darwin-368    | [greedy-1.0-darwin-386.tar.bz2](https://github.com/rogierlommers/greedy/releases/download/1.0/greedy-1.0-darwin-386.tar.bz2)     |

History
=======
- 1.0
  - uses BOLT as storage engine
  - display articles as RSS
  - scrapes title and page description
  - single binary containing all (static) files, easy to install
  - multiple platforms: linux and darwin

Todo
=======
- [x] use sqlite3, instead of xml file
- [x] make use of makefile
- [x] use godep as dependency management
- [x] use [spf13/viper](https://github.com/spf13/viper) package to read environment vars
- [x] use [GeertJohan/go.rice](https://github.com/GeertJohan/go.rice) to add static files to binary
- [x] implement [inconshreveable/log15](https://github.com/inconshreveable/log15) as logger, replacing glog
- [x] update readme
- [x] fix injection of build date: https://ariejan.net/2015/10/12/building-golang-cli-tools-update/
- [x] implement native, embeddable database (https://www.reddit.com/r/golang/comments/3m1xcu/embeddable_database_for_go/)
- [x] automatic releases --> https://github.com/miekg/mmark/blob/master/.rel.sh
- [x] extract serverlocation from header
- [x] update installation instructions --> add binary section
- [ ] create Dockerfile
- [ ] some kind of authentication?
- [ ] finish cleanup routine
- [ ] create new screenshots
- [ ] add to [avelino/awesome-go](https://github.com/avelino/awesome-go)

Screenshots
=======
204/232 <--> 847/480

![home page](./docs/gui-01.png)

![stats page](./docs/gui-02.png)
