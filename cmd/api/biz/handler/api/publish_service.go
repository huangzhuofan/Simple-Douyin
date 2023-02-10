// Code generated by hertz generator.

package api

import (
	"context"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"

	api "Simple-Douyin/cmd/api/biz/model/api"
	"Simple-Douyin/cmd/api/rpc"
	"Simple-Douyin/kitex_gen/publish"
	"Simple-Douyin/pkg/constants"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// PublishAction .
// @router /douyin/publish/action/ [POST]
func PublishAction(ctx context.Context, c *app.RequestContext) {
	var err error
	var req multipart.FileHeader
	err = c.BindAndValidate(&req)
	if err != nil {
		log.Println("[ypx debug] api BindAndValidate err", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		log.Println("[ypx debug] api c.FormFile(\"data\") err", err)
		c.JSON(consts.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	fileName := filepath.Base(data.Filename)
	formPos := strings.LastIndex(fileName, ".")
	fileForm := fileName[formPos:]

	title := c.PostForm("title")

	fileContent, _ := data.Open()
	byteContainer := make([]byte, constants.MaxPublishSize)
	totalLen, err := fileContent.Read(byteContainer)
	if err != nil {
		log.Println("[ypx debug] api resultContainer err", err)
		c.JSON(consts.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	log.Println("[ypx debug] len ", totalLen)

	resultContainer := byteContainer[:totalLen]

	// var resultContainer []byte
	// for {
	// 	tmpContainer := make([]byte, 100)
	// 	tmplen, err := fileContent.Read(tmpContainer)
	// 	log.Println("[ypx debug] tmplen ", tmplen)
	// 	if err != nil {
	// 		log.Println("[ypx debug] api resultContainer err", err)
	// 		c.JSON(consts.StatusOK, Response{
	// 			StatusCode: 1,
	// 			StatusMsg:  err.Error(),
	// 		})
	// 		return
	// 	}
	// 	log.Println("[ypx debug] tmplen ", tmplen)

	// 	resultContainer = append(resultContainer, byteContainer...)
	// 	if tmplen == 0 {
	// 		break
	// 	}
	// }

	log.Println("[ypx debug] api fileName ", fileName)
	log.Println("[ypx debug] api fileForm ", fileForm)
	log.Println("[ypx debug] api title ", title)
	log.Println("[ypx debug] api len(fileContent) ", len(resultContainer))
	// fmt.Println("[ypx debug] api fileName %s, fileForm %s, title %s, len(fileContent) %s ", fileName, fileForm, title, len(byteContainer))

	log.Println("[ypx debug] api prepare to rpc.PublishAction")
	v, exist := c.Get(constants.IdentityKey)
	if !exist {
		log.Println("[ypx debug] api token does not exist")
	}
	err = rpc.PublishAction(context.Background(), &publish.PublishActionRequest{
		UserId: v.(*api.User).ID,
		Data:   resultContainer,
		Title:  title + fileForm,
	})
	if err != nil {
		log.Println("[ypx debug] api rpc.PublishAction err", err)
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	log.Println("[ypx debug] api rpc.PublishAction success")
	c.JSON(consts.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  fileName + " uploaded successfully",
	})
}

// PublishList .
// @router /douyin/publish/list/ [GET]
func PublishList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.PublishListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		log.Println("[ypx debug] api PublishList BindAndValidate err", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.PublishListResponse)

	log.Println("[ypx debug] api prepare to rpc.PublishList")
	v, _ := c.Get(constants.IdentityKey)
	videos, err := rpc.PublishList(context.Background(), &publish.PublishListRequest{
		UserId: v.(*api.User).ID,
	})
	if err != nil {
		log.Println("[ypx debug] api rpc.PublishList err", err)
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	resp.StatusCode = 0
	resp.StatusMsg = "success"

	respVs := make([]*api.Video, 0)

	for _, pv := range videos {
		pAuthor := pv.Author
		author := &api.User{
			ID:            pAuthor.Id,
			Name:          pAuthor.Name,
			FollowCount:   pAuthor.FollowCount,
			FollowerCount: pAuthor.FollowerCount,
			IsFollow:      pAuthor.IsFollow,
		}

		respVs = append(respVs, &api.Video{
			ID:            pv.Id,
			Author:        author,
			PlayURL:       pv.PlayUrl,
			CoverURL:      pv.CoverUrl,
			FavoriteCount: pv.FavoriteCount,
			CommentCount:  pv.CommentCount,
			IsFavorite:    pv.IsFavorite,
			Title:         pv.Title,
		})
	}

	resp.VideoList = respVs

	log.Println("[ypx debug] api rpc.PublishList success")

	c.JSON(consts.StatusOK, resp)
}
