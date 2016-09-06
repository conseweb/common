package chaincode

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/conseweb/common/fabricGoSDK/models"
)

// ChaincodeDeployReader is a Reader for the ChaincodeDeploy structure.
type ChaincodeDeployReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the recieved o.
func (o *ChaincodeDeployReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewChaincodeDeployOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewChaincodeDeployDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	}
}

// NewChaincodeDeployOK creates a ChaincodeDeployOK with default headers values
func NewChaincodeDeployOK() *ChaincodeDeployOK {
	return &ChaincodeDeployOK{}
}

/*ChaincodeDeployOK handles this case with default header values.

Successfully deployed chainCode
*/
type ChaincodeDeployOK struct {
	Payload *models.OK
}

func (o *ChaincodeDeployOK) Error() string {
	return fmt.Sprintf("[POST /devops/deploy][%d] chaincodeDeployOK  %+v", 200, o.Payload)
}

func (o *ChaincodeDeployOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.OK)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewChaincodeDeployDefault creates a ChaincodeDeployDefault with default headers values
func NewChaincodeDeployDefault(code int) *ChaincodeDeployDefault {
	return &ChaincodeDeployDefault{
		_statusCode: code,
	}
}

/*ChaincodeDeployDefault handles this case with default header values.

Unexpected error
*/
type ChaincodeDeployDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the chaincode deploy default response
func (o *ChaincodeDeployDefault) Code() int {
	return o._statusCode
}

func (o *ChaincodeDeployDefault) Error() string {
	return fmt.Sprintf("[POST /devops/deploy][%d] chaincodeDeploy default  %+v", o._statusCode, o.Payload)
}

func (o *ChaincodeDeployDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
