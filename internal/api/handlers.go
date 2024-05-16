package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/online_order/internal/models"
	"github.com/tanush-128/openzo_backend/online_order/internal/service"
)

type Handler struct {
	online_orderService service.OnlineOrderService
}

func NewHandler(online_orderService *service.OnlineOrderService) *Handler {
	return &Handler{online_orderService: *online_orderService}
}

func (h *Handler) CreateOnlineOrder(ctx *gin.Context) {
	var online_order models.OnlineOrder
	if err := ctx.BindJSON(&online_order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdOnlineOrder, err := h.online_orderService.CreateOnlineOrder(ctx, online_order)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdOnlineOrder)

}

func (h *Handler) GetOnlineOrderByID(ctx *gin.Context) {
	id := ctx.Param("id")

	online_order, err := h.online_orderService.GetOnlineOrderByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, online_order)
}



func (h *Handler) GetOnlineOrdersByStoreID(ctx *gin.Context) {
	store_id := ctx.Param("store_id")

	online_orders, err := h.online_orderService.GetOnlineOrdersByStoreID(ctx, store_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, online_orders)
}

func (h *Handler) ChangeOrderStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	status := ctx.Query("status")

	updatedOnlineOrder, err := h.online_orderService.ChangeOrderStatus(ctx, id, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedOnlineOrder)
}

func (h *Handler) UpdateOnlineOrder(ctx *gin.Context) {
	var online_order models.OnlineOrder
	if err := ctx.BindJSON(&online_order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedOnlineOrder, err := h.online_orderService.UpdateOnlineOrder(ctx, online_order)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedOnlineOrder)
}
