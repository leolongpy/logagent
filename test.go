package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v", err)
		return
	}

	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	str := `[{"path":"e:/desktop/log3.log","topic":"log3"},{"path":"e:/desktop/log2.log","topic":"log2"}]`
	_, err = cli.Put(ctx, "log_192.168.0.106", str)
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v", err)
		return
	}
	cancel()
}
