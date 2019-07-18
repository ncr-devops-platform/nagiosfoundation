# User and Group Check

The user and group check (`check_user_group`) will check for the existence of a user, of a group, or of a user being in a group, depending on the command line flags given.

All features of this check are applicable to both Linux and Windows.

## Check for a User
Use the `--user` flag to check for the existence of a user. It is the only command like flag required for this check.

```
check_user_group --user nobody
```

## Check for a Group

Use the `--group` flag to check for the existence of a group. It is the only command line flag required for this check.

```
check_user_group --group sudo
```

## Check for a User Belonging to a Group
Use both the `--user` and `--group` flags to check if a user exists and verify the user is in a group.

```
check_user_group --user adm --group syslog 
```
