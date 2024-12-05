package cart

import (
	"net/http"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> b0fe977 (place order and make payment)
=======
>>>>>>> b0fe977 (place order and make payment)
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
<<<<<<< HEAD
<<<<<<< HEAD
=======

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
>>>>>>> c449ffc (feat: cart api)
=======
>>>>>>> 21608b5 (cart and order api)
=======

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
>>>>>>> c449ffc (feat: cart api)
=======
>>>>>>> 21608b5 (cart and order api)
)

type CartControllerInterface interface {
	AddToCart(ctx *gin.Context)
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	GetCartItems(ctx *gin.Context)
	CreateOrder(c *gin.Context)
	GetOrders(c *gin.Context)
	GetOrderByID(c *gin.Context)
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	RemoveItemFromCart(c *gin.Context)
	GetAllOrders(c *gin.Context)
<<<<<<< HEAD
}

func (c *CartController) AddToCart(ctx *gin.Context) {
	var req CartItemRequest
=======
=======
	GetCartItems(ctx *gin.Context)
	CreateOrder(c *gin.Context)
>>>>>>> 21608b5 (cart and order api)
=======
>>>>>>> b0fe977 (place order and make payment)
=======
	RemoveItemFromCart(c *gin.Context)
>>>>>>> 4a16bfc (remove item in cart)
=======
>>>>>>> c449ffc (feat: cart api)
=======
	GetCartItems(ctx *gin.Context)
	CreateOrder(c *gin.Context)
>>>>>>> 21608b5 (cart and order api)
=======
>>>>>>> b0fe977 (place order and make payment)
=======
	RemoveItemFromCart(c *gin.Context)
>>>>>>> 4a16bfc (remove item in cart)
}

func (c *CartController) AddToCart(ctx *gin.Context) {
	var req CartItem
<<<<<<< HEAD
>>>>>>> c449ffc (feat: cart api)
=======
}

func (c *CartController) AddToCart(ctx *gin.Context) {
	var req CartItemRequest
>>>>>>> dc47646 (Optimize SQL query)
=======
>>>>>>> c449ffc (feat: cart api)
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	res, err := c.service.AddToCartService(ctx, req, authPayload.Username)
=======
	err = c.service.AddToCartService(ctx, req, authPayload.Username)
>>>>>>> c449ffc (feat: cart api)
=======
	res, err := c.service.AddToCartService(ctx, req, authPayload.Username)
>>>>>>> 21608b5 (cart and order api)
=======
	err = c.service.AddToCartService(ctx, req, authPayload.Username)
>>>>>>> c449ffc (feat: cart api)
=======
	res, err := c.service.AddToCartService(ctx, req, authPayload.Username)
>>>>>>> 21608b5 (cart and order api)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 21608b5 (cart and order api)
=======
>>>>>>> 21608b5 (cart and order api)
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

<<<<<<< HEAD
<<<<<<< HEAD
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
=======
	ctx.JSON(http.StatusOK, gin.H{"message": "Item added to cart successfully"})
>>>>>>> c449ffc (feat: cart api)
=======
>>>>>>> 21608b5 (cart and order api)
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
=======
	ctx.JSON(http.StatusOK, gin.H{"message": "Item added to cart successfully"})
>>>>>>> c449ffc (feat: cart api)
=======
>>>>>>> 21608b5 (cart and order api)
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
