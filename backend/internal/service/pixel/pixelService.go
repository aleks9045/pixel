package pixel

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/aleks9045/pixels/internal/model"
	"github.com/aleks9045/pixels/internal/repository"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

type PixelService struct {
	pixelRepository *repository.PixelRepository
}

func NewPixelService(repository *repository.PixelRepository) *PixelService {
	return &PixelService{
		pixelRepository: repository,
	}
}

func (s *PixelService) GetField(c *gin.Context) {
	pixelField, err := s.pixelRepository.GetField()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	data, err := proto.Marshal(&pixelField)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Data(http.StatusOK, "application/x-protobuf", data)
}

func (s *PixelService) PutPixel(wsMessagePixel *model.WsPixel, userId string) error {
	return s.pixelRepository.PutPixel(wsMessagePixel, userId)
}

func (s *PixelService) GetAuthor(c *gin.Context) {
	x, err := strconv.Atoi(c.Query("x"))
	if err != nil {
		c.String(http.StatusBadRequest, "param x is invalid")
		return
	}
	y, err := strconv.Atoi(c.Query("y"))
	if err != nil {
		c.String(http.StatusBadRequest, "param y is invalid")
		return
	}

	user, err := s.pixelRepository.GetPixelAuthor(x, y)
	if errors.Is(err, sql.ErrNoRows) {
		c.Status(http.StatusNoContent)
		return
	}
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	data, err := proto.Marshal(user)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Data(http.StatusOK, "application/x-protobuf", data)
}

func (s *PixelService) PutPixelRegion(wsMessagePixel *model.WsPixelRegion, userId string) error {
	return s.pixelRepository.PutPixelRegion(wsMessagePixel, userId)
}
