package handler

import (
	"errors"
	"net/http"
	"strconv"
	"vote/enum/result"
	"vote/errno"
	"vote/model"
	"vote/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TopicHandler struct {
	TopicService service.TopicService
}

func (h *TopicHandler) Insert(c *gin.Context) {
	var topicVo model.TopicVO

	token := c.GetHeader("Authorization")
	if err := c.ShouldBindJSON(&topicVo); err != nil {
		logrus.Error("bind json to model.topic error: ", err)
		c.JSON(http.StatusBadRequest, result.NewWithCode(result.BAD_REQUEST))
		return
	}
	logrus.Infof("topicVo: %+v\n", topicVo)
	topicVo.Deadline = topicVo.Deadline.Local()

	t := &model.Topic{
		Title:       topicVo.Title,
		Description: topicVo.Description,
		Deadline:    topicVo.Deadline,
	}

	s := &model.TopicSet{
		SelectType: topicVo.SelectType,
		Anonymous:  topicVo.Anonymous,
		ShowResult: topicVo.ShowResult,
		Password:   topicVo.Password,
	}

	var o []*model.TopicOption
	for _, v := range topicVo.Option {
		op := &model.TopicOption{
			OptionContent: v.OptionContent,
		}
		o = append(o, op)
	}

	if err := h.TopicService.Insert(t, s, o, token); err != nil {
		if errors.Is(err, errno.TokenInvalid) {
			logrus.Error(err)
			c.JSON(http.StatusOK, result.NewWithCode(result.TOKEN_INVALID))
			return
		}
		c.JSON(http.StatusOK, result.NewWithCode(result.SERVER_ERROR))
		return
	}

	c.JSON(http.StatusOK,
		result.NewWithCodeAndData(result.SUCCESS, nil))
}

func (h *TopicHandler) QueryAllWithTopicSet(c *gin.Context) {
	page := c.Query("page")
	size := c.Query("size")

	pint, _ := strconv.Atoi(page)
	sint, _ := strconv.Atoi(size)

	topics, err := h.TopicService.QueryAllWithTopicSet(pint, sint)
	// 将返回的 json 中的密码字段隐藏
	for _, topic := range topics {
		topic.TopicSets.Password = "******回复后可查看******"
	}
	if err != nil {
		if errors.Is(err, errno.MysqlLimitParamError) {
			c.JSON(http.StatusBadRequest, result.NewWithCode(result.BAD_REQUEST))
			return
		}
		c.JSON(http.StatusInternalServerError, result.NewWithCode(result.SERVER_ERROR))
		return
	}
	c.JSON(http.StatusOK, result.NewWithCodeAndData(result.SUCCESS, topics))
}

// QueryAllFriendlyData 返回友好数据方便前端使用（例如直接显示 name 而不是 id）
// 这个脱裤子放屁的接口完全是因为前端不会写而产生的
func (h *TopicHandler) QueryAllFriendlyData(c *gin.Context) {
	page := c.Query("page")
	size := c.Query("size")

	pint, _ := strconv.Atoi(page)
	sint, _ := strconv.Atoi(size)

	topics, total, err := h.TopicService.QueryAllFriendlyData(pint, sint)
	if err != nil {
		if errors.Is(err, errno.MysqlLimitParamError) {
			c.JSON(http.StatusBadRequest, result.NewWithCode(result.BAD_REQUEST))
			return
		}
		c.JSON(http.StatusInternalServerError, result.NewWithCode(result.SERVER_ERROR))
		return
	}
	c.JSON(http.StatusOK, result.NewWithCodeAndData(result.SUCCESS, gin.H{
		"total": total,
		"topic": topics,
	}))
}

func (h *TopicHandler) QueryByTitleFriendlyData(c *gin.Context) {
	title := c.Param("title")
	data, total, err := h.TopicService.QueryByTitleFriendlyData(title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.NewWithCode(result.SERVER_ERROR))
		return
	}
	c.JSON(http.StatusOK, result.NewWithCodeAndData(result.SUCCESS, gin.H{
		"total": total,
		"topic": data,
	}))
}

func (h *TopicHandler) QueryById(c *gin.Context) {
	topicId := c.Param("topicId")
	topic, err := h.TopicService.QueryByIdWithFmtTime(topicId)
	logrus.Info(topic)
	if err != nil {
		if errors.Is(err, errno.MysqlSelectNoData) {
			c.JSON(http.StatusOK, result.NewWithCode(result.NO_DATA))
			return
		}
		c.JSON(http.StatusInternalServerError, result.NewWithCode(result.SERVER_ERROR))
		return
	}
	c.JSON(http.StatusOK, result.NewWithCodeAndData(result.SUCCESS, topic))
}

func (h *TopicHandler) ShowResult(c *gin.Context) {
	topicId := c.Param("topicId")
	// FIXME &{1 0 NaN} float32 序列化错误
	// Panic info is: json: unsupported value: NaN
	r, err := h.TopicService.ShowResultById(topicId)
	for _, vo := range r {
		logrus.Info(vo)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, result.NewWithCode(result.SERVER_ERROR))
		return
	}
	c.JSON(http.StatusOK, result.NewWithCodeAndData(result.SUCCESS, r))
}
