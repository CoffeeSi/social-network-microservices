package handler

import (
	"net/http"

	authpb "github.com/IsFariza/maxat-protobuf/auth-service-pb"
	"github.com/gin-gonic/gin"
)

func (g *Gateway) Register(c *gin.Context) {
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		DOB       string `json:"dob"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	if !bindJSON(c, &req) {
		return
	}
	resp, err := g.auth.RegisterUser(c.Request.Context(), &authpb.RegisterUserRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Dob:       req.DOB,
		Email:     req.Email,
		Password:  req.Password,
	})
	writeGRPC(c, resp, err, http.StatusCreated)
}

func (g *Gateway) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if !bindJSON(c, &req) {
		return
	}
	resp, err := g.auth.LoginUser(c.Request.Context(), &authpb.LoginUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if !bindJSON(c, &req) {
		return
	}
	resp, err := g.auth.RefreshToken(c.Request.Context(), &authpb.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) Verify(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
		Token string `json:"token"`
		Code  string `json:"code"`
	}
	if !bindJSON(c, &req) {
		return
	}
	// Support both 'token' and 'code' field names
	code := req.Code
	if code == "" {
		code = req.Token
	}
	resp, err := g.auth.VerifyUser(c.Request.Context(), &authpb.VerifyUserRequest{
		Email: req.Email,
		Code:  code,
	})
	writeGRPC(c, resp, err, http.StatusOK)
}
