package nagiosfoundation

import (
	"errors"
	"os/user"
	"testing"
)

const (
	goodUserGroupString  = "good"
	badUserGroupString   = "bad"
	errorUserGroupString = "error"
)

type userGroupTestHandler struct{}

func (u userGroupTestHandler) testString(value string) error {
	var retval error

	if value != badUserGroupString {
		retval = nil
	} else {
		retval = errors.New("testing error")
	}

	return retval
}

func (u userGroupTestHandler) Lookup(userName string) (*user.User, error) {
	user := user.User{
		Username: userName,
	}

	return &user, u.testString(userName)
}

func (u userGroupTestHandler) LookupGroup(groupName string) (*user.Group, error) {
	return nil, u.testString(groupName)
}

func (u userGroupTestHandler) LookupGroupID(groupID string) (*user.Group, error) {
	return &user.Group{Name: goodUserGroupString}, nil
}

func (u userGroupTestHandler) GroupIds(userInfo *user.User) ([]string, error) {
	groupIDList := []string{"1"}
	var err error

	if userInfo != nil && userInfo.Username == errorUserGroupString {
		groupIDList = nil
		err = errors.New("Error fetching Group IDs")
	}
	return groupIDList, err
}

func TestCheckUser(t *testing.T) {
	var retval int

	if _, retval = (UserGroupCheck{
		UserName:  goodUserGroupString,
		GroupName: "",
		Service:   new(userGroupTestHandler),
	}).CheckUser(); retval != statusCodeOK {
		t.Error("CheckUser() with good user failed")
	}

	if _, retval = (UserGroupCheck{
		UserName:  badUserGroupString,
		GroupName: "",
		Service:   new(userGroupTestHandler),
	}).CheckUser(); retval != statusCodeCritical {
		t.Error("CheckUser() with bad user failed")
	}
}

func TestCheckGroup(t *testing.T) {
	var retval int
	handler := new(userGroupTestHandler)

	if _, retval = (UserGroupCheck{
		UserName:  "",
		GroupName: goodUserGroupString,
		Service:   handler,
	}).CheckGroup(); retval != statusCodeOK {
		t.Error("CheckGroup() with good group failed")
	}

	if _, retval = (UserGroupCheck{
		UserName:  "",
		GroupName: badUserGroupString,
		Service:   handler,
	}).CheckGroup(); retval != statusCodeCritical {
		t.Error("CheckGroup() with bad group failed")
	}
}

func TestCheckUserGroup(t *testing.T) {
	var retval int
	handler := new(userGroupTestHandler)

	if _, retval = (UserGroupCheck{
		UserName:  goodUserGroupString,
		GroupName: goodUserGroupString,
		Service:   handler,
	}).CheckUserGroup(); retval != statusCodeOK {
		t.Error("CheckGroup() with good user and good group failed")
	}

	if _, retval = (UserGroupCheck{
		UserName:  errorUserGroupString,
		GroupName: goodUserGroupString,
		Service:   handler,
	}).CheckUserGroup(); retval != statusCodeCritical {
		t.Error("CheckGroup() with GroupIds() returning error failed")
	}

	if _, retval = (UserGroupCheck{
		UserName:  goodUserGroupString,
		GroupName: badUserGroupString,
		Service:   handler,
	}).CheckUserGroup(); retval != statusCodeCritical {
		t.Error("CheckGroup() with good user and bad group failed")
	}

	if _, retval = (UserGroupCheck{
		UserName:  badUserGroupString,
		GroupName: badUserGroupString,
		Service:   handler,
	}).CheckUserGroup(); retval != statusCodeCritical {
		t.Error("CheckGroup() with bad user and bad group failed")
	}

	if _, retval = (UserGroupCheck{
		UserName:  badUserGroupString,
		GroupName: goodUserGroupString,
		Service:   handler,
	}).CheckUserGroup(); retval != statusCodeCritical {
		t.Error("CheckGroup() with bad user and good group failed")
	}
}

func TestCheckUserGroupWithFlags(t *testing.T) {
	handler := new(userGroupTestHandler)

	_, err := CheckUserGroupWithHandler(goodUserGroupString, "", handler)
	if err != statusCodeOK {
		t.Error("CheckUserGroup with --user flag failed")
	}

	_, err = CheckUserGroupWithHandler("", goodUserGroupString, handler)
	if err != statusCodeOK {
		t.Error("CheckUserGroup with --group flag failed")
	}

	_, err = CheckUserGroupWithHandler(goodUserGroupString, goodUserGroupString, handler)
	if err != statusCodeOK {
		t.Error("CheckUserGroup with --user and --group flags failed")
	}
}

func TestOsUserCalls(t *testing.T) {
	ugh := UserGroupHandler{}

	currentUser, err := user.Current()
	if err != nil {
		t.Error("Lookingup current user")
	}

	// Lookup Group by ID
	groupByID, err := ugh.LookupGroupID(currentUser.Gid)
	if err != nil {
		t.Error("Lookup current user primary group ID")
	}

	// Lookup Group by Name
	groupByName, err := ugh.LookupGroup(groupByID.Name)
	if err != nil {
		t.Error("Lookup primary group ID by name")
	}

	if groupByID.Gid != groupByName.Gid {
		t.Error("Lookup group by ID does not match lookup group by name")
	}

	_, result := CheckUserGroup(currentUser.Username, groupByID.Name)
	if result != 0 {
		t.Error("Check for current user name and group ID")
	}
}
