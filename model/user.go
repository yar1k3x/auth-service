package model

type User struct {
	ID       int64  `db:"id"`
	FIO      string `db:"fio"`
	Email    string `db:"email"`
	Password string `db:"password"` // hashed
	Role     string `db:"role"`     // "admin" or "user"
}
