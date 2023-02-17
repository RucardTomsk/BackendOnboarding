package enum

type Roles int

const (
	BEGINNER Roles = iota
	ADMIN
	EMPLOYEE
	DIRECTOR
	HR
)

func (i Roles) String() string {
	return [...]string{"beginner", "admin", "employee", "director", "hr"}[i]
}

func ParseRoles(rolesString string) Roles {
	switch rolesString {
	case "beginner":
		return BEGINNER
	case "admin":
		return ADMIN
	case "employee":
		return EMPLOYEE
	case "director":
		return DIRECTOR
	case "hr":
		return HR
	default:
		return BEGINNER
	}
}
