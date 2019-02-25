package nagiosfoundation

import (
	"flag"
	"fmt"
	"os"
	"os/user"
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
	var formatString string

	_, err := ugc.Service.Lookup(ugc.UserName)

	if err == nil {
		formatString = "%s OK - User %s exists"
		retCode = 0
	} else {
		formatString = "%s CRITICAL - User %s does not exist"
		retCode = 3
	}

	return fmt.Sprintf(formatString, "CheckUser", ugc.UserName), retCode
}

// CheckGroup checks for the existence of a group on the host
// operating system.
//
// Returns a string containing plaintext result and an integer
// with the return code. 0 indicates the group exists, 3 indicates
// the group does not exist.
func (ugc UserGroupCheck) CheckGroup() (string, int) {
	var retCode int
	var formatString string

	_, err := ugc.Service.LookupGroup(ugc.GroupName)

	if err == nil {
		formatString = "%s OK - Group %s exists"
		retCode = 0
	} else {
		formatString = "%s CRITICAL - Group %s does not exist"
		retCode = 3
	}

	return fmt.Sprintf(formatString, "CheckGroup", ugc.GroupName), retCode
}

// CheckUserGroup checks for the existence of a user and if it
// belongs to the named group on the host operating system.
//
// Returns a string containing plaintext result and an integer
// with the return code. 0 indicates the user exists and is in
// the group, 3 indicates otherwise.
func (ugc UserGroupCheck) CheckUserGroup() (string, int) {
	var msg string
	retcode := 3

	userInfo, err := ugc.Service.Lookup(ugc.UserName)
	if err != nil {
		msg = fmt.Sprintf("CheckUserGroup CRITICAL - User %s does not exist", ugc.UserName)
	} else {
		groupIds, err := ugc.Service.GroupIds(userInfo)
		if err != nil {
			msg = fmt.Sprintf("CheckUserGroup CRITICAL - Could not get Group IDs for user %s", ugc.UserName)
		} else {
			msg = fmt.Sprintf("CheckUserGroup CRITICAL - User %s exists but is not in Group %s",
				ugc.UserName, ugc.GroupName)

			for i := range groupIds {
				groupInfo, _ := ugc.Service.LookupGroupID(groupIds[i])
				if groupInfo.Name == ugc.GroupName {
					msg = fmt.Sprintf("CheckUserGroup OK - User %s exists and is in Group %s",
						ugc.UserName, ugc.GroupName)
					retcode = 0
				}
			}
		}
	}

	return msg, retcode
}

// CheckUserGroupFlagsWithHandler checks for the existence of a user and
// if it belongs to the named group on the host operating system.
// It does this based on command line flags.
//
// The user and group are provided with
// the command line flags, -user, and -group, respectively.
//
// Returns - result message and return code.
func CheckUserGroupFlagsWithHandler(userGroupHandler UserGroupService) (string, int) {
	userPtr := flag.String("user", "", "user name")
	groupPtr := flag.String("group", "", "group name")

	flag.Parse()

	var msg string
	var retCode int

	userGroupCheck := UserGroupCheck{
		UserName:  *userPtr,
		GroupName: *groupPtr,
		Service:   userGroupHandler,
	}

	if *userPtr != "" && *groupPtr != "" {
		msg, retCode = userGroupCheck.CheckUserGroup()
	} else if *userPtr != "" {
		msg, retCode = userGroupCheck.CheckUser()
	} else if *groupPtr != "" {
		msg, retCode = userGroupCheck.CheckGroup()
	}

	return msg, retCode
}

// CheckUserGroupFlagsWithExit executes the normal user/group
// checks in CheckUserGroupFlags() then prints the results
// and exits.
func CheckUserGroupFlagsWithExit() {
	msg, retCode := CheckUserGroupFlagsWithHandler(new(UserGroupHandler))

	fmt.Println(msg)
	os.Exit(retCode)
}
