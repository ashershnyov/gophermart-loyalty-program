package model

// IntUser defines user data internally.
type IntUser struct {
	ID       int64  `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
}
