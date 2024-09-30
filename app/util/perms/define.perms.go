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
	PERMS_ALL                 Permission = "ALL"
	PERMS_VIEW                Permission = "VIEW"
	PERMS_VIEW_GROUP          Permission = "VIEW_GROUP"
	PERMS_VIEW_USER           Permission = "VIEW_USER"
	PERMS_VIEW_PERMISSION     Permission = "VIEW_PERMISSION"
	PERMS_VIEW_DEPARTMENT     Permission = "VIEW_DEPARTMENT"
	PERMS_VIEW_BRANCH         Permission = "VIEW_BRANCH"
	PERMS_VIEW_ALL            Permission = "VIEW_ALL"
	PERMS_MUTATION_GROUP      Permission = "MUTATION_GROUP"
	PERMS_MUTATION_USER       Permission = "MUTATION_USER"
	PERMS_MUTATION_PERMISSION Permission = "MUTATION_PERMISSION"
	PERMS_MUTATION_ALL        Permission = "MUTATION_ALL"
	PERMS_SETUP_GROUP         Permission = "SETUP_GROUP"
	PERMS_SETUP_USER          Permission = "SETUP_USER"
	PERMS_SETUP_ALL           Permission = "SETUP_ALL"
)
