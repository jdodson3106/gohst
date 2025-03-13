package auth

type Enviornment int

const (
	// List of environments
	Local Enviornment = iota
	Dev
	QA
	Stage
	Prod
)

type Profile struct {
	Name        string
	Enviornment Enviornment
}

func (p Profile) EnvName() string {
	switch p.Enviornment {
	case 0:
		return "local"
	case 1:
		return "dev"
	case 2:
		return "qa"
	case 3:
		return "stage"
	case 4:
		return "prod"
	default:
		return ""
	}
}
