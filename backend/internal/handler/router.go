package handler

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(writerHandler *WriterHandler, workHandler *WorkHandler, opinionHandler *OpinionHandler) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")
	writers := api.Group("/writers")
	writers.POST("", writerHandler.Create)
	writers.GET("", writerHandler.List)
	writers.GET("/:id", writerHandler.GetByID)
	writers.PUT("/:id", writerHandler.Update)
	writers.DELETE("/:id", writerHandler.Delete)

	works := api.Group("/works")
	works.POST("", workHandler.Create)
	works.GET("", workHandler.List)
	works.GET("/:id", workHandler.GetByID)
	works.GET("/author/:author_id", workHandler.GetByAuthor)
	works.PUT("/:id", workHandler.Update)
	works.DELETE("/:id", workHandler.Delete)

	opinions := api.Group("/opinions")
	opinions.POST("", opinionHandler.Create)
	opinions.GET("", opinionHandler.List)
	opinions.GET("/writer/:writer_id", opinionHandler.GetByWriter)
	opinions.GET("/work/:work_id", opinionHandler.GetByWork)
	opinions.GET("/writer/:writer_id/work/:work_id", opinionHandler.GetByWriterAndWork)
	opinions.PUT("/writer/:writer_id/work/:work_id", opinionHandler.Update)
	opinions.DELETE("/writer/:writer_id/work/:work_id", opinionHandler.Delete)

	return router
}
