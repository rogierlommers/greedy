package articles

import "time"

// Stats holds information about imported assets
type Stats struct {
	LastCrawled time.Time
	LastCrawler string
	CrawlCount  int
}

func (s *Stats) incCrawlCount() {
	s.CrawlCount++
}

func (s *Stats) setLastCrawler(refer string) {
	s.LastCrawler = refer
	s.LastCrawled = time.Now()
}
