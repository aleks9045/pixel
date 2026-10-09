package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/aleks9045/pixels/internal/model"
	"github.com/aleks9045/pixels/internal/service/auth"
	"github.com/aleks9045/pixels/internal/service/pixel"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type WsService struct {
	mu           *sync.RWMutex
	userRegistry map[string]*websocket.Conn
	pixelService *pixel.PixelService
}

func NewWsService(pixelService *pixel.PixelService) *WsService {
	return &WsService{
		mu:           &sync.RWMutex{},
		userRegistry: make(map[string]*websocket.Conn),
		pixelService: pixelService,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *WsService) NewConnection(c *gin.Context) {
	w, r := c.Writer, c.Request
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error during protocol upgrade:"+err.Error())
		return
	}
	userId, ok := auth.GetUserIdFromCtx(c)
	if !ok {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
	s.addUserToRegistry(userId, ws)
	defer func() {
		_ = ws.Close()
		s.removeUserFromRegistry(userId)
	}()
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if mt != websocket.BinaryMessage {
			break
		}
		wsMessage := &model.Ws{}

		proto.Unmarshal(message, wsMessage)
		switch payload := wsMessage.Payload.(type) {

		case *model.Ws_Pixel:
			log.Println("pixel update:", payload.Pixel)

			err := s.pixelService.PutPixel(payload.Pixel, userId)
			if err != nil {
				s.sendPixelError(ws, payload.Pixel, err.Error())
			}
			binaryMessage, err := proto.Marshal(payload.Pixel)
			if err != nil {
				log.Fatal("Marshal error:", err.Error())
				return
			}
			s.broadcastMessage(userId, websocket.BinaryMessage, binaryMessage)

		case *model.Ws_PixelRegion:
			log.Println("pixel region update:", payload.PixelRegion)

			err := s.pixelService.PutPixelRegion(payload.PixelRegion, userId)
			if err != nil {
				log.Println(err)
				s.sendPixelRegionError(ws, payload.PixelRegion, err.Error())
			}

			binaryMessage, err := proto.Marshal(payload.PixelRegion)
			if err != nil {
				log.Fatal("Marshal error:", err)
			}

			s.broadcastMessage(userId, websocket.BinaryMessage, binaryMessage)

			if err != nil {
				break
			}
		}
		err = ws.WriteMessage(mt, message)
	}
}

func (s *WsService) sendPixelError(ws *websocket.Conn, wsMessage *model.WsPixel, errorMessage string) {

	wsError := &model.Ws{
		Payload: &model.Ws_PixelError{
			PixelError: &model.WsPixelError{
				X:       wsMessage.X,
				Y:       wsMessage.Y,
				Message: errorMessage,
			},
		},
	}

	data, err := proto.Marshal(wsError)
	if err != nil {
		log.Fatal("Marshal error:", err)
		return
	}

	ws.WriteMessage(websocket.BinaryMessage, data)
}

func (s *WsService) sendPixelRegionError(ws *websocket.Conn, wsMessage *model.WsPixelRegion, errorMessage string) {

	wsError := &model.Ws{
		Payload: &model.Ws_PixelRegionError{
			PixelRegionError: &model.WsPixelRegionError{
				X1:      wsMessage.X1,
				X2:      wsMessage.X2,
				Y1:      wsMessage.Y1,
				Y2:      wsMessage.Y2,
				Message: errorMessage,
			},
		},
	}

	data, err := proto.Marshal(wsError)
	if err != nil {
		log.Fatal("Marshal error:", err)
	}

	ws.WriteMessage(websocket.BinaryMessage, data)
}

func (s *WsService) broadcastMessage(userId string, messageType int, data []byte) {
	s.mu.Lock()
	for key := range s.userRegistry {
		if key == userId {
			continue
		}
		err := s.userRegistry[key].WriteMessage(messageType, data)
		if err != nil {
			continue
		}
	}
	s.mu.Unlock()
}

func (s *WsService) addUserToRegistry(userId string, ws *websocket.Conn) {
	s.mu.Lock()
	s.userRegistry[userId] = ws
	s.mu.Unlock()
}

func (s *WsService) removeUserFromRegistry(userId string) {
	s.mu.Lock()
	delete(s.userRegistry, userId)
	s.mu.Unlock()
}
