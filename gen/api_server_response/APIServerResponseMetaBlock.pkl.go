// Code generated from Pkl module `org.kdeps.pkl.APIServerResponse`. DO NOT EDIT.
package apiserverresponse

// Contains metadata related to an API response.
//
// This block includes essential details such as the request ID, response headers,
// and custom properties, providing additional context for API interactions.
type APIServerResponseMetaBlock struct {
	// A unique identifier (UUID) for the request.
	//
	// This ID helps track and correlate API requests.
	RequestID *string `pkl:"requestID"`

	// HTTP headers included in the API response.
	//
	// Contains key-value pairs representing response headers.
	Headers *map[string]string `pkl:"headers"`

	// Custom key-value properties included in the JSON response.
	//
	// Used to store additional metadata or context-specific details.
	Properties *map[string]string `pkl:"properties"`
}
