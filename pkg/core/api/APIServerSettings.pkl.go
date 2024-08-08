// Code generated from Pkl module `org.kdeps.pkl.API`. DO NOT EDIT.
package api

type APIServerSettings struct {
	ServerPort int `pkl:"serverPort"`

	Routes []*APIServerRoutes `pkl:"routes"`
}
