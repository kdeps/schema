// Code generated from Pkl module `org.kdeps.pkl.APIServer`. DO NOT EDIT.
package apiserver

type APIServerSettings struct {
	ServerPort int `pkl:"serverPort"`

	Routes []*APIServerRoutes `pkl:"routes"`
}
