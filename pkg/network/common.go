// Package network @Author evan_yxw
// @Date 2024/8/16 15:07:00
// @Desc
package network

type PackerOptions func(*option)
type option struct {
	RawID uint32
}

func WithRaID(rawId uint32) PackerOptions {
	return func(option *option) {
		option.RawID = rawId
	}
}
