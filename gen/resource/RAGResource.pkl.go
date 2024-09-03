// Code generated from Pkl module `org.kdeps.pkl.Resource`. DO NOT EDIT.
package resource

type RAGResource struct {
	Id string `pkl:"id"`

	Name string `pkl:"name"`

	Description string `pkl:"description"`

	Category string `pkl:"category"`

	Requires *[]string `pkl:"requires"`

	Action *[]*RAGAction `pkl:"action"`
}
