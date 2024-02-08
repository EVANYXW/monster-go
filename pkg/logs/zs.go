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
	"github.com/sirupsen/logrus"
	client "github.com/zinclabs/sdk-go-zincsearch"
)

type zsHook struct {
	address    string
	user       string
	password   string
	serverName string
}

func newCsHook(zsAddress string, zsUser, zsPassword, serverName string) *zsHook {
	return &zsHook{address: zsAddress, user: zsUser, password: zsPassword, serverName: serverName}
}

func (hook *zsHook) Fire(entry *logrus.Entry) error {
	doc := newEsLog(entry)
	go hook.sendZs(doc)
	return nil
}

func (hook *zsHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

func (hook *zsHook) sendZs(doc appLogDocModel) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("send entry to es failed: ", r)
		}
	}()
	configuration := client.NewConfiguration()
	ctx := context.WithValue(context.Background(), client.ContextBasicAuth, client.BasicAuth{
		UserName: hook.user,
		Password: hook.password,
	})
	configuration.Servers = client.ServerConfigurations{
		client.ServerConfiguration{
			URL: hook.address,
		},
	}
	doc["serverName"] = hook.serverName
	apiClient := client.NewAPIClient(configuration)
	resp1, r, err := apiClient.Document.Index(ctx, "index_log").Document(doc).Execute()
	if err != nil {
		fmt.Sprintf("Error when calling `Index.Create``: %v\n", err)
	}
	// response from `Create`: MetaHTTPResponseIndex
	fmt.Sprintf("Response from `Index.Create`: %v\n", resp1)
	if r.StatusCode == 200 {

	} else {
		e, _ := err.(*client.GenericOpenAPIError)
		me, _ := e.Model().(client.MetaHTTPResponseError)
		fmt.Sprintf("`Index.Create` error: %v\n", me.GetError())
	}
}
