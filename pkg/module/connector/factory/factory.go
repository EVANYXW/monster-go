// Package connector @Author evan_yxw
// @Date 2024/9/25 23:33:00
// @Desc
package factory

import (
	"github.com/evanyxw/monster-go/pkg/module/connector"
)

type AbstractFactory interface {
	CreateManager(id int32) connector.TcpConnectorManager
}
