package controller

import (
	"back-end-notes/ginExample/dataBindingAndValidation/entity"
	"back-end-notes/ginExample/dataBindingAndValidation/service"
	"back-end-notes/ginExample/dataBindingAndValidation/validators"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type VideoController interface {
	FindAll() []entity.Video
	Save(ctx *gin.Context) error
}

type controller struct {
	service  service.VideoService
	validate *validator.Validate
}

func New(videoService service.VideoService) VideoController {
	validate := validator.New()
	validate.RegisterValidation("is-cool", validators.ValidateCoolTitle)
	return &controller{service: videoService, validate: validate}
}

func (this *controller) FindAll() []entity.Video {
	return this.service.FindAll()
}
func (this *controller) Save(ctx *gin.Context) error {
	var video entity.Video
	err := ctx.ShouldBindJSON(&video)
	if err != nil {
		return err
	}
	err = this.validate.Struct(video)
	if err != nil {
		return err
	}
	this.service.Save(video)
	return nil
	//return this.service.Save(video)
}
