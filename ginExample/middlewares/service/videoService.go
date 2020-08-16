package service

import "back-end-notes/ginExample/middlewares/entity"

type VideoService interface {
	Save(video entity.Video) entity.Video
	FindAll() []entity.Video
}

type videoService struct {
	videos []entity.Video
}

func New() VideoService {
	return &videoService{}
}

func (this *videoService) Save(video entity.Video) entity.Video {
	this.videos = append(this.videos, video)
	return video
}

func (this *videoService) FindAll() []entity.Video {
	return this.videos
}
