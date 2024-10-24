// Code generated from Pkl module `org.kdeps.pkl.APIServerRequest`. DO NOT EDIT.
package apiserverrequest

// Represents metadata for an uploaded file, including its file path and MIME type.
type APIServerRequestUploads struct {
	// The file path where the uploaded file is stored.
	Filepath string `pkl:"filepath"`

	// The MIME type of the uploaded file.
	Filetype string `pkl:"filetype"`
}
