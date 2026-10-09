package user

import (
	"io"
	"net/http"
	"strconv"

	"google.golang.org/protobuf/proto"

	"github.com/aleks9045/pixels/internal/model"
	"github.com/aleks9045/pixels/internal/repository"
	"github.com/aleks9045/pixels/internal/service/auth"
	"github.com/gin-gonic/gin"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetAll(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		c.String(http.StatusBadRequest, "query param page is invalid")
		return
	}
	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil {
		c.String(http.StatusBadRequest, "query param page is invalid")
		return
	}

	users, err := s.userRepository.GetUsers(page, size)
	if err != nil {
		c.String(http.StatusInternalServerError, "error")
		return
	}
	protoUsers := &model.Users{User: users}
	data, err := proto.Marshal(protoUsers)
	if err != nil {
		c.String(http.StatusInternalServerError, "error")
		return
	}
	c.Data(http.StatusOK, "application/x-protobuf", data)

}

func (s *UserService) GetMe(c *gin.Context) {
	userId, ok := auth.GetUserIdFromCtx(c)
	if !ok {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := s.userRepository.GetByUserId(userId)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	data, err := proto.Marshal(user)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Data(http.StatusOK, "application/x-protobuf", data)

}

func (s *UserService) PatchUser(c *gin.Context) {
	userId := c.Param("user_id")
	authUserId, ok := auth.GetUserIdFromCtx(c)
	if !ok {
		c.String(http.StatusUnauthorized, "unauthorized")
		return
	}
	isAdmin, ok := auth.GetIsAdminFromCtx(c)
	if !ok {
		c.String(http.StatusUnauthorized, "unauthorized")
		return
	}
	if userId != authUserId || !isAdmin {
		c.String(http.StatusForbidden, "forbidden")
		return
	}

	bodyBinary, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	patchedUserModel := &model.User{}
	err = proto.Unmarshal(bodyBinary, patchedUserModel)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid body")
		return
	}
	if patchedUserModel.TimeLimit < 0 {
		c.String(http.StatusBadRequest, "invalid body")
		return
	}

	userModel, err := s.userRepository.GetByUserId(userId)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if userModel.Name != patchedUserModel.Name {
		userModel.Name = patchedUserModel.Name
	}
	// Only for admins - (isAdmin && userModel.TimeLimit != patchedUserModel.TimeLimit)
	if userModel.TimeLimit != patchedUserModel.TimeLimit {
		userModel.TimeLimit = patchedUserModel.TimeLimit
	}
	if userModel.StudentGroup != patchedUserModel.StudentGroup {
		userModel.StudentGroup = patchedUserModel.StudentGroup
	}

	err = s.userRepository.UpdateUser(userModel)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "updated")
}

func (s *UserService) GetUser(c *gin.Context) {
	userId := c.Param("user_id")

	user, err := s.userRepository.GetByUserId(userId)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	data, err := proto.Marshal(user)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Data(http.StatusOK, "application/x-protobuf", data)

}
