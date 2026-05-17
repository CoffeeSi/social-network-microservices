package handler

import (
	"net/http"

	authpb "github.com/CoffeeSi/social-network-microservices/api-gateway/proto/auth"
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
