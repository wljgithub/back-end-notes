package controller

import (
	"back-end-notes/ginExample/middlewares/entity"
	"back-end-notes/ginExample/middlewares/service"
	"github.com/gin-gonic/gin"
)

type VideoController interface {
	FindAll() []entity.Video
	Save(ctx *gin.Context) entity.Video
}

type controller struct {
	service service.VideoService
}

func New(videoService service.VideoService) VideoController {
	return &controller{service: videoService}
}

func (this *controller) FindAll() []entity.Video {
	return this.service.FindAll()
}
func (this *controller) Save(ctx *gin.Context) entity.Video {
	var video entity.Video
	ctx.BindJSON(&video)
	return this.service.Save(video)
}
