package nagiosfoundation

import (
	"flag"
	"fmt"
	"os"
	"os/user"
)

type UserGroupService interface {
	Lookup(string) (*user.User, error)
	LookupGroup(string) (*user.Group, error)
	LookupGroupId(string) (*user.Group, error)
	GroupIds(*user.User) ([]string, error)
}

type UserGroupHandler struct{}

func (u UserGroupHandler) Lookup(userName string) (*user.User, error) {
	return user.Lookup(userName)
}

func (u UserGroupHandler) LookupGroup(groupName string) (*user.Group, error) {
	return user.LookupGroup(groupName)
}

func (u UserGroupHandler) LookupGroupId(groupName string) (*user.Group, error) {
	return user.LookupGroupId(groupName)
}

func (u UserGroupHandler) GroupIds(userInfo *user.User) ([]string, error) {
	return userInfo.GroupIds()
}

type UserGroupCheck struct {
	UserName  string
	GroupName string

	Service UserGroupService
}

// Checks for the existence of a user on the host operating system.
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

// Checks for the existence of a group on the host operating system.
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

// Checks for the existence of a user and if it belongs to the
// named group on the host operating system.
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
		}

		msg = fmt.Sprintf("CheckUserGroup CRITICAL - User %s exists but is not in Group %s",
			ugc.UserName, ugc.GroupName)
		for i := range groupIds {
			groupInfo, _ := ugc.Service.LookupGroupId(groupIds[i])
			if groupInfo.Name == ugc.GroupName {
				msg = fmt.Sprintf("CheckUserGroup OK - User %s exists and is in Group %s",
					ugc.UserName, ugc.GroupName)
				retcode = 0
			}
		}
	}

	return msg, retcode
}

// Checks for the existence of a user and if it belongs to the
// named group on the host operating system.
//
// The user and group are provided with
// the command line flags, -user, and -group, respectively.
// This function then exist with the return code of the
// user, group, or usergroup function.
//
// Returns - see CheckUser(), CheckGroup, and CheckUserGroup()
func CheckUserGroupFlags() {
	userPtr := flag.String("user", "", "user name")
	groupPtr := flag.String("group", "", "group name")

	flag.Parse()

	var msg string
	var retCode int

	userGroupCheck := UserGroupCheck{
		UserName:  *userPtr,
		GroupName: *groupPtr,
		Service:   new(UserGroupHandler),
	}

	if *userPtr != "" && *groupPtr != "" {
		msg, retCode = userGroupCheck.CheckUserGroup()
	} else if *userPtr != "" {
		msg, retCode = userGroupCheck.CheckUser()
	} else if *groupPtr != "" {
		msg, retCode = userGroupCheck.CheckGroup()
	}

	fmt.Println(msg)
	os.Exit(retCode)
}
