package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/godcong/wego-auth-manager/model"
	"golang.org/x/exp/xerrors"
	"log"
	"strings"
)

// NotifyServerBack ...
func NotifyServerBack(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println(ctx.Params)
		paths := strings.Split(ctx.Request.URL.Path, "/")
		if len(paths) < 6 {
			log.Println("path error", paths)
			Error(ctx, xerrors.New("path error"))
			return
		}
		var back model.UserNotify
		back.Ver = paths[1]
		back.Sign = paths[3]
		back.URI = paths[4]
		back.BackType = paths[5]
		b, e := model.Get(nil, &back)
		if e != nil || !b {
			log.Println("wrong back address", paths)
			Error(ctx, xerrors.New("wrong back address"))
			return
		}
		log.Println(ctx.HandlerName())
		log.Println(ctx.Request.URL.RawPath)
		log.Println(ctx.Request.URL.Path)
		log.Println(ctx.Request.URL.EscapedPath())

		Success(ctx, back)
		log.Println(back)
		return
	}
}