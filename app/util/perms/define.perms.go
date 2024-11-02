package perms

type (
	TypeApi    string
	Permission string
)

const (
	USER_MANAGEMENT TypeApi = "USER_MANAGEMENT"
	DASHBOARD       TypeApi = "DASHBOARD"
	PROJECT         TypeApi = "PROJECT"
	LAYOUT          TypeApi = "LAYOUT"
	SERVICE         TypeApi = "SERVICE"
	VERSIONS        TypeApi = "VERSIONS"
	QUEUE           TypeApi = "QUEUE"

	ALL TypeApi = "ALL"
)

const (
	PERMS_ALL        Permission = "ALL"
	PERMS_VIEW       Permission = "VIEW"
	PERMS_VIEW_GROUP Permission = "VIEW_GROUP"
	PERMS_VIEW_USER  Permission = "VIEW_USER"
)
