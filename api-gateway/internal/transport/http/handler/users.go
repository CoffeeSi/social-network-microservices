package handler

import (
	"net/http"

	userpb "github.com/IsFariza/maxat-protobuf/user-service-pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (g *Gateway) CreateUser(c *gin.Context) {
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		DOB       string `json:"dob"`
	}
	if !bindJSON(c, &req) {
		return
	}
	dob, err := parseDate(req.DOB)
	if err != nil {
		writeError(c, http.StatusBadRequest, "dob must be RFC3339 or YYYY-MM-DD")
		return
	}
	resp, err := g.users.CreateUser(c.Request.Context(), &userpb.CreateUserRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
		Dob:       timestamppb.New(dob),
	})
	writeGRPC(c, resp, err, http.StatusCreated)
}

func (g *Gateway) ListUsers(c *gin.Context) {
	resp, err := g.users.ListUsers(c.Request.Context(), &userpb.ListUsersRequest{
		PageSize: int32(queryInt(c, "page_size", 20)),
		Page:     int32(queryInt(c, "page", 1)),
	})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) GetUser(c *gin.Context) {
	resp, err := g.users.GetUser(c.Request.Context(), &userpb.GetUserRequest{Id: c.Param("id")})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) UpdateUser(c *gin.Context) {
	var req struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		DOB       *string `json:"dob"`
	}
	if !bindJSON(c, &req) {
		return
	}

	data := &userpb.UserUpdateData{}
	paths := make([]string, 0, 4)
	if req.FirstName != nil {
		data.FirstName = *req.FirstName
		paths = append(paths, "first_name")
	}
	if req.LastName != nil {
		data.LastName = *req.LastName
		paths = append(paths, "last_name")
	}
	if req.Email != nil {
		data.Email = *req.Email
		paths = append(paths, "email")
	}
	if req.DOB != nil {
		dob, err := parseDate(*req.DOB)
		if err != nil {
			writeError(c, http.StatusBadRequest, "dob must be RFC3339 or YYYY-MM-DD")
			return
		}
		data.Dob = timestamppb.New(dob)
		paths = append(paths, "dob")
	}
	resp, err := g.users.UpdateUser(c.Request.Context(), &userpb.UpdateUserRequest{
		Id:         c.Param("id"),
		Data:       data,
		UpdateMask: &fieldmaskpb.FieldMask{Paths: paths},
	})
	writeGRPC(c, resp, err, http.StatusOK)
}

func (g *Gateway) DeleteUser(c *gin.Context) {
	_, err := g.users.DeleteUser(c.Request.Context(), &userpb.DeleteUserRequest{Id: c.Param("id")})
	writeEmptyGRPC(c, err)
}

func (g *Gateway) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if !bindJSON(c, &req) {
		return
	}
	_, err := g.users.ChangePassword(c.Request.Context(), &userpb.ChangePasswordRequest{
		Id:          c.Param("id"),
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	writeEmptyGRPC(c, err)
}
