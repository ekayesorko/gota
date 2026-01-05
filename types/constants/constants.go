package constants

import "github.com/ekayesorko/gota/types"

const UserIdContextKey = "user_id"
const RoleContextKey = "role"

const (
	RoleUser      types.Role = "user"
	RoleAdmin     types.Role = "admin"
	RoleAnonymous types.Role = "anonymous"
)

const (
	String types.Datatype = "string"
	Int    types.Datatype = "int"
	Float  types.Datatype = "float"
)
