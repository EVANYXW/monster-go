/**
 * @api post logs.
 *
 * User: yunshengzhu
 * Date: 2022/5/10
 * Time: 7:35 PM
 */
package logs

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

type esHook struct {
	client *elastic.Client
}

//newEsHook 初始化
func newEsHook(esAddress []string, esUser, esPassword string) *esHook {
	es, err := elastic.NewClient(
		elastic.SetURL(esAddress...),
		elastic.SetBasicAuth(esUser, esPassword),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(15*time.Second),
		elastic.SetErrorLog(log.New(os.Stderr, "ES:", log.LstdFlags)),
	)
	if err != nil {
		log.Fatal("failed to create Elastic V6 Client: ", err)
	}
	return &esHook{client: es}
}

//Fire logrus hook interface 方法
func (hook *esHook) Fire(entry *logrus.Entry) error {
	doc := newEsLog(entry)
	go hook.sendEs(doc)
	return nil
}

//Levels logrus hook interface 方法
func (hook *esHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

//sendEs 异步发送日志到es
func (hook *esHook) sendEs(doc appLogDocModel) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("send entry to es failed: ", r)
		}
	}()
	_, err := hook.client.Index().Index(doc.indexName()).Type("_doc").BodyJson(doc).Do(context.Background())
	if err != nil {
		log.Println(err)
	}

}

//appLogDocModel es model
type appLogDocModel map[string]interface{}

func newEsLog(e *logrus.Entry) appLogDocModel {
	ins := map[string]interface{}{}
	for kk, vv := range e.Data {
		ins[kk] = vv
	}
	ins["time"] = time.Now().Local()
	ins["lvl"] = e.Level
	ins["message"] = e.Message
	if e.Caller != nil {
		ins["file"] = e.Caller.File
		ins["line"] = e.Caller.Line
		ins["func"] = e.Caller.Function
	}
	return ins
}

// indexName es index name 时间分割
func (m *appLogDocModel) indexName() string {
	return "hlga-" + time.Now().Local().Format("2006-01-02")
}
