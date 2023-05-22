package admin

type Admin struct {
	ID       uint32
	Login    string
	Password string
}

type AdminsRepo interface {
	FindAdmin(login string) (*Admin, error)
	Authorize(login, pass string) (*Admin, error)
	NewAdminID() uint32
	Add(u *Admin) error
}
