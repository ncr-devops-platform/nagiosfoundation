package nagiosfoundation

import (
	"errors"
	"os/user"
	"testing"
)

const (
	goodString = "good"
	badString  = "bad"
)

type UserGroupTestHandler struct{}

func (u UserGroupTestHandler) testString(value string) error {
	var retval error

	if value == goodString {
		retval = nil
	} else {
		retval = errors.New("testing error")
	}

	return retval
}

func (u UserGroupTestHandler) Lookup(userName string) (*user.User, error) {
	return nil, u.testString(userName)
}

func (u UserGroupTestHandler) LookupGroup(groupName string) (*user.Group, error) {
	return nil, u.testString(groupName)
}

func (u UserGroupTestHandler) LookupGroupId(groupId string) (*user.Group, error) {
	return &user.Group{Name: goodString}, nil
}

func (u UserGroupTestHandler) GroupIds(userInfo *user.User) ([]string, error) {
	return []string{"1"}, nil
}

func TestCheckUser(t *testing.T) {
	var retval int

	if _, retval = (UserGroupCheck{
		UserName:  goodString,
		GroupName: "",
		Service:   new(UserGroupTestHandler),
	}).CheckUser(); retval != 0 {
		t.Error("CheckUser() with good user failed")
	}

	if _, retval = (UserGroupCheck{
		UserName:  badString,
		GroupName: "",
		Service:   new(UserGroupTestHandler),
	}).CheckUser(); retval != 3 {
		t.Error("CheckUser() with bad user failed")
	}
}

func TestCheckGroup(t *testing.T) {
	var retval int

	if _, retval = (UserGroupCheck{
		UserName:  "",
		GroupName: goodString,
		Service:   new(UserGroupTestHandler),
	}).CheckGroup(); retval != 0 {
		t.Error("CheckGroup() with good group failed")
	}

	if _, retval = (UserGroupCheck{
		UserName:  "",
		GroupName: badString,
		Service:   new(UserGroupTestHandler),
	}).CheckGroup(); retval != 3 {
		t.Error("CheckGroup() with bad group failed")
	}
}
func TestCheckUserGroup(t *testing.T) {
	var retval int

	if _, retval = (UserGroupCheck{
		UserName:  goodString,
		GroupName: goodString,
		Service:   new(UserGroupTestHandler),
	}).CheckUserGroup(); retval != 0 {
		t.Error("CheckGroup() with good user and good group failed")
	}

	if _, retval = (UserGroupCheck{
		UserName:  goodString,
		GroupName: badString,
		Service:   new(UserGroupTestHandler),
	}).CheckUserGroup(); retval != 3 {
		t.Error("CheckGroup() with good user and bad group failed")
	}

	if _, retval = (UserGroupCheck{
		UserName:  badString,
		GroupName: badString,
		Service:   new(UserGroupTestHandler),
	}).CheckUserGroup(); retval != 3 {
		t.Error("CheckGroup() with bad user and bad group failed")
	}

	if _, retval = (UserGroupCheck{
		UserName:  badString,
		GroupName: goodString,
		Service:   new(UserGroupTestHandler),
	}).CheckUserGroup(); retval != 3 {
		t.Error("CheckGroup() with bad user and good group failed")
	}
}
