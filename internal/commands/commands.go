package commands

import (
	"os/exec"
	"strings"
)

/***********************************************
* These currecntly only work for Unix. Must be *
* expanded to support Linux and Win in future. *
***********************************************/

func GetMachineUsersList() ([]string, error) {
	users := make([]string, 0)

	cmd := exec.Command("users")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	users = append(users, strings.Split(string(out), " ")...)
	return users, nil
}

func GetCurrentUser() (string, error) {
	cmd := exec.Command("whoami")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
