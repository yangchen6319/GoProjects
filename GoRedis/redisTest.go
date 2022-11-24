package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
)

func main() {
	// 初始化结构体链接redis服务
	// 这里是一个连接池
	rdb1 := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	fmt.Println(rdb1)
	// 初始化一个context
	ctx := context.Background()
	val, err := rdb1.Get(ctx, "key").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val)
	getResult := rdb1.Get(ctx, "book")
	fmt.Println(getResult.Val(), getResult.Err())
	fmt.Println("执行一些操作")
	val1, err1 := rdb1.Do(ctx, "set", "please", "dont").Result()
	if err1 != nil {
		fmt.Println(err1)
	}
	fmt.Println(val1)
	// 这里从连接池中抽取一个链接
	cn := rdb1.Conn()
	defer cn.Close()

	if err := cn.ClientSetName(ctx, "myclient").Err(); err != nil {
		panic(err)
	}
	name, err2 := cn.ClientGetName(ctx).Result()
	if err2 != nil {
		panic(err2)
	}
	fmt.Println(name)

}
