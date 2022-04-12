package api

import (
	"fmt"
	"soulight/model"
	"soulight/response"
	"soulight/utils/errmsg"

	"github.com/gin-gonic/gin"
)

//订单评论
func OrderReview(c *gin.Context) {
	var review model.Review
	//1.绑定参数
	if err := c.ShouldBindJSON(&review); err != nil {
		response.SendResponse(c, errmsg.INVALID_PARAMS)
		return
	}
	//2.查询订单
	order, _ := model.GetOneOrder(model.Db, map[string]interface{}{"order_id": review.OrderID})
	if order == nil {
		response.SendResponse(c, errmsg.ERROR_ORDER_NOT_EXIST)
		return
	} else {
		if order.Status != 1 {
			response.SendResponse(c, errmsg.ERROR_ORDER_STATUS_WRONG)
			return
		}
	}
	//3.插入评论，修改订单星级
	conn, _ := model.Db.Begin()
	if _, err := conn.Exec("insert into comment(order_id,content) values(?,?)", review.OrderID, review.Content); err != nil {
		fmt.Println(err)
		conn.Rollback()
		response.SendResponse(c, errmsg.ERROR_DATABASE)
		return
	}
	if _, err := conn.Exec("update orders set rate=? where order_id=?", review.Rate, review.OrderID); err != nil {
		fmt.Println(err)
		conn.Rollback()
		response.SendResponse(c, errmsg.ERROR_DATABASE)
		return
	}
	conn.Commit()
	comment, _ := model.GetOneComment(model.Db, map[string]interface{}{"order_id": review.OrderID})
	response.SendResponse(c, errmsg.SUCCSE, comment)
}

//订单打赏
