package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type errorResponse struct {
	Error string `json:"error"`
}

func bindJSON(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid json: " + err.Error()})
		return false
	}
	return true
}

func writeGRPC(c *gin.Context, data any, err error, successStatus int) {
	if err != nil {
		c.JSON(grpcHTTPStatus(err), errorResponse{Error: status.Convert(err).Message()})
		return
	}
	c.JSON(successStatus, data)
}

func writeEmptyGRPC(c *gin.Context, err error) {
	if err != nil {
		c.JSON(grpcHTTPStatus(err), errorResponse{Error: status.Convert(err).Message()})
		return
	}
	c.Status(http.StatusNoContent)
}

func writeError(c *gin.Context, status int, msg string) {
	c.JSON(status, errorResponse{Error: msg})
}

func grpcHTTPStatus(err error) int {
	switch status.Code(err) {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.Unavailable:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}
