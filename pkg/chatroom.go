package chatapp

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

const SystemID = "system"

// SystemMessage generates a message with the given content for the system.
//
// Parameters:
// - message: the content of the message.
//
// Returns:
// - Message: the generated message with the system ID as the username and the provided message.
func SystemMessage(message string) Message {
	return Message{
		Username: SystemID,
		Message:  message,
	}
}

// NewMessage creates a new Message object.
//
// It takes in a username and a message as parameters, both of type string.
// It returns a Message object.
func NewMessage(username, message string) Message {
	return Message{
		Username: username,
		Message:  message,
	}
}
