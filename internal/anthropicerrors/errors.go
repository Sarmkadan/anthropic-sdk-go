package anthropicerrors

import "fmt"

// AnthropicSdkGoException is the base error type for this SDK.
type AnthropicSdkGoException struct {
	Message string
	Err     error
}

func (e *AnthropicSdkGoException) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AnthropicSdkGoException) Unwrap() error {
	return e.Err
}

// ConfigurationException indicates an error in the SDK configuration.
type ConfigurationException struct {
	AnthropicSdkGoException
}

func NewConfigurationException(message string, err error) *ConfigurationException {
	return &ConfigurationException{
		AnthropicSdkGoException{Message: message, Err: err},
	}
}

// ValidationException indicates an error in request validation.
type ValidationException struct {
	AnthropicSdkGoException
}

func NewValidationException(message string, err error) *ValidationException {
	return &ValidationException{
		AnthropicSdkGoException{Message: message, Err: err},
	}
}
