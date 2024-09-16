// Code generated from Pkl module `org.kdeps.pkl.APIServer`. DO NOT EDIT.
package apiserver

type APIServerSettings struct {
	HostIP string `pkl:"hostIP"`

	PortNum uint16 `pkl:"portNum"`

	Routes []*APIServerRoutes `pkl:"routes"`
}
