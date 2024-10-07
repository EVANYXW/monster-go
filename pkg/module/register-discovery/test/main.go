// Package register_discovery @Author evan_yxw
// @Date 2024/10/7 09:08:00
// @Desc

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"go.etcd.io/etcd/client/v3"
)

const (
	leaseTTL = 5 // 租约的生存时间（秒）
)

func registerService(etcdClient *clientv3.Client, serviceName, serviceAddr string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 生成唯一的服务名称
	uniqueServiceName := fmt.Sprintf("%s-%s", serviceName, uuid.New().String())

	// 将服务地址注册到 etcd
	_, err := etcdClient.Put(ctx, uniqueServiceName, serviceAddr)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}
	fmt.Printf("Service %s registered with address %s\n", uniqueServiceName, serviceAddr)
}

func registerServiceKeepAlived(etcdClient *clientv3.Client, serviceName, serviceAddr string) (clientv3.LeaseID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 创建租约
	resp, err := etcdClient.Grant(ctx, leaseTTL)
	if err != nil {
		return 0, fmt.Errorf("failed to create lease: %v", err)
	}

	// 生成唯一的服务名称
	uniqueServiceName := fmt.Sprintf("%s-%s", serviceName, uuid.New().String())

	// 将服务地址注册到 etcd，并与租约绑定
	_, err = etcdClient.Put(ctx, uniqueServiceName, serviceAddr, clientv3.WithLease(resp.ID))
	if err != nil {
		return 0, fmt.Errorf("failed to register service: %v", err)
	}
	fmt.Printf("Service %s registered with address %s\n", uniqueServiceName, serviceAddr)

	// 启动一个 goroutine 定期续租
	go keepAlive(etcdClient, resp.ID)

	return resp.ID, nil
}

func keepAlive(etcdClient *clientv3.Client, leaseID clientv3.LeaseID) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		_, err := etcdClient.KeepAliveOnce(ctx, leaseID)
		if err != nil {
			log.Printf("Failed to keep alive for lease %v: %v", leaseID, err)
			return
		}
		time.Sleep(leaseTTL / 2 * time.Second) // 每半个TTL时间续租一次
	}
}

func discoverServices(etcdClient *clientv3.Client, serviceName string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 从 etcd 获取所有服务实例，使用 WithPrefix 进行前缀匹配
	resp, err := etcdClient.Get(ctx, serviceName, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %v", err)
	}

	var services []string
	for _, kv := range resp.Kvs {
		services = append(services, string(kv.Value))
	}

	return services, nil
}
func discoverService(etcdClient *clientv3.Client, serviceName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 从 etcd 获取服务地址
	resp, err := etcdClient.Get(ctx, serviceName)
	if err != nil {
		return "", fmt.Errorf("failed to get service: %v", err)
	}

	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("service not found")
	}

	return string(resp.Kvs[0].Value), nil
}

func register() {
	//// 创建 etcd 客户端
	//etcdClient, err := clientv3.New(clientv3.Config{
	//	Endpoints:   []string{"localhost:2379"},
	//	DialTimeout: 5 * time.Second,
	//})
	//if err != nil {
	//	log.Fatalf("Failed to connect to etcd: %v", err)
	//}
	//defer etcdClient.Close()
	//
	//// 注册服务
	//registerService(etcdClient, "my-service", "localhost:8081")
	// 创建 etcd 客户端
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer etcdClient.Close()

	// 注册多个服务实例
	for i := 0; i < 5; i++ {
		addr := fmt.Sprintf("localhost:%d", 8080+i)
		registerServiceKeepAlived(etcdClient, "my-service", addr)
	}
	time.Sleep(20 * time.Second)
}

func discover() {
	// 创建 etcd 客户端
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer etcdClient.Close()

	// 发现服务
	addr, err := discoverService(etcdClient, "my-service")
	if err != nil {
		log.Fatalf("Failed to discover service: %v", err)
	}
	fmt.Printf("Discovered service address: %s\n", addr)
}

func discovers() {
	// 创建 etcd 客户端
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer etcdClient.Close()

	// 发现所有服务实例
	services, err := discoverServices(etcdClient, "my-service")
	if err != nil {
		log.Fatalf("Failed to discover services: %v", err)
	}

	// 随机选择一个服务实例
	if len(services) > 0 {
		selectedService := services[rand.Intn(len(services))]
		fmt.Printf("Selected service address: %s\n", selectedService)
	} else {
		fmt.Println("No services found")
	}
}

func main() {
	register()

}
