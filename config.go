package main

const jsonFileName = "config.json"

type configJson struct {
	Login  loginParameters
	Logout logoutParameters

	Endpoints string
	Results   string

	Domains [2]string

	Timeout int
}

type loginParameters struct {
	Endpoint string
	Fields   loginFields
}

type loginFields struct {
	Username [2]string
	Password [2]string
}

type logoutParameters struct {
	Endpoint string
}

func (f loginFields) getUsernameField() string {
	return f.Username[0]
}

func (f loginFields) getUsernameValue() string {
	return f.Username[1]
}

func (f loginFields) getPasswordField() string {
	return f.Password[0]
}

func (f loginFields) getPasswordValue() string {
	return f.Password[1]
}
