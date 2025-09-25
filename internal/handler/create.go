package handler

import (
	"net/http"

	"github.com/Komilov31/sales-tracker/internal/dto"
	validate "github.com/Komilov31/sales-tracker/internal/validator"
	"github.com/go-playground/validator/v10"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// CreateItem godoc
//
//	@Summary		Create a new item
//	@Description	Creates a new expense or income entry in the sales tracker
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.CreateItem	true	"Item to create"
//	@Success		200		{object}	dto.ItemWithoutAggregated	"Created item"
//	@Failure		400		{object}	map[string]string			"Invalid payload"
//	@Failure		500		{object}	map[string]string			"Internal server error"
//	@Router			/items [post]
func (h *Handler) CreateItem(c *ginext.Context) {
	var createItem dto.CreateItem
	if err := c.BindJSON(&createItem); err != nil {
		zlog.Logger.Error().Msg("could not unmarshal json: " + err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid payload"})
		return
	}

	if err := validate.Validator.Struct(createItem); err != nil {
		errors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid payload: " + errors.Error()})
		return
	}

	item, err := h.service.CreateItem(h.ctx, createItem)
	if err != nil {
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "could not create item"})
		return
	}

	zlog.Logger.Info().Msg("successfully handled POST request and created item")
	c.JSON(http.StatusOK, convertWithoutAggregated(item))
}
