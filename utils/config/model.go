package config

import "strings"

type Environment string

func (e Environment) String() string {
	return string(e)
}

func (e Environment) ToLower() string {
	return strings.ToLower(string(e))
}
