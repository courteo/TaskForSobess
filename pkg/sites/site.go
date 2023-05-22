package sites

import "time"

type Site struct {
	ID 	uint32
	Name string
	AccessTime time.Duration
}

type SiteRepo interface {
	Add(u *Site) error
	FindMaxAccessTimeSite() (*Site, error)
	FindMinAccessTimeSite() (*Site, error)
	FindSite(name string) (*Site, error)
	Update(name string, accessTime time.Duration) error
}