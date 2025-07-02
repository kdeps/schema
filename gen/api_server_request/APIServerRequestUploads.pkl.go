// Code generated from Pkl module `org.kdeps.pkl.APIServerRequest`. DO NOT EDIT.
package apiserverrequest

// Class representing metadata for an uploaded file in an API request.
type APIServerRequestUploads struct {
	// The file path where the uploaded file is stored on the server.
	Filepath string `pkl:"Filepath"`

	// The MIME type of the uploaded file.
	Filetype string `pkl:"Filetype"`
}
