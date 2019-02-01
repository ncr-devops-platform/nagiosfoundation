package nagiosfoundation

import (
	"flag"
	"fmt"
	"os"
	"os/user"
)

// Checks for the existence of a user on the host operating system.
//
// Returns a string containing plaintext result and an integer
// with the return code. 0 indicates the user exists, 3 indicates
// the user does not exist.
func CheckUser(userName string) (string, int) {
	var retCode int
	var formatString string

	_, err := user.Lookup(userName)

	if err == nil {
		formatString = "%s OK - User %s exists"
		retCode = 0
	} else {
		formatString = "%s CRITICAL - User %s does not exist"
		retCode = 3
	}

	return fmt.Sprintf(formatString, "CheckUser", userName), retCode
}

// Checks for the existence of a group on the host operating system.
//
// Returns a string containing plaintext result and an integer
// with the return code. 0 indicates the group exists, 3 indicates
// the group does not exist.
func CheckGroup(groupName string) (string, int) {
	var retCode int
	var formatString string

	_, err := user.LookupGroup(groupName)

	if err == nil {
		formatString = "%s OK - Group %s exists"
		retCode = 0
	} else {
		formatString = "%s CRITICAL - Group %s does not exist"
		retCode = 3
	}

	return fmt.Sprintf(formatString, "CheckGroup", groupName), retCode
}

// Checks for the existence of a user and if it belongs to the
// named group on the host operating system.
//
// Returns a string containing plaintext result and an integer
// with the return code. 0 indicates the user exists and is in
// the group, 3 indicates otherwise.
func CheckUserGroup(userName string, groupName string) (string, int) {
	var msg string
	retcode := 3

	userInfo, err := user.Lookup(userName)
	if err != nil {
		msg = fmt.Sprintf("CheckUserGroup CRITICAL - User %s does not exist", userName)
	} else {
		groupIds, err := userInfo.GroupIds()
		if err != nil {
			msg = fmt.Sprintf("CheckUserGroup CRITICAL - Could not get Group IDs for user %s", userName)
		}

		msg = fmt.Sprintf("CheckUserGroup CRITICAL - User %s exists but is not in Group %s", userName, groupName)
		for i := range groupIds {
			groupInfo, _ := user.LookupGroupId(groupIds[i])
			if groupInfo.Name == groupName {
				msg = fmt.Sprintf("CheckUserGroup OK - User %s exists and is in Group %s", userName, groupName)
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

	if *userPtr != "" && *groupPtr != "" {
		msg, retCode = CheckUserGroup(*userPtr, *groupPtr)
	} else if *userPtr != "" {
		msg, retCode = CheckUser(*userPtr)
	} else if *groupPtr != "" {
		msg, retCode = CheckGroup(*groupPtr)
	}

	fmt.Println(msg)
	os.Exit(retCode)
}
