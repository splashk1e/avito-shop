package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/splashk1e/avito-shop/internal/services"
)

type Handler struct {
	authservice        *services.AuthService
	transactionService *services.TransactionService
}

func NewHandler(authservice *services.AuthService, transactionService *services.TransactionService) *Handler {
	return &Handler{
		authservice:        authservice,
		transactionService: transactionService,
	}
}
func (handler *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.POST("/api/auth", handler.Auth)
	api := router.Group("/api", handler.userIndentity)
	{

		api.GET("/buy/:item", handler.BuyItem)
		api.POST("/sendCoin", handler.SendCoin)
		api.GET("/info", handler.Info)
	}

	return router
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type responseToken struct {
	token string `json:"token"`
}

func (h *Handler) Auth(ctx *gin.Context) {
	var user User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	token, err := h.authservice.Login(ctx, user.Username, user.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		}
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	response := responseToken{token: token}
	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) BuyItem(ctx *gin.Context) {
	item := ctx.Param("item")
	username, ok := ctx.Keys["username"].(string)
	if !ok {
		newErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
	}
	if _, err := h.transactionService.SaveTransaction(ctx, username, "", item, 0); err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		}
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	ctx.AbortWithStatus(http.StatusOK)
}

type transacRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

func (h *Handler) SendCoin(ctx *gin.Context) {
	var req transacRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	username, ok := ctx.Keys["username"].(string)
	if !ok {
		newErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
	}
	_, err := h.transactionService.SaveTransaction(ctx, username, req.ToUser, "", req.Amount)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		}
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	ctx.AbortWithStatus(http.StatusOK)
}

type Recevied struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}
type Sent struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type CoinHistory struct {
	Recevied []Recevied `json:"recieved"`
	Sent     []Sent     `json:"sent"`
}
type infoResponse struct {
	Coins       int         `json:"coins"`
	Inventory   []string    `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

func (h *Handler) Info(ctx *gin.Context) {
	var infoResponse infoResponse
	sent := make([]Sent, 0)
	recieved := make([]Recevied, 0)
	username, ok := ctx.Keys["username"].(string)
	if !ok {
		newErrorResponse(ctx, http.StatusUnauthorized, "unauthorized")
	}
	transactions, err := h.transactionService.GetTransactions(ctx, username)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		}
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	coins, err := h.authservice.GetCoinsInfo(ctx, username)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		}
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	items, err := h.transactionService.GetItems(ctx, username)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		}
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	infoResponse.Inventory = items
	for _, val := range transactions {
		if val.Sender == username {
			sent = append(sent, Sent{Amount: val.Amount, ToUser: val.Reciever})
		}
		if val.Reciever == username {
			recieved = append(recieved, Recevied{Amount: val.Amount, FromUser: val.Sender})
		}
	}
	infoResponse.CoinHistory.Recevied = recieved
	infoResponse.CoinHistory.Sent = sent
	infoResponse.Coins = coins
	ctx.JSON(http.StatusOK, infoResponse)
}
