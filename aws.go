package main

type awsCreds struct {
}

type LoginManager interface {
	func(awsCreds) (bool, error)
}

type ssoLogin struct {
}

func getAllProfiles() ([]string, error) {
	return nil, nil
}
