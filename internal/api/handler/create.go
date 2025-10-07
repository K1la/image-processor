package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/K1la/image-processor/internal/api/response"
	"github.com/K1la/image-processor/internal/model"
	"github.com/K1la/image-processor/internal/service"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"io"
)

func (h *Handler) CreateImage(c *ginext.Context) {
	imageHeader, err := c.FormFile("image")
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("invalid image")
		response.BadRequest(c.Writer, fmt.Errorf("invalid image: %w", err))
		return

	}

	file, err := imageHeader.Open()
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("could not open the image")
		response.BadRequest(c.Writer, fmt.Errorf("could not open the image: %w", err))
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("could not read the image")
		response.BadRequest(c.Writer, fmt.Errorf("could not read the image: %w", err))
		return
	}

	metadataStr := c.PostForm("metadata")
	var message model.Message
	if err = json.Unmarshal([]byte(metadataStr), &message); err != nil {
		zlog.Logger.Error().Err(err).Msg("could not unmarshal the metadata")
		response.BadRequest(c.Writer, fmt.Errorf("could not unmarshal the metadata: %w", err))
		return
	}

	id, err := h.service.CreateImage(c.Request.Context(), fileBytes, message)
	if err != nil {
		if errors.Is(err, service.ErrInvalidImageFormat) || errors.Is(err, service.ErrInvalidTask) {
			zlog.Logger.Error().Err(err).Msg("could not create the image")
			response.BadRequest(c.Writer, fmt.Errorf("could not create the image: %w", err))
		}
		zlog.Logger.Error().Err(err).Msg("could not create the image")
		response.Internal(c.Writer, fmt.Errorf("could not create the image: %w", err))
		return
	}

	zlog.Logger.Info().Interface("id", id).Msg("successfully created the image with id <-")
	response.Created(c.Writer, id)
}
