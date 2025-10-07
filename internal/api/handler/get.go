package handler

import (
	"errors"
	"fmt"
	"github.com/K1la/image-processor/internal/api/response"
	"github.com/K1la/image-processor/internal/service"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"os"
)

func (h *Handler) GetImageById(c *ginext.Context) {
	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("could not parse uuid from id")
		response.BadRequest(c.Writer, fmt.Errorf("could not parse uuid from id: %w", err))
		return
	}

	image, err := h.service.GetImageById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotProcessed) {
			response.OK(c.Writer, "in processing, wait please")
			return
		}

		if errors.Is(err, service.ErrNoSuchImage) {
			response.BadRequest(c.Writer, err)
			return
		}

		zlog.Logger.Error().Err(err).Msg("could not get image")
		response.Internal(c.Writer, fmt.Errorf("could not get image: %w", err))
		return
	}

	zlog.Logger.Info().Msg("successfully handled GET request and returned image to user")
	c.File(image)

	if err = os.Remove(image); err != nil {
		zlog.Logger.Error().Msg("could not delete processed image from local storage after sending to user: " + err.Error())
	}
}
