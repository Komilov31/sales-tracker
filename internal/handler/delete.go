package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Komilov31/sales-tracker/internal/repository"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// DeleteItem godoc
//
//	@Summary		Delete an item
//	@Description	Remove an item by its ID
//	@Tags			items
//	@Produce		json
//	@Param			id	path		int		true	"Item ID"
//	@Success		200	{object}	map[string]string	"Success message"
//	@Failure		400	{object}	map[string]string	"Item not found or invalid ID"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/items/{id} [delete]
func (h *Handler) DeleteItem(c *ginext.Context) {
	itemID := c.Param("id")
	id, err := strconv.Atoi(itemID)
	if err != nil {
		zlog.Logger.Error().Msg("invalid id: " + err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteItem(h.ctx, id); err != nil {
		if errors.Is(err, repository.ErrNoSuchItem) {
			zlog.Logger.Error().Msg(err.Error())
			c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
			return
		}
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Msg("successfully handled DELETE request and deleted item")
	c.JSON(http.StatusOK, ginext.H{"status": "successfully delete item"})
}
