package mertics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

type DiscoveredLotsStatus struct {
	Address     string `json:"__address__"`
	MetricsPath string `json:"__metrics_path__"`
	Scheme      string `json:"__scheme__"`
	Job         string `json:"job"`
}

type LabelsLotsStatus struct {
	Instance string `json:"instance"`
	Job      string `json:"job"`
}

type DataOne struct {
	DiscoveredLabels DiscoveredLotsStatus `json:"discoveredLabels"`
	labels           LabelsLotsStatus     `json:"labels"`
	ScrapeUrl        string               `json:"scrapeUrl"`
	LastError        string               `json:"lastError"`
	LastScrape       string               `json:"lastScrape"`
	Health           string               `json:"health"`
}

type DataTwo struct {
	ActiveTargets []DataOne `json:"activeTargets"`
}

type LotsStatusAll struct {
	Status string  `json:"status"` //返回值
	Data   DataTwo `json:"data"`   //返回值
}

type NewDataOne struct {
	ScrapeUrl string `json:"scrapeUrl"`
	Health    string `json:"health"`
}

type NewData struct {
	Data []NewDataOne `json:"data"`
}

type Alert struct {
	GenericDataList []DataCell
	dataSelectQuery *DataSelectQuery
}

type DataCell interface {
	GetName() string
	GetHealth() string
}

type DataSelectQuery struct {
	FilterQuery   *FilterQuery
	PaginateQuery *PaginateQuery
}

type FilterQuery struct {
	Name string
}

type PaginateQuery struct {
	Limit int
	Page  int
}

type podCell NewDataOne

func (p NewDataOne) GetName() string {
	return p.ScrapeUrl
}

func (p NewDataOne) GetHealth() string {
	return p.Health
}

type deployment struct{}

func (d *deployment) toCells(std []NewDataOne) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = NewDataOne(std[i])
	}
	return cells
}

func (d *deployment) fromCells(cells []DataCell) []NewDataOne {
	deployments := make([]NewDataOne, len(cells))
	for i := range cells {
		deployments[i] = NewDataOne(cells[i].(NewDataOne))
	}

	return deployments
}

type DeploymentsResp struct {
	Items []NewDataOne `json:"items"`
}

func (d *Alert) Filter() *Alert {
	if d.dataSelectQuery.FilterQuery.Name == "" {
		return d
	}
	filteredList := []DataCell{}
	for _, value := range d.GenericDataList {
		matches := true
		objName := value.GetName()
		healthName := value.GetHealth()
		if !strings.Contains(objName, d.dataSelectQuery.FilterQuery.Name) && !strings.Contains(healthName, d.dataSelectQuery.FilterQuery.Name) {
			matches = false
			continue
		}
		if matches {
			filteredList = append(filteredList, value)
		}
	}

	d.GenericDataList = filteredList
	return d
}

func (d *Alert) Paginate() *Alert {
	limit := d.dataSelectQuery.PaginateQuery.Limit
	page := d.dataSelectQuery.PaginateQuery.Page
	//验证参数合法
	if limit <= 0 || page <= 0 {
		return d
	}

	startIndex := limit * (page - 1)
	endIndex := limit * page

	//处理最后一页
	if len(d.GenericDataList) < endIndex {
		endIndex = len(d.GenericDataList)
	}

	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	return d
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {

		method := c.Request.Method
		c.Header("Content-Type", "application/json")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Max-Age", "86400")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "X-Token, Content-Type, accessToken, authoms,Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Header("Access-Control-Allow-Credentials", "false")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func GetMertics(c *gin.Context) {
	var lotsStatus LotsStatusAll
	var NewStatus NewData

	params := new(struct {
		Title    string `form:"filter_name"`
		PageNo   int    `form:"pageNo"`
		PageSize int    `form:"pageSize"`
	})
	if err := c.Bind(params); err != nil {
		fmt.Println("params bind faid")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	err := json.Unmarshal([]byte(GetFile()), &lotsStatus)
	if err != nil {
		fmt.Println("error:", err)
	}
	// 在NewStatus结构体中循环插入多条数据
	for _, v := range lotsStatus.Data.ActiveTargets {
		NewStatus.Data = append(NewStatus.Data, NewDataOne{v.ScrapeUrl, v.Health})
	}

	// //请求参数(query)：{"pageSize":"100","filter_name":"Wljf","pageNo":"1"}
	// for _, i := range NewStatus.Data {
	// 	fmt.Println(i.ScrapeUrl)
	// }
	d := deployment{}

	selectableData := &Alert{
		GenericDataList: d.toCells(NewStatus.Data),
		dataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: params.Title},
			PaginateQuery: &PaginateQuery{
				Limit: params.PageSize,
				Page:  params.PageNo,
			},
		},
	}

	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	data := filtered.Paginate()

	//将[]DataCell类型的deployment列表转为appsv1.deployment列表
	deployments := d.fromCells(data.GenericDataList)

	da := &DeploymentsResp{
		Items: deployments,
	}

	fmt.Println(da)

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "success",
		"data":  da,
		"total": total,
	})
}

type Ids struct {
	Ids []string `json:"ids"`
}

func DeleteMertics(c *gin.Context) {
	// 打印vue传过来的参数 req.body
	id := c.Request.Body
	body, _ := ioutil.ReadAll(id)
	//body返回为{"ids":"http://prometheus-mongodb-exporter:9216/metrics,http://prometheus-mysql-exporter:9104/metrics"}
	//取出ids
	ids := gjson.Get(string(body), "ids") // go get github.com/tidwall/gjson
	//将ids转为数组
	idsArr := strings.Split(ids.String(), ",")
	//将数组转为map
	idsMap := make(map[string]string)
	for _, v := range idsArr {
		idsMap[v] = "default"
	}

	// 将数据传给Deletest函数
	Deletest(idsMap)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "接口删除成功",
	})
}

// Deletest 删除数据
func Deletest(idsMap map[string]string) {
	for k, v := range idsMap {
		fmt.Println(k, v)
	}
}

func GetFile() string {
	r, err := ioutil.ReadFile("./prom.json")
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(string(r))
	str := string(r)
	return str
}
func GetUrl() string {
	apiUrl := "https://prom.i-counting.cn/api/v1/targets"
	resp, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
	}
	r, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(r))
	str := string(r)
	ioutil.WriteFile("./prom.json", r, 0777) //注释放开写入文件，方便观看数据
	return str

}
