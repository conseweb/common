package registrar

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/conseweb/common/fabricGoSDK/models"
)

// GetUserRegistrationReader is a Reader for the GetUserRegistration structure.
type GetUserRegistrationReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the recieved o.
func (o *GetUserRegistrationReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetUserRegistrationOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewGetUserRegistrationDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	}
}

// NewGetUserRegistrationOK creates a GetUserRegistrationOK with default headers values
func NewGetUserRegistrationOK() *GetUserRegistrationOK {
	return &GetUserRegistrationOK{}
}

/*GetUserRegistrationOK handles this case with default header values.

Confirm registration for target user
*/
type GetUserRegistrationOK struct {
	Payload *models.OK
}

func (o *GetUserRegistrationOK) Error() string {
	return fmt.Sprintf("[GET /registrar/{enrollmentID}][%d] getUserRegistrationOK  %+v", 200, o.Payload)
}

func (o *GetUserRegistrationOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.OK)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetUserRegistrationDefault creates a GetUserRegistrationDefault with default headers values
func NewGetUserRegistrationDefault(code int) *GetUserRegistrationDefault {
	return &GetUserRegistrationDefault{
		_statusCode: code,
	}
}

/*GetUserRegistrationDefault handles this case with default header values.

Unexpected error
*/
type GetUserRegistrationDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the get user registration default response
func (o *GetUserRegistrationDefault) Code() int {
	return o._statusCode
}

func (o *GetUserRegistrationDefault) Error() string {
	return fmt.Sprintf("[GET /registrar/{enrollmentID}][%d] getUserRegistration default  %+v", o._statusCode, o.Payload)
}

func (o *GetUserRegistrationDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
