# Greedy
Greedy allows you to run your own collection of `urls to read`. For example: you are reading a webpage and you want to mark it as `read later`, you can use this service to quickly save the page and read it later. From there, the service generates an RSS feed containing all these urls. You can import this RSS feed into your own, favorite RSS client (f.e. [TinyTinyRSS](https://tt-rss.org "TinyTinyRSS")).

# Storage
Greedy uses a local database file ([boltdb](https://github.com/boltdb/bolt)) as it's storage. You can specify the location of the database file by setting the `databasefile` environment variable. See environment section below for more information.

# Installation instructions
###### I don't have Go installed
If you don't have a working Go environment, then you can simply download one of the pre-built binaries. For download links, see below. After putting the binary in your path, start the service by running the binary: `./greedy`. By default, it binds to 0.0.0.0:8080, but you can change it's configuration by setting some environment variables which are described below. After starting, the output should be like this:

    INFO[0000] environment loaded [host: 0.0.0.0], [port: 8080], [databasefile: articles.bolt]
    INFO[0000] greedy info [builddate: ], [git commit hash: ]
    INFO[0000] bucket initialized with 0 records
    INFO[0000] deamon running on host 0.0.0.0 and port 8080

###### I have Go installed
If you have Go installed, simply `go get` it:

    go get github.com/rogierlommers/greedy

this will download the sources to your `$gopath`, build a binary and puts it in your Go binary directory. You can leave it there or you can put it in a more convenient place. You can manually start a build by running `make build`. Please notice that you will need the Go-Rice tool (https://github.com/GeertJohan/go.rice) to embed the static files to your binary.

###### Usage
Once you have installed and started Greedy, open a browser and point to the `host:port` you have configurated. The greedy homepage should appear. Drag the button to your favorites/bookmarks bar. It is a bookmarklet which redirects to the service and stores the current page to your reading list. Next step is to add the /rss endpoint to your RSS aggregator.

###### Configuration
You can change the default configuration by changing environment vars. For example, running on port 9090 can be done with: `GREEDY_PORT=9090 ./greedy`.

| environment var     | description               | default           |
| --------------------|:-------------------------:| ------------------|
| GREEDY_HOST         | host it binds to          | 0.0.0.0           |
| GREEDY_PORT         | http port                 | 8080              |
| GREEDY_DATABASEFILE | location of database file | ./articles.bolt   |

###### Need help?
For more information, please don't hesitate to contact me [@rogierlommers](https://twitter.com/rogierlommers).

# History
- 1.0
  - uses BOLT as storage engine
  - display articles as RSS
  - scrapes title and page description
  - single binary containing all (static) files, easy to install
  - multiple platforms: linux and darwin

# Screenshots
![home page](./docs/gui-01.png)

![stats page](./docs/gui-02.png)
