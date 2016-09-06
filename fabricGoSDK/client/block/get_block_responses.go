package block

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/conseweb/common/fabricGoSDK/models"
)

// GetBlockReader is a Reader for the GetBlock structure.
type GetBlockReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the recieved o.
func (o *GetBlockReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetBlockOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewGetBlockDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	}
}

// NewGetBlockOK creates a GetBlockOK with default headers values
func NewGetBlockOK() *GetBlockOK {
	return &GetBlockOK{}
}

/*GetBlockOK handles this case with default header values.

Individual Block contents
*/
type GetBlockOK struct {
	Payload *models.Block
}

func (o *GetBlockOK) Error() string {
	return fmt.Sprintf("[GET /chain/blocks/{Block}][%d] getBlockOK  %+v", 200, o.Payload)
}

func (o *GetBlockOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Block)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetBlockDefault creates a GetBlockDefault with default headers values
func NewGetBlockDefault(code int) *GetBlockDefault {
	return &GetBlockDefault{
		_statusCode: code,
	}
}

/*GetBlockDefault handles this case with default header values.

Unexpected error
*/
type GetBlockDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the get block default response
func (o *GetBlockDefault) Code() int {
	return o._statusCode
}

func (o *GetBlockDefault) Error() string {
	return fmt.Sprintf("[GET /chain/blocks/{Block}][%d] getBlock default  %+v", o._statusCode, o.Payload)
}

func (o *GetBlockDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
