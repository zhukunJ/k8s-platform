package controller

import (
	"fmt"
	"k8s-platform/service"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/wonderivan/logger"
)

var Login login

type login struct{}

// PingExample godoc
// @Summary 登陆接口
// @Schemes
// @Description 登陆信息
// @Tags  登陆接口
// @Param User+Password body string true "账号名和密码"
// @Success 200 { } json "{ "code": 200,"data": {"token": "intel-accessToken-e2dcf178-42de-415a-98d3-d1721a4ac58a-1670318553870"},"msg": "success"}"
// @Router /api/login [POST]
func (l *login) Auth(ctx *gin.Context) {
	params := new(struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Login.Auth(params.UserName, params.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	token := generateusertoken(data.Username)
	ctx.JSON(http.StatusOK, gin.H{

		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"token": token,
		},
	})
}

// 账号权限控制
func (l *login) UserInfo(ctx *gin.Context) {

	token := ctx.Request.Header.Get("Authorization") // Bearer admin-accessToken-xxxxxx

	username := strings.Split(strings.Replace(token, "Bearer ", "", -1), "-accessToken-")[0]

	data, _ := service.Iopsflow.GetByName(username)

	roles := []string{}

	roles = append(roles, data.Editor)

	permissions := []string{}

	if data.Read {
		permissions = append(permissions, "read:system")
	}
	if data.Write {
		permissions = append(permissions, "write:system")
	}
	if data.Delete {
		permissions = append(permissions, "delete:system")
	}
	avatar := data.Avatar

	fmt.Println(username)
	fmt.Println(roles)
	fmt.Println(permissions)
	fmt.Println(data.Avatar)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"username":    username,
			"roles":       roles,
			"permissions": permissions,
			"avatar":      avatar,
		},
	})

}

// 退出登录
func (l *login) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

// workflow名字转换成ingress名字，添加-ing后缀
func generateusertoken(username string) (ingressName string) {
	useruuid := uuid.NewV4().String()
	usertime := time.Now().UnixNano() / 1e6
	return username + "-accessToken-" + useruuid + "-" + fmt.Sprintf("%d", usertime)
}
