package domain

type Role string

const (
	RoleWaiter  Role = "waiter"
	RoleChef    Role = "chef"
	RoleManager Role = "manager"
)

type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Role      Role   `json:"role"`
	CreatedAt string `json:"created_at"`
}
