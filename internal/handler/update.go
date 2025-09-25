package handler

import (
	"net/http"
	"strconv"

	"github.com/Komilov31/sales-tracker/internal/dto"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// UpdateItem godoc
//
//	@Summary		Update an item
//	@Description	Partially update an existing item by ID
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int		true	"Item ID"
//	@Param			body	body		dto.UpdateItem	true	"Fields to update"
//	@Success		200		{object}	map[string]string	"Success message"
//	@Failure		400		{object}	map[string]string	"Invalid ID or payload"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/items/{id} [put]
func (h *Handler) UpdateItem(c *ginext.Context) {
	itemID := c.Param("id")
	id, err := strconv.Atoi(itemID)
	if err != nil {
		zlog.Logger.Error().Msg("invalid id: " + err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid id"})
		return
	}

	var updateItem dto.UpdateItem
	if err := c.BindJSON(&updateItem); err != nil {
		zlog.Logger.Error().Msg("could not unmarshal json: " + err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid payload"})
		return
	}

	if err := h.service.UpdateItem(h.ctx, id, updateItem); err != nil {
		zlog.Logger.Error().Msg("could not update item: " + err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Msg("successfully handled PUT request and updated item")
	c.JSON(http.StatusOK, ginext.H{"status": "successfully updated item"})
}
