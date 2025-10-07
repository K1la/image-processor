package handler

import (
	"fmt"

	"github.com/K1la/image-processor/internal/api/response"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func (h *Handler) DeleteImageByID(c *ginext.Context) {
	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("could not parse id to uuid")
		response.BadRequest(c.Writer, fmt.Errorf("invalid id was provided: %w", err))
		return
	}

	if err = h.service.DeleteImage(c.Request.Context(), id); err != nil {
		zlog.Logger.Error().Err(err).Msg("could not delete image")
		response.Internal(c.Writer, fmt.Errorf("could not delete image: %w", err))
		return
	}

	zlog.Logger.Info().Msg("successfully handled DELETE request and deleted image")
	response.OK(c.Writer, "successfully deleted image")
}
