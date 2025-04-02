package invoice

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

func NewInvoiceController(store db.Store) *InvoiceController {
	return &InvoiceController{
		service: NewInvoiceService(store),
	}
}

func (c *InvoiceController) CreateInvoice(ctx *gin.Context) {
	var req CreateInvoiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invoice, err := c.service.CreateInvoice(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, invoice)
}

func (c *InvoiceController) GetInvoiceByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	invoice, err := c.service.GetInvoiceByID(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, invoice)
}

func (c *InvoiceController) GetInvoiceByNumber(ctx *gin.Context) {
	invoiceNumber := ctx.Param("number")

	invoice, err := c.service.GetInvoiceByNumber(ctx, invoiceNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, invoice)
}

func (c *InvoiceController) ListInvoices(ctx *gin.Context) {

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	invoices, err := c.service.ListInvoices(ctx, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, invoices)
}

func (c *InvoiceController) UpdateInvoiceStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	var req UpdateInvoiceStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.service.UpdateInvoiceStatus(ctx, int32(id), req.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Invoice status updated successfully"})
}

func (c *InvoiceController) DeleteInvoice(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	err = c.service.DeleteInvoice(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Invoice deleted successfully"})
}

func (c *InvoiceController) AddInvoiceItem(ctx *gin.Context) {
	idStr := ctx.Param("id")
	invoiceID, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	var req InvoiceItem
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := c.service.AddInvoiceItem(ctx, int32(invoiceID), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, item)
}

func (c *InvoiceController) GetInvoiceItems(ctx *gin.Context) {
	idStr := ctx.Param("id")
	invoiceID, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	items, err := c.service.GetInvoiceItems(ctx, int32(invoiceID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

func (c *InvoiceController) UpdateInvoiceItem(ctx *gin.Context) {
	itemIDStr := ctx.Param("item_id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req UpdateInvoiceItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := c.service.UpdateInvoiceItem(ctx, int32(itemID), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (c *InvoiceController) DeleteInvoiceItem(ctx *gin.Context) {
	itemIDStr := ctx.Param("item_id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	err = c.service.DeleteInvoiceItem(ctx, int32(itemID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Invoice item deleted successfully"})
}

func (c *InvoiceController) CreateInvoiceFromTestOrder(ctx *gin.Context) {

	id := ctx.Param("test_order_id")

	newID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test order ID"})
		return
	}
	invoice, err := c.service.CreateInvoiceFromTestOrder(ctx, int32(newID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, invoice)
}
