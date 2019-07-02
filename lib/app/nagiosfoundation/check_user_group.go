package nagiosfoundation

import (
	"fmt"
	"os/user"
)

const (
	checkUserName      = "CheckUser"
	checkGroupName     = "CheckGroup"
	checkUserGroupName = "CheckUserGroup"
)

// UserGroupService is an interface that allows for overriding
// methods for looking up a user, a group, finding the ID for
// a group name, and fetching a list of groups to which a user
// belongs.
type UserGroupService interface {
	Lookup(string) (*user.User, error)
	LookupGroup(string) (*user.Group, error)
	LookupGroupID(string) (*user.Group, error)
	GroupIds(*user.User) ([]string, error)
}

// UserGroupHandler is a specific implementation of the methods
// in the UserGroupService interface. The methods implemented use
// Go's user package which should be compatible with both Linux
// and Windows.
type UserGroupHandler struct{}

// Lookup is given a user name and returns the details of a user
// on the host OS. This implementation calls Lookup() in Go's
// user package.
func (u UserGroupHandler) Lookup(userName string) (*user.User, error) {
	return user.Lookup(userName)
}

// LookupGroup is given a group name and returns the details of
// the group on the host OS. This implementation calls
// LookupGroup() in Go's user package.
func (u UserGroupHandler) LookupGroup(groupName string) (*user.Group, error) {
	return user.LookupGroup(groupName)
}

// LookupGroupID is given a group ID and returns the details of
// the group on the host OS. This implmentation calls
// LookupGroupId() in Go's user package.
func (u UserGroupHandler) LookupGroupID(groupID string) (*user.Group, error) {
	return user.LookupGroupId(groupID)
}

// GroupIds is given the User struct from Go's user package and
// returns a list of group IDs to which the user belongs.
func (u UserGroupHandler) GroupIds(userInfo *user.User) ([]string, error) {
	return userInfo.GroupIds()
}

// UserGroupCheck is a struct in which a user and group
// on the host OS can be populated, then information about
// the user and group can be retrieved.
type UserGroupCheck struct {
	UserName  string
	GroupName string

	Service UserGroupService
}

// CheckUser checks for the existence of a user on the host
// operating system.
//
// Returns a string containing plaintext result and an integer
// with the return code. 0 indicates the user exists, 3 indicates
// the user does not exist.
func (ugc UserGroupCheck) CheckUser() (string, int) {
	var retCode int
	var formatString, status string

	_, err := ugc.Service.Lookup(ugc.UserName)

	if err == nil {
		formatString = "User %s exists"
		status = statusTextOK
		retCode = 0
	} else {
		formatString = "User %s does not exist"
		status = statusTextCritical
		retCode = 3
	}

	msg, _ := resultMessage(checkUserName, status, fmt.Sprintf(formatString, ugc.UserName))
	return msg, retCode
}

// CheckGroup checks for the existence of a group on the host
// operating system.
//
// Returns a string containing plaintext result and an integer
// with the return code. 0 indicates the group exists, 3 indicates
// the group does not exist.
func (ugc UserGroupCheck) CheckGroup() (string, int) {
	var retCode int
	var formatString, status string

	_, err := ugc.Service.LookupGroup(ugc.GroupName)

	if err == nil {
		formatString = "Group %s exists"
		status = statusTextOK
		retCode = 0
	} else {
		formatString = "Group %s does not exist"
		status = statusTextCritical
		retCode = 3
	}

	msg, _ := resultMessage(checkGroupName, status, fmt.Sprintf(formatString, ugc.GroupName))
	return msg, retCode
}

// CheckUserGroup checks for the existence of a user and if it
// belongs to the named group on the host operating system.
//
// Returns a string containing plaintext result and an integer
// with the return code. 0 indicates the user exists and is in
// the group, 3 indicates otherwise.
func (ugc UserGroupCheck) CheckUserGroup() (string, int) {
	var msg, status string
	retcode := 3

	userInfo, err := ugc.Service.Lookup(ugc.UserName)
	if err != nil {
		msg = fmt.Sprintf("User %s does not exist", ugc.UserName)
		status = statusTextCritical
	} else {
		groupIds, err := ugc.Service.GroupIds(userInfo)
		if err != nil {
			msg = fmt.Sprintf("Could not get Group IDs for user %s", ugc.UserName)
			status = statusTextCritical
		} else {
			msg = fmt.Sprintf("User %s exists but is not in Group %s",
				ugc.UserName, ugc.GroupName)
			status = statusTextCritical

			for i := range groupIds {
				groupInfo, _ := ugc.Service.LookupGroupID(groupIds[i])
				if groupInfo.Name == ugc.GroupName {
					msg = fmt.Sprintf("User %s exists and is in Group %s",
						ugc.UserName, ugc.GroupName)
					status = statusTextOK
					retcode = 0
				}
			}
		}
	}

	msg, _ = resultMessage(checkUserGroupName, status, msg)
	return msg, retcode
}

// CheckUserGroupWithHandler checks for the existence of a user and
// if it belongs to the named group on the host operating system.
//
// The user and group are provided with the user and group parameters.
//
// Returns - result message and return code.
func CheckUserGroupWithHandler(user, group string, userGroupHandler UserGroupService) (string, int) {
	var msg string
	var retCode int

	userGroupCheck := UserGroupCheck{
		UserName:  user,
		GroupName: group,
		Service:   userGroupHandler,
	}

	if user != "" && group != "" {
		msg, retCode = userGroupCheck.CheckUserGroup()
	} else if user != "" {
		msg, retCode = userGroupCheck.CheckUser()
	} else if group != "" {
		msg, retCode = userGroupCheck.CheckGroup()
	}

	return msg, retCode
}

// CheckUserGroup executes the normal user/group
// checks in CheckUserGroupFlags() then prints the results
// and exits.
func CheckUserGroup(user, group string) (string, int) {
	return CheckUserGroupWithHandler(user, group, new(UserGroupHandler))
}
