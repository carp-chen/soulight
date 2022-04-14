package model

import (
	"encoding/json"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

//将订单详情数据插入Redis，并设置过期时间为5分钟
func InsertOrderInfoToRedis(orderinfo OrderInfo) {
	k := "order:" + orderinfo.Order.OrderID
	v, _ := json.Marshal(orderinfo)
	r := Pool.Get()
	defer r.Close()
	r.Do("setex", k, 300, v)
}

//从Redis获取订单详情数据
func GetOrderInfoFromRedis(order_id string) (OrderInfo, error) {
	k := "order:" + order_id
	r := Pool.Get()
	defer r.Close()
	v, err := redis.String(r.Do("get", k))
	if err != nil {
		return OrderInfo{}, err
	}
	var orderinfo OrderInfo
	err = json.Unmarshal([]byte(v), &orderinfo)
	if err != nil {
		return OrderInfo{}, err
	}
	return orderinfo, nil
}

//删除redis中的订单详情数据
func RemoveOrderInfoFromRedis(order_id string) {
	k := "order:" + order_id
	r := Pool.Get()
	defer r.Close()
	r.Do("del", k)
}

//插入评论数据到redis adviser对应的哈希表中
func InsertCommentToRedis(adviser_id int, comment Comment) {
	k := "adviser:" + strconv.Itoa(adviser_id)
	v, _ := json.Marshal(comment)
	r := Pool.Get()
	defer r.Close()
	r.Do("hset", k, comment.OrderID, v)
}

//从Redis获取adviser对应的评论数据
func GetCommentsFromRedis(adviser_id int) ([]*Comment, error) {
	k := "adviser:" + strconv.Itoa(adviser_id)
	r := Pool.Get()
	defer r.Close()
	v, err := redis.Strings(r.Do("hgetall", k))
	if err != nil {
		return nil, err
	}
	var comments []*Comment
	for i := 0; i < len(v); i += 2 {
		var comment Comment
		err = json.Unmarshal([]byte(v[i+1]), &comment)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

//从Redis删除adviser对应的order评论数据
func RemoveCommentsFromRedis(adviser_id int, order_id string) {
	k := "adviser:" + strconv.Itoa(adviser_id)
	r := Pool.Get()
	defer r.Close()
	r.Do("hdel", k, order_id)
}
