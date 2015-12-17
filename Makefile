default: run

VERSION := 1.0
LDFLAGS := -X github.com/rogierlommers/greedy/internal/common.CommitHash=`git rev-parse HEAD` -X github.com/rogierlommers/greedy/internal/common.BuildDate=`date +"%d-%B-%Y/%T"`
BINARY := ./bin/greedy-${VERSION}

build:
	rm -rf ./target
	mkdir -p ./target
	rice embed-go -i ./internal/render/
	CGO_ENABLED=0 go build -ldflags "-s $(LDFLAGS)" -a -installsuffix cgo -o ./target/greedy main.go

run:
	go run *.go

release:
	CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -ldflags "-s $(LDFLAGS)" -a -installsuffix cgo -o $(BINARY)-darwin-386 main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -a -installsuffix cgo -o $(BINARY)-darwin-amd64 main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "$(LDFLAGS)" -a -installsuffix cgo -o $(BINARY)-linux-386 main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -a -installsuffix cgo -o $(BINARY)-linux-amd64 main.go
	zip -m -9 $(BINARY)-darwin-386.zip $(BINARY)-darwin-386
	zip -m -9 $(BINARY)-darwin-amd64.zip $(BINARY)-darwin-amd64
	zip -m -9 $(BINARY)-linux-386.zip $(BINARY)-linux-386
	zip -m -9 $(BINARY)-linux-amd64.zip $(BINARY)-linux-amd64

test:
	go run ./systemtest/main.go
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=http%3A%2F%2Fubuntuforums.org%2Fforumdisplay.php%3Ff%3D339
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=http%3A%2F%2Fwww.iculture.nl%2Ficulture-vandaag-26-mei-2015%2F
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2Fonedirection
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2FHarry_Styles
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2Fjimmyfallon
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2FPink
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2Fcnnbrk
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2Fddlovato
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2FArianaGrande
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2Fselenagomez
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2FKimKardashian
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2FCristiano
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2Ftwitter
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2FTheEllenShow
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2Frihanna
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2Fladygaga
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2FYouTube
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2Ftaylorswift13
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2FBarackObama
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2Fjustinbieber
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2Fkatyperry
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Fwww.wikipedia.de%2F
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Ftwitter.com%2Fwikipedia
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Fen.wikipedia.org%2Fwiki%2FMain_Page
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Fwww.wikipedia.org%2F
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=https%3A%2F%2Fnl.wikipedia.org%2Fwiki%2FHoofdpagina
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=http%3A%2F%2Fstackoverflow.com%2Fquestions%2F8350609%2Fhow-do-you-time-a-function-in-go-and-return-its-runtime-in-milliseconds
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=http%3A%2F%2Ftweakers.net%2Fdownloads%2F34843%2Fvirtualbox-50-beta-4.html
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=http%3A%2F%2Ftweakers.net%2Fgeek%2F103305%2Fgoogle-maakt-roboto-font-helemaal-opensource.html
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=http%3A%2F%2Ftweakers.net%2Fnieuws%2F103304%2Fgerucht-apple-gaat-force-touch-actief-inzetten-bij-nieuwe-iphone.html
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=http%3A%2F%2Ftweakers.net%2Fnieuws%2F103302%2Fsamsung-patentaanvraag-toont-smartphone-met-laptopdock.html
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=http%3A%2F%2Ftweakers.net%2Fnieuws%2F103300%2Fgemeenschap-brengt-final-van-fedora-22-uit.html
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=http%3A%2F%2Ftweakers.net%2Fnieuws%2F103301%2Fopera-laat-android-app-nu-ook-dataverbruik-bij-wifi-terugdringen.html
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=http%3A%2F%2Ftweakers.net%2Fnieuws%2F103299%2Ftwitter-brengt-periscope-app-voor-android-uit.html
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=http%3A%2F%2Ftweakers.net%2Fnieuws%2F103298%2Fbe-quiet-introduceert-dark-rock-top-flow-koeler.html
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=http%3A%2F%2Ftweakers.net%2Fvideo%2F10359%2Fautos-maken-de-man-in-mad-max.html
	curl http://localhost:8080/add?url=http%3A%2F%2Fen.wikipedia.org%2Fwiki%2FWikipedia%3ALists_of_protected_pages http://localhost:8080/add?url=http%3A%2F%2Ftweakers.net%2Fgeek%2F103297%2Fhardwarefan-zet-55-jaar-oude-ibm-mainframe-in-voor-minen-van-bitcoins.html
	curl http://localhost:8080/add?url=http://www.bol.com/nl/l/sport-vrije-tijd/sport-hardlopen-hardloopschoenen-dames/N/4283445807+17553+17257+17322/index.html
	curl http://localhost:8080/add?url=https://github.com/koding/tunnel
	curl http://localhost:8080/add?url=http://esnooker.pl/turnieje/2015/me/en/mem_2015.php
	curl http://localhost:8080/add?url=https://events.google.com/io2015/#
	curl http://localhost:8080/add?url=https://www.reddit.com/r/golang/comments/37jwo9/mongolar_cms/
	curl http://localhost:8080/add?url=https://www.reddit.com/r/golang/comments/37jvaz/go_web_page_scraper/
	curl http://localhost:8080/add?url=https://www.socketloop.com/tutorials/golang-capture-stdout-of-a-child-process-and-act-according-to-the-result
	curl http://localhost:8080/add?url=https://news.ycombinator.com/item?id=5742470
	curl http://localhost:8080/add?url=http://stackoverflow.com/questions/10042388/program-that-convert-html-to-image
	curl http://localhost:8080/add?url=http://www.paulhammond.org/webkit2png/
	curl http://localhost:8080/add?url=https://www.google.nl/search?client=safari&rls=en&q=linux+html+to+image&ie=UTF-8&oe=UTF-8&gfe_rd=cr&ei=w8NlVfj2IsvD0wXD_oDADA
	curl http://localhost:8080/add?url=http://stackoverflow.com/questions/10042388/program-that-convert-html-to-image
	curl http://localhost:8080/add?url=http://cutycapt.sourceforge.net/
	curl http://localhost:8080/add?url=http://stackoverflow.com/questions/3318490/generate-image-e-g-jpg-of-a-web-page
	curl http://localhost:8080/add?url=http://www.packal.org/workflow-list
	curl http://localhost:8080/add?url=http://blog.doit.st/
	curl http://localhost:8080/add?url=https://omgwtfnzbs.org/
	curl http://localhost:8080/add?url=http://www.bol.com/nl/p/zwarte-premium-aluminium-card-case-aluminium-portemonee-hard-case-creditkaarthouder-card-vision-huismerk/9200000010826286/?bltg=itm_event%3dclick%26cnt_disp%3dMHP%26mcm_ccd%3dLHOT%26mcm_refdate%3d20140507%26mcm_cstep%3d1%26slt_subch%3dhs1%26pg_nm%3dmain%26slt_id%3d802%26slt_nm%3dHorizontaal-MCM-slot-hs1%26slt_pos%3dB2%26slt_owner%3dmcm%26itm_type%3dproduct%26itm_lp%3d3%26itm_id%3d9200000010826286%26itm_role%3din&promo=main_802_LHOT-Horizontaal-MCM-slot-hs1_B2_product_3_9200000010826286
	curl http://localhost:8080/add?url=http://talks.golang.org/2015/gogo.slide#1
	curl http://localhost:8080/add?url=http://talks.golang.org/2015/state-of-go-may.slide#1
	curl http://localhost:8080/add?url=https://overcast.fm/+CHZ-pZJYk
	curl http://localhost:8080/add?url=http://gathering.tweakers.net/forum/list_bookmarks/////default
	curl http://localhost:8080/add?url=https://www.google.nl/search?q=create+image+from+page+golang&ie=UTF-8&oe=UTF-8&hl=nl&client=safari#hl=nl&q=golang+save+html+as+image
	curl http://localhost:8080/add?url=http://superuser.com/questions/272265/getting-curl-to-output-http-status-code
	curl http://localhost:8080/add?url=https://www.linkedin.com/
	curl http://localhost:8080/add?url=https://walledcity.com/supermighty/building-go-projects-with-gb
	curl http://localhost:8080/add?url=http://www.waterpop.nl/
	curl http://localhost:8080/add?url=http://golang-basic.blogspot.nl/2014/06/golang-database-step-by-step-guide-on.html
	curl http://localhost:8080/add?url=http://golangtutorials.blogspot.nl/2011/06/goroutines.html
	curl http://localhost:8080/add?url=https://www.reddit.com/r/golang/comments/35v5bm/best_way_to_run_go_server_as_a_daemon/cr86r5m
	curl http://localhost:8080/add?url=http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/
	curl http://localhost:8080/add?url=http://www.klimbosgarderen.nl/
	curl http://localhost:8080/add?url=https://github.com/angular-ui/ui-ace/issues/15
	curl http://localhost:8080/add?url=http://blog.doit.st/
	curl http://localhost:8080/add?url=http://www.geenstijl.nl/mt/archieven/2015/05/video_alcomobilist_rijdt_daiha.html
	curl http://localhost:8080/add?url=http://gathering.tweakers.net/forum/list_bookmarks/////default
	curl http://localhost:8080/add?url=https://www.facebook.com/?sk=h_chr
	curl http://localhost:8080/add?url=http://blog.doit.st/
	curl http://localhost:8080/add?url=https://github.com/angular-ui/ui-ace/issues/15
	curl http://localhost:8080/add?url=http://www.klimbosgarderen.nl/
	curl http://localhost:8080/add?url=http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/
	curl http://localhost:8080/add?url=https://www.reddit.com/r/golang/comments/35v5bm/best_way_to_run_go_server_as_a_daemon/cr86r5m
	curl http://localhost:8080/add?url=http://golangtutorials.blogspot.nl/2011/06/goroutines.html
	curl http://localhost:8080/add?url=http://golang-basic.blogspot.nl/2014/06/golang-database-step-by-step-guide-on.html
	curl http://localhost:8080/add?url=http://www.waterpop.nl/
	curl http://localhost:8080/add?url=https://walledcity.com/supermighty/building-go-projects-with-gb
	curl http://localhost:8080/add?url=http://www.geenstijl.nl/mt/archieven/2015/05/video_alcomobilist_rijdt_daiha.html

