Greedy
===========
Greedy allows you to run your own collection of `urls to read`. For example: you are reading a webpage and you want to mark it as `read later`, you can use this service to quickly save the page and read it later. From there, the service generates an RSS feed containing all these urls. You can import this RSS feed into your own, favorite RSS client (f.e. [TinyTinyRSS](https://tt-rss.org "TinyTinyRSS")).

Storage
============
Greedy uses a local database file (sqlite3) as it's storage. You can specify the location of the database file by setting the `databasefile` environment variable. See section below for more information.

Installation instructions
=========================
Simply install by:

    go get github.com/rogierlommers/greedy

this will download the sources to your `$gopath`, build a binary and puts it in your Go binary directory. You can leave it there or you can put it in a more convenient place. I have put it in `/srv/greedy`. Now you can start the service by running the binary: `./greedy`. By default, it binds to 0.0.0.0:8080, but you can change it's configuration by setting some environment variables.

| environment var     | description               | default           |
| --------------------|---------------------------| ------------------|
| GREEDY_HOST         | host it bids to           | 0.0.0.0           |
| GREEDY_PORT         | http port                 | 8080              |
| GREEDY_DATABASEFILE | location of database file | ./articles.sqlite |

After starting, the output should be like this:

    INFO[10-09|21:23:19] environment vars                         host=0.0.0.0 port=8080 databasefile=articles.sqlite
    WARN[10-09|21:23:19] greedy meta info                         builddate=09-October-2015/21:16:25
    DBUG[10-09|21:23:19] check if database file exists            check result=true
    INFO[10-09|21:23:19] number of records in database            amount=666
    INFO[10-09|21:23:19] deamon listening                         host=0.0.0.0 port=8080
    INFO[10-09|21:23:19] cleanup database file                    amount deleted=1234

Now open a browser and point to the `host:port` you have configurated. The greedy homepage should appear. Drag the button to your favorites/bookmarks bar. It is a bookmarklet which redirects to the service and stores the current page to your reading list. Next step is to add the /rss endpoint to your RSS aggregator.

For more information, please don't hesitate to contact me [@rogierlommers](https://twitter.com/rogierlommers).

Releases
=========================
| version     | download    | features    |
| ------------|-------------|-------------|
| 1.0         | greedy1.0   | see main    |

Todo
=======
- [x] use sqlite3, instead of xml file
- [x] make use of makefile
- [x] use godep as dependency management
- [x] use [spf13/viper](https://github.com/spf13/viper) package to read environment vars
- [?] extract serverlocation from header (find out if host starts with www?) 
- [x] use [GeertJohan/go.rice](https://github.com/GeertJohan/go.rice) to add static files to binary
- [x] implement [inconshreveable/log15](https://github.com/inconshreveable/log15) as logger, replacing glog
- [ ] implement github's releases feature and add versions for all platforms
- [ ] some kind of authentication
- [ ] finish cleanup routine
- [ ] create Dockerfile
- [ ] create new screenshots
- [x] update readme

Screenshots
=======
204/232 <--> 847/480

![home page](./docs/gui-01.png)

![stats page](./docs/gui-02.png)
