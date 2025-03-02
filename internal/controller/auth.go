package controller

import (
	"TaskList/internal/lib/http/response"
	"TaskList/internal/models"
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strings"
)

type Auth interface {
	Registration(ctx context.Context, email string, password string) (int64, error)
	Login(ctx context.Context, email string, password []byte) (string, error)
}

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	response.Response
	Token string `json:"token,omitempty"`
}

type RegisterResponse struct {
	response.Response
	ID int64 `json:"id,omitempty"`
}

// Login ...
func (c Controller) Login(w http.ResponseWriter, r *http.Request) {
	const op = "controller.Login"
	a := AuthRequest{}

	if err := render.DecodeJSON(r.Body, &a); err != nil {
		c.log.Error(
			"failed decode json",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, RegisterResponse{
			Response: response.Error("failed decode json"),
		})
		return
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.log.Warn(
				"failed close body",
				slog.String("op", op),
			)
		}
	}()

	if err := validateAuthRequest(a); err != nil {
		c.log.Warn(
			"invalid request",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, LoginResponse{
			Response: response.Error(err.Error()),
		})
		return
	}

	c.log.Info(
		"attempting login",
		slog.String("op", op),
		slog.String("email", a.Email),
	)

	token, err := c.auth.Login(context.Background(), a.Email, []byte(a.Password))
	if err != nil {

		if errors.Is(err, models.ErrUserNotFound) {
			c.log.Warn(
				"user not found or invalid cred",
				slog.String("op", op),
				slog.String("email", a.Email),
			)

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, LoginResponse{
				Response: response.Error("user not found"),
			})
			return
		}

		c.log.Error(
			"failed login user",
			slog.String("op", op),
			slog.String("email", a.Email),
			slog.String("err", err.Error()),
		)

		resp := &LoginResponse{
			Response: response.Error("bad request"),
		}
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp)
		return
	}

	c.log.Info(
		"success login",
		slog.String("op", op),
		slog.String("email", a.Email),
	)

	resp := &LoginResponse{
		Response: response.OK(),
		Token:    token,
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)

}

// Registration new user
func (c Controller) Registration(w http.ResponseWriter, r *http.Request) {
	const op = "controller.Registration"
	a := &AuthRequest{}

	if err := render.DecodeJSON(r.Body, a); err != nil {
		c.log.Error(
			"failed decode json",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, RegisterResponse{
			Response: response.Error("failed decode json"),
		})
		return
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			c.log.Warn(
				"failed close body",
				slog.String("op", op),
			)
		}
	}()

	if err := validateAuthRequest(a); err != nil {
		c.log.Warn(
			"invalid request",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, LoginResponse{
			Response: response.Error(err.Error()),
		})
		return
	}

	c.log.Info(
		"attempt registration new user",
		slog.String("op", op),
		slog.String("email", a.Email),
	)

	uid, err := c.auth.Registration(context.Background(), a.Email, a.Password)
	if err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			c.log.Warn(
				"user already exists",
				slog.String("op", op),
				slog.String("email", a.Email),
			)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, RegisterResponse{
				Response: response.Error("user already exists"),
			})

			return
		}

		c.log.Error(
			err.Error(),
			slog.String("op", op),
			slog.String("email", a.Email),
		)

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, RegisterResponse{
			Response: response.Error("internal error"),
		})

		return
	}

	c.log.Info(
		"success registration user",
		slog.String("op", op),
		slog.String("email", a.Email),
	)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, RegisterResponse{
		Response: response.OK(),
		ID:       uid,
	})
}

func validateAuthRequest(a interface{}) error {
	var errMsgs []string
	validate := validator.New()
	if err := validate.Struct(a); err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				switch e.ActualTag() {
				case "required":
					errMsgs = append(
						errMsgs,
						fmt.Sprintf("field %s is a required field", e.Field()),
					)
				case "min":
					errMsgs = append(
						errMsgs,
						fmt.Sprintf("field %s must consist of at least 8 characters", e.Field()),
					)
				default:
					errMsgs = append(
						errMsgs,
						fmt.Sprintf("field %s is not valid", e.Field()),
					)
				}
			}
		}
	}

	if len(errMsgs) > 0 {
		return errors.New(strings.Join(errMsgs, ", "))
	}

	return nil
}
