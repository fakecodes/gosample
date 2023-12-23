package http

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/fakecodes/gosample/domain"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// TaskHandler represent the httphandler for task list
type TaskHandler struct {
	TUsecase domain.TaskUsecase
}

// NewRoleHandler will initialize the role/ resources endpoint
func NewTaskHandler(e *echo.Echo, us domain.TaskUsecase) {
	handler := &TaskHandler{
		TUsecase: us,
	}
	e.POST("/task", handler.Create)
	// e.GET("/task/:id", handler.GetByID)
}

func isRequestValid(m *domain.Task) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Create will create the task by given request body
func (a *TaskHandler) Create(c echo.Context) (err error) {
	var task domain.Task
	err = c.Bind(&task)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&task); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = a.TUsecase.Create(ctx, &task)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, task)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
