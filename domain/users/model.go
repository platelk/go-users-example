package users

import "time"

// Operation define the action taken to the model. for example, did we create a new user.
type Operation string

var (
	// CreateOp is the operation by which the model start.
	// after this operation, the model can be requested.
	//
	// No `Before` model will be set on this operation, because no model existed previously
	//
	// Note: No ordering is guaranted by design
	CreateOp Operation = "create"

	// UpdateOp define an update operation of the model. only the ID is consistent and wont be change by an update.
	//
	// `Before` and `After` model will be filled on the event
	//
	// Note: No ordering is guaranted by design
	UpdateOp Operation = "update"

	// DeleteOp define a deletion of the model from the stores.
	//
	// No `After` model will be filled on the event, as the operation delete the user
	//
	// Note: No ordering is guaranted by design
	DeleteOp Operation = "delete"
)

// ChangeEvent will be emitted on each change in the user base to notify other systems of the changes.
type ChangeEvent struct {
	Time   time.Time
	Op     Operation
	Before *User
	After  *User
}

// User hold the definition of what is a user in the system
type User struct {
	// ID is a unique id across the system which is hidden to external users
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	NickName  string `json:"nick_name"`
	// Password is the hash representation of the user password
	// Note: as is, its printed on if logged, but as it contains hash representation is not that bad, but still should be improved
	Password string `json:"password"`
	Email    string `json:"email"`
	Country  string `json:"country"` // Note: here it should be a defined list and not an open field
}
