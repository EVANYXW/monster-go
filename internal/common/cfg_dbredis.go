// ////////////////////////////////////////////////////////////////////////
//
// 文件：schema/cfg_dbredis.go
// 作者：xiexx
// 时间：2023/4/28
// 描述：Redis配置
// 说明：
//
// ////////////////////////////////////////////////////////////////////////
package common

import (
	"encoding/xml"
	"fmt"
)

type DBRedisParam struct {
	Key  string
	Type string
}

type DBRedisCfg struct {
	Id     uint32
	Op     string
	FmtKey string
	Key    string

	InParams  []DBRedisParam
	OutParams []DBRedisParam
}

type ScpDBRedisParam struct {
	Name string `xml:"name,attr"`
	Key  string `xml:"key,attr"`
	Type string `xml:"type,attr"`
}

type ScpDBRedis struct {
	XMLName   xml.Name          `xml:"request"`
	Id        uint32            `xml:"id,attr"`
	Name      string            `xml:"name,attr"`
	Desc      string            `xml:"desc,attr"`
	Op        string            `xml:"op,attr"`
	FmtKey    string            `xml:"fmt_key,attr"`
	Key       string            `xml:"key,attr"`
	InParams  []ScpDBRedisParam `xml:"in"`
	OutParams []ScpDBRedisParam `xml:"out"`
}

type ScpDBRedisRoot struct {
	XMLName xml.Name     `xml:"root"`
	Request []ScpDBRedis `xml:"request"`
}

type CfgDBRedis struct {
	datas map[uint32]*DBRedisCfg
}

func (s *CfgDBRedis) Load() error {
	data, loadErr := xml_loader("DBRedis")
	if loadErr != nil {
		return loadErr
	}

	root := new(ScpDBRedisRoot)
	err := xml.Unmarshal(data, &root)
	if err != nil {
		return fmt.Errorf("CfgDBRedis load redis request error, " + err.Error())
	}

	s.datas = make(map[uint32]*DBRedisCfg)

	for i := 0; i < len(root.Request); i++ {
		scp := &(root.Request[i])
		_, ok := s.datas[scp.Id]
		if ok {
			return fmt.Errorf("CfgDBRedis id exist, id=%d", scp.Id)
		}

		redisCfg := new(DBRedisCfg)
		redisCfg.Id = scp.Id
		redisCfg.Op = scp.Op
		redisCfg.FmtKey = scp.FmtKey
		redisCfg.Key = scp.Key

		redisCfg.InParams = make([]DBRedisParam, len(scp.InParams))
		for i, p := range scp.InParams {
			param := DBRedisParam{}
			param.Key = p.Key
			param.Type = p.Type
			redisCfg.InParams[i] = param
		}

		redisCfg.OutParams = make([]DBRedisParam, len(scp.OutParams))
		for i, p := range scp.OutParams {
			param := DBRedisParam{}
			param.Key = p.Key
			param.Type = p.Type
			redisCfg.OutParams[i] = param
		}

		s.datas[scp.Id] = redisCfg
	}

	return nil
}

func (s *CfgDBRedis) Get(id uint32) *DBRedisCfg {
	return s.datas[id]
}
