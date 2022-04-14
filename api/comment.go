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
	//3.数据库插入评论，修改订单星级
	conn, _ := model.Db.Begin()
	if _, err := conn.Exec("insert into comment(order_id,rate,content) values(?,?,?)", review.OrderID, review.Rate, review.Content); err != nil {
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
	if _, err := conn.Exec("update adviser set reviews_num=reviews_num+1,total_rates=total_rates+? where id=?",
		review.Rate, order.AdviserID); err != nil {
		fmt.Println(err)
		conn.Rollback()
		response.SendResponse(c, errmsg.ERROR_DATABASE)
		return
	}
	conn.Commit()
	comment, _ := model.GetOneComment(model.Db, map[string]interface{}{"order_id": review.OrderID})
	//4.redis缓存评论信息,删除订单详情信息
	model.RemoveOrderInfoFromRedis(order.OrderID)
	model.InsertCommentToRedis(order.AdviserID, *comment)
	response.SendResponse(c, errmsg.SUCCSE, comment)
}

//订单打赏
func OrderReward(c *gin.Context) {
	var reward model.Reward
	//1.绑定参数
	if err := c.ShouldBindJSON(&reward); err != nil {
		response.SendResponse(c, errmsg.INVALID_PARAMS)
		return
	}
	//2.查询订单
	order, _ := model.GetOneOrder(model.Db, map[string]interface{}{"order_id": reward.OrderID})
	if order == nil {
		response.SendResponse(c, errmsg.ERROR_ORDER_NOT_EXIST)
		return
	} else {
		if order.Status != 1 {
			response.SendResponse(c, errmsg.ERROR_ORDER_STATUS_WRONG)
			return
		}
	}
	//3.检查用户是否有足够的金币
	u, _ := c.Get("user")
	user, _ := u.(*model.User)
	if user.Coins < reward.Reward {
		response.SendResponse(c, errmsg.ERROR_COINS_NOT_ENOUGH)
		return
	}
	//4.用户账户扣除金额，顾问账户增加金额
	adviser, _ := model.GetOneAdviser(model.Db, map[string]interface{}{"id": order.AdviserID})
	conn, _ := model.Db.Begin()
	if _, err := conn.Exec("update user set coins=coins-? where id=?", reward.Reward, order.UserID); err != nil {
		fmt.Println(err)
		conn.Rollback()
		response.SendResponse(c, errmsg.ERROR_DATABASE)
		return
	}
	if _, err := conn.Exec(`insert into transaction_user(action,id,order_id,service_type,coins,credits) 
	    values(?,?,?,?,?,?)`, 5, order.UserID, order.OrderID, order.ServiceType, user.Coins-reward.Reward, -reward.Reward); err != nil {
		conn.Rollback()
		response.SendResponse(c, errmsg.ERROR_DATABASE)
		return
	}
	if _, err := conn.Exec("update adviser set coins=coins+? where id=?", reward.Reward, order.AdviserID); err != nil {
		fmt.Println(err)
		conn.Rollback()
		response.SendResponse(c, errmsg.ERROR_DATABASE)
		return
	}
	if _, err := conn.Exec(`insert into transaction_adviser(action,id,order_id,service_type,coins,credits) 
	    values(?,?,?,?,?,?)`, 3, order.AdviserID, order.OrderID, order.ServiceType, adviser.Coins+reward.Reward, reward.Reward); err != nil {
		conn.Rollback()
		response.SendResponse(c, errmsg.ERROR_DATABASE)
		return
	}
	conn.Commit()
	response.SendResponse(c, errmsg.SUCCSE, order)
}
