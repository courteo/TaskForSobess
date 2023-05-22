package admin

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrNoAdmin  = errors.New(" No Admin found")
	ErrBadPass = errors.New(" Invald password")
)

type AdminsMemoryRepository struct {
	data   *sql.DB
	LastID uint32
}

func (repo *AdminsMemoryRepository) NewAdminID() uint32 {
	repo.LastID++
	return repo.LastID
}

func NewMemoryRepo(db *sql.DB) *AdminsMemoryRepository {
	return &AdminsMemoryRepository{
		data:   db,
		LastID: 1,
	}
}

func (repo *AdminsMemoryRepository) FindAdmin(login string) (*Admin, error) {
	u := &Admin{}
	row := repo.data.QueryRow("SELECT id, login, password FROM admins WHERE login = ?", login)
	err := row.Scan(&u.ID, &u.Login, &u.Password)
	if err != nil {
		return nil, ErrNoAdmin
	}
	return u, nil
}

func (repo *AdminsMemoryRepository) Authorize(login, pass string) (*Admin, error) {
	u := &Admin{}
	row := repo.data.QueryRow("SELECT id, login, password FROM admins WHERE login = ?", login)
	err := row.Scan(&u.ID, &u.Login, &u.Password)
	fmt.Println("row ", row)
	if err != nil {
		return nil, ErrNoAdmin
	}

	if u.Password != pass {
		return nil, ErrBadPass
	}

	return u, nil
}

func (repo *AdminsMemoryRepository) Add(u *Admin) error {

	_, err := repo.data.Exec(
		"INSERT INTO admins (`login`, `password`, `ID`) VALUES (?, ?, ?)",
		u.Login,
		u.Password,
		repo.NewAdminID(),
	)
	return err
}