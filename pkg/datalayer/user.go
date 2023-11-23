package datalayer

type User struct {
	ID             uint64 `db:"user_id"`
	Email          string `db:"email"`
	PasswordDigest string `db:"passwd"`
	Name           string `db:"username"`
}
