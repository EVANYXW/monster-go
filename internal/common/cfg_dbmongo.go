package common

import (
	"encoding/xml"
	"fmt"
)

type DBMongoCfg struct {
	Id        uint32
	Op        string
	ColName   string
	Sort      string
	PageCount uint32

	Filters []string
	Datas   []string
}

type ScpDBMongoRoot struct {
	XMLName xml.Name `xml:"root"`
	Db      []ScpDB  `xml:"db"`
}

type ScpDB struct {
	XMLName    xml.Name   `xml:"db"`
	Name       string     `xml:"name,attr"`
	Collection []ScpDBCol `xml:"collection"`
}

type ScpDBCol struct {
	XMLName xml.Name       `xml:"collection"`
	Name    string         `xml:"name,attr"`
	Request []ScpDBReqeust `xml:"request"`
}

type ScpDBReqeust struct {
	XMLName   xml.Name      `xml:"request"`
	Id        uint32        `xml:"id,attr"`
	Name      string        `xml:"name,attr"`
	Desc      string        `xml:"desc,attr"`
	Op        string        `xml:"op,attr"`
	PageCount uint32        `xml:"page_count,attr"`
	Sort      string        `xml:"sort,attr"`
	Filters   []ScpDBFilter `xml:"filter"`
	Datas     []ScpDBData   `xml:"data"`
}

type ScpDBFilter struct {
	XMLName xml.Name `xml:"filter"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
}

type ScpDBData struct {
	XMLName xml.Name `xml:"data"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
}

type CfgDBMongo struct {
	datas map[uint32]*DBMongoCfg
}

func (s *CfgDBMongo) Load() error {
	data, loadErr := xml_loader("DBMongo")
	if loadErr != nil {
		return loadErr
	}

	root := new(ScpDBMongoRoot)
	err := xml.Unmarshal(data, &root)
	if err != nil {
		return fmt.Errorf("SchemaDBMongo load mongo request error, " + err.Error())
	}

	s.datas = make(map[uint32]*DBMongoCfg)

	for i := 0; i < len(root.Db); i++ {
		db := root.Db[i]
		for J := 0; J < len(db.Collection); J++ {
			collection := db.Collection[J]
			for Z := 0; Z < len(collection.Request); Z++ {
				request := collection.Request[Z]

				_, ok := s.datas[request.Id]
				if ok {
					return fmt.Errorf("CfgDBMongo id exist, id=%d", request.Id)
				}

				cfg := new(DBMongoCfg)
				cfg.Id = request.Id
				cfg.Op = request.Op
				cfg.ColName = fmt.Sprintf("%s-%s", db.Name, collection.Name)
				cfg.Sort = request.Sort
				cfg.PageCount = request.PageCount

				s.datas[request.Id] = cfg

				for _, f := range request.Filters {
					cfg.Filters = append(cfg.Filters, f.Name)
				}
				for _, d := range request.Datas {
					cfg.Datas = append(cfg.Datas, d.Name)
				}
			}
		}
	}

	return nil
}

func (s *CfgDBMongo) Get(id uint32) *DBMongoCfg {
	return s.datas[id]
}
