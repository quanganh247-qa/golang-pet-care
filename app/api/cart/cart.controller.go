package cart

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type CartControllerInterface interface {
	AddToCart(ctx *gin.Context)
	GetCartItems(ctx *gin.Context)
	CreateOrder(c *gin.Context)
	GetOrders(c *gin.Context)
	GetOrderByID(c *gin.Context)
	RemoveItemFromCart(c *gin.Context)
	GetAllOrders(c *gin.Context)
	GetOrderHistory(c *gin.Context)
	DeleteCart(c *gin.Context)
	// IncreaseItemQuantity(c *gin.Context)
	// DecreaseItemQuantity(c *gin.Context)
	// UpdateItemQuantity(c *gin.Context)
}

func (c *CartController) AddToCart(ctx *gin.Context) {
	var req CartItemRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := c.service.AddToCartService(ctx, req, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Item added to cart successfully", res))
}

func (c *CartController) GetCartItems(ctx *gin.Context) {

	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := c.service.GetCartItemsService(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Item added to cart successfully", res))
}

func (c *CartController) CreateOrder(ctx *gin.Context) {

	var req PlaceOrderRequest

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := c.service.CreateOrderService(ctx, authPayload.Username, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Order is added successfully", res))

}

func (c *CartController) GetOrders(ctx *gin.Context) {
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := c.service.GetOrdersService(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Orders fetched successfully", res))
}

func (c *CartController) GetOrderByID(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	// Convert orderID to int64
	id, _ := strconv.ParseInt(orderID, 10, 64)
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := c.service.GetOrderByIdService(ctx, authPayload.Username, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Order fetched successfully", res))
}

func (c *CartController) RemoveItemFromCart(ctx *gin.Context) {
	id := ctx.Param("product_id")

	// Convert orderID to int64
	idInt, _ := strconv.ParseInt(id, 10, 64)

	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = c.service.DeleteItemFromCartService(ctx, authPayload.Username, idInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Item removed from cart successfully", nil))
}

func (c *CartController) GetAllOrders(ctx *gin.Context) {

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := c.service.GetAllOrdersService(ctx, pagination, "paid")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Orders fetched successfully", res))
}

func (c *CartController) GetOrderHistory(ctx *gin.Context) {
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := c.service.GetOrderHistoryService(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Order history fetched successfully", res))
}

// func (c *CartController) IncreaseItemQuantity(ctx *gin.Context) {
// 	var req UpdateItemQuantityRequest
// 	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	authPayload, err := middleware.GetAuthorizationPayload(ctx)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	err = c.service.IncreaseItemQuantityService(ctx, authPayload.Username, req.ItemID, req.Quantity)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, util.SuccessResponse("Item quantity increased successfully", nil))
// }

// func (c *CartController) DecreaseItemQuantity(ctx *gin.Context) {
// 	var req UpdateItemQuantityRequest
// 	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	authPayload, err := middleware.GetAuthorizationPayload(ctx)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	err = c.service.DecreaseItemQuantityService(ctx, authPayload.Username, req.ItemID, req.Quantity)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, util.SuccessResponse("Item quantity decreased successfully", nil))
// }

// func (c *CartController) UpdateItemQuantity(ctx *gin.Context) {
// 	var req UpdateItemQuantityRequest
// 	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	authPayload, err := middleware.GetAuthorizationPayload(ctx)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	err = c.service.UpdateItemQuantityService(ctx, authPayload.Username, req.ItemID, req.Quantity)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, util.SuccessResponse("Item quantity updated successfully", nil))
// }

func (c *CartController) DeleteCart(ctx *gin.Context) {
	cartID := ctx.Param("cart_id")
	idInt, _ := strconv.ParseInt(cartID, 10, 64)
	err := c.service.DeleteCartService(ctx, idInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Cart deleted successfully", nil))
}
