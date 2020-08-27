package contexts

import (
	"avito-chat_service/internal/api/db"
	"avito-chat_service/internal/api/jsonapi"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Context ...
const (

	// Users ...
	Users = iota

	// Chats ...
	Chats

	// Messages ...
	Messages
)

const (

	// Add ...
	Add = iota

	// Create ...
	Create

	// Get ...
	Get

	// Send ...
	Send
)

// SyntaxJSONErr ...
type (
	//	InvalidJSONAPiFormatErr struct {}
	SyntaxJSONErr         struct{}
	UnknownFieldErr       struct{}
	InvalidFieldTypeErr   struct{}
	InvalidFieldFormatErr struct{}
	NotUniqueErr          struct{}
	EmptyFieldErr         struct{}
	NotFoundErr           struct{}
	InternalErr           struct{}
)

// Error ...
type Error struct {
	Type            interface{}
	UserDescription string
	Details         string
	Err             error
}

// Error ...
func (e *Error) Error() string {
	return e.Err.Error()
}

func getContext(context uint8) string {
	switch context {
	case Users:
		return "Users"
	case Chats:
		return "Chats"
	case Messages:
		return "Messages"
	default:
		return "undefined"
	}
}

func getModule(module uint8) string {
	switch module {
	case Add:
		return "Add"
	case Create:
		return "Create"
	case Get:
		return "Get"
	case Send:
		return "Send"
	default:
		return "undefined"
	}
}

// ErrorHandler ...
func ErrorHandler(context, module uint8, req *http.Request, err []error) (
	resp jsonapi.ResponseObject) {
	errStr := fmt.Sprintf("ERROR IN CONTEXT: %s/Delivery/%s by REQ: %v\n",
		getContext(context), getModule(module), req)

	for i := range err {
		newErr := &Error{}
		errors.As(err[i], &newErr)
		errStr += fmt.Sprintf(
			"TYPE: %s; DESCRIPTION: %s; DETAILS: %s; "+"RAW.ERR: %s\n",
			reflect.TypeOf(newErr.Type), newErr.UserDescription, newErr.Details, newErr.Err,
		)

		switch newErr.Type.(type) {
		// case InvalidJSONAPiFormatErr:
		//	resp.Errors = append(resp.Errors, jsonapi.ErrorObject{
		//		Status: http.StatusUnsupportedMediaType,
		//		Title:  "Invalid JSON API format",
		//		Detail: newErr.UserDescription,
		//	})
		case SyntaxJSONErr:
			resp.Errors = append(resp.Errors, jsonapi.ErrorObject{
				Status: http.StatusBadRequest,
				Title:  "Invalid JSON syntax",
				Detail: newErr.UserDescription,
			})
		case UnknownFieldErr:
			resp.Errors = append(resp.Errors, jsonapi.ErrorObject{
				Status: http.StatusBadRequest,
				Title:  "Unknown field",
				Detail: newErr.UserDescription,
			})
		case InvalidFieldTypeErr:
			resp.Errors = append(resp.Errors, jsonapi.ErrorObject{
				Status: http.StatusBadRequest,
				Title:  "Invalid field type",
				Detail: newErr.UserDescription,
			})
		case InvalidFieldFormatErr:
			resp.Errors = append(resp.Errors, jsonapi.ErrorObject{
				Status: http.StatusBadRequest,
				Title:  "Invalid field format",
				Detail: newErr.UserDescription,
			})
		case db.SyntaxErr:
			resp.Errors = append(resp.Errors, jsonapi.ErrorObject{
				Status: http.StatusBadRequest,
				Title:  "Invalid field format",
				Detail: "Something wrong with your request data",
			})
		case EmptyFieldErr:
			resp.Errors = append(resp.Errors, jsonapi.ErrorObject{
				Status: http.StatusBadRequest,
				Title:  "Empty field",
				Detail: newErr.UserDescription,
			})
		case NotUniqueErr:
			resp.Errors = append(resp.Errors, jsonapi.ErrorObject{
				Status: http.StatusConflict,
				Title:  "Object exists, must be unique",
				Detail: newErr.UserDescription,
			})
		case NotFoundErr:
			resp.Errors = append(resp.Errors, jsonapi.ErrorObject{
				Status: http.StatusNotFound,
				Title:  "Object not found",
				Detail: newErr.UserDescription,
			})
		default: // internal, db.InternalErr ...
			resp.Errors = append(resp.Errors, jsonapi.ErrorObject{
				Status: http.StatusInternalServerError,
				Title:  "Internal server error",
				Detail: "oops, something goes wrong :(",
			})
		}
	}

	if gin.Mode() == "debug" ||
		resp.Errors[0].Status == http.StatusInternalServerError {
		log.Println(errStr)
	}

	return resp
}

// ErrHandlerUseCase ...
func ErrHandlerUseCase(context, module uint8, req *http.Request, err error) error {
	newErr := &Error{}
	errors.As(err, &newErr)

	errStr := fmt.Sprintf("ERROR %s/UseCase/%s by REQ: %v\n",
		getContext(context), getModule(module), req)
	errStr += fmt.Sprintf("TYPE: %s; DESCRIPTION: %s; DETAILS: %s; RAW.ERR: %s\n",
		reflect.TypeOf(newErr.Type), newErr.UserDescription, newErr.Details, newErr.Err)

	if gin.Mode() == "debug" {
		log.Println(errStr)
	}

	switch newErr.Type.(type) {
	case db.InternalErr:
		newErr.Type = InternalErr{}
	case db.NotUniqueErr:
		newErr.Type = NotUniqueErr{}
	case db.ForeignKeyViolation, db.NotFound:
		newErr.Type = NotFoundErr{}
	default: // users.InternalErr, db.SyntaxErr
		return err
	}

	return newErr
}

// ParseErr ...
func ParseErr(err error) []error {
	var errs []error

	errCopy := err
	switch errCopy.(type) {
	case *json.UnmarshalTypeError:
		{
			errTyped := err.(*json.UnmarshalTypeError)
			errs = append(errs, &Error{
				Type: InvalidFieldTypeErr{},
				UserDescription: "Field '" + errTyped.Field + "' must be " +
					errTyped.Type.Name() + ", have " + errTyped.Value,
				Err: err,
			})
		}
	case *json.SyntaxError:
		{
			errTyped := err.(*json.SyntaxError)
			errs = append(errs, &Error{
				Type: SyntaxJSONErr{},
				UserDescription: "Check byte: " +
					strconv.FormatInt(errTyped.Offset, 10),
				Err: err,
			})
		}
	case validator.ValidationErrors:
		{
			ers := err.(validator.ValidationErrors)
			for i := range ers {
				errNew := &Error{
					Err: err,
				}
				switch ers[i].Tag() {
				case "uuid5":
					errNew.UserDescription = "Check format of '" +
						strings.ToLower(ers[i].Field()) + "'"
					errNew.Type = InvalidFieldFormatErr{}
				case "required":
					errNew.UserDescription = "Field '" +
						strings.ToLower(ers[i].Field()) +
						"' must be filled"
					errNew.Type = EmptyFieldErr{}
				case "unique":
					errNew.UserDescription = "Field '" + strings.ToLower(ers[i].Field()) +
						"' must have unique params"
					errNew.Type = NotUniqueErr{}
				case "min":
					errNew.UserDescription = "Field '" + strings.ToLower(ers[i].Field()) +
						"' must have at least " + ers[i].Param() + " args"
					errNew.Type = InvalidFieldFormatErr{}
					// case "jsonapi":
					//	errNew.Type = InvalidJSONAPiFormatErr{}
					//	errNew.UserDescription = "Can't find '" +
					//	strings.ToLower(ers[i].Field()) + "' field"
				}
				errs = append(errs, errNew)
			}
		}
	default:
		{
			errStr := err.Error()
			if strings.Contains(errStr, "json: unknown field ") {
				fields := strings.SplitAfter(errStr, "json: unknown field ")
				errStr = fields[1][1:]
				errStr = errStr[:len(errStr)-1]
				errs = append(errs, &Error{
					Type:            UnknownFieldErr{},
					UserDescription: "'" + errStr + "'",
					Err:             err,
				})
			} else {
				errs = append(errs, &Error{
					Type: InternalErr{},
					Err:  err,
				})
			}
		}
	}

	return errs
}
