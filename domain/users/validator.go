package users

import (
	"errors"
	"strings"
)

func validateUser(usr *User) error {
	if err := validateFirstName(usr.FirstName); err != nil {
		return err
	}
	if err := validateLastName(usr.LastName); err != nil {
		return err
	}
	if err := validateNickName(usr.NickName); err != nil {
		return err
	}
	if err := validateEmail(usr.Email); err != nil {
		return err
	}
	return nil
}

func validateFirstName(firstName string) error {
	if len(firstName) < 2 || len(firstName) > 20 {
		return errors.New("firstname should have a len greater than 2 and less than 20")
	}
	return nil
}

func validateLastName(lastName string) error {
	if len(lastName) < 4 || len(lastName) > 40 {
		return errors.New("lastname should have a len greater than 4 and less than 40")
	}
	return nil
}

func validateNickName(nickName string) error {
	if len(nickName) < 4 || len(nickName) > 20 {
		return errors.New("nickname should have a len greater than 4 and less than 20")
	}
	return nil
}

func validateEmail(email string) error {
	part := strings.Split(email, "@")
	if len(part) != 2 {
		return errors.New("an email should contains one '@'")
	}
	part = strings.Split(part[1], ".")
	if len(part) != 2 {
		return errors.New("an email should contains one '.' on host part")
	}
	return nil
}
