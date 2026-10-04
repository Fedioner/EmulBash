package host

import (
	"os"
	"os/user"
)

type Info struct {
	Username string 
	Hostname string
}

func GetInfo() (Info, error) {
	user, err := user.Current()
	if err != nil {
		return Info{}, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return Info{}, err
	}

	return Info{
		Username: user.Username,
		Hostname: hostname,
	}, nil


}