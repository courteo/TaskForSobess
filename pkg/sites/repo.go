package sites

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNoSite  = errors.New("No site found")
)

type SiteMemoryRepository struct {
	data   *sql.DB
	LastID uint32
}

func (repo *SiteMemoryRepository) NewUserID() uint32 {
	repo.LastID++
	return repo.LastID
}

func NewMemoryRepo(db *sql.DB) *SiteMemoryRepository {
	return &SiteMemoryRepository{
		data:   db,
		LastID: 1,
	}
}

func (repo *SiteMemoryRepository) FindSite(name string) (*Site, error) {
	u := &Site{}
	row := repo.data.QueryRow("SELECT id, name, accessTime FROM sites WHERE name = ?", name)
	err := row.Scan(&u.ID, &u.Name, &u.AccessTime)
	if err != nil {
		return nil, ErrNoSite
	}
	return u, nil
}

func (repo *SiteMemoryRepository) Add(u *Site) error {
	_, err := repo.data.Exec(
		"INSERT INTO sites (`name`, `accessTime`, `ID`) VALUES (?, ?, ?)",
		u.Name,
		u.AccessTime,
		repo.NewUserID(),
	)
	return err
}

func (repo *SiteMemoryRepository) Update(name string, accessTime time.Duration) error {
	_, err := repo.data.Exec(
		"UPDATE sites SET accessTime = ? WHERE name = ?", 
		accessTime, 
		name,
	)
	return err
}

func (repo *SiteMemoryRepository) FindMinAccessTimeSite() (*Site, error) {
	u := &Site{}
	row := repo.data.QueryRow("SELECT name FROM sites ORDER BY accessTime ASC LIMIT 1")
	err := row.Scan(&u.Name)
	if err != nil {
		return nil, ErrNoSite
	}
	return u, nil
}

func (repo *SiteMemoryRepository) FindMaxAccessTimeSite() (*Site, error) {
	u := &Site{}
	row := repo.data.QueryRow("SELECT name FROM sites ORDER BY accessTime DESC LIMIT 1")
	err := row.Scan(&u.Name)
	if err != nil {
		return nil, ErrNoSite
	}
	return u, nil
}