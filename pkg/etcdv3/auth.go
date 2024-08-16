/**
 * @api post etcdv3.
 *
 * User: yunshengzhu
 * Date: 2021/12/26
 * Time: 下午7:34
 */
package etcdv3

import (
	"context"
	"fmt"
)

type Auth struct {
	etcd *Etcd
}

func NewAuth(etcd *Etcd) *Auth {
	return &Auth{etcd: etcd}
}

type User struct {
	Name     string
	Password string
	Role     string
}

func (a *Auth) OpenAuth() error {
	if _, err := a.etcd.client.AuthEnable(context.TODO()); err != nil {
		fmt.Println("etcd auth enable " + err.Error())
		return err
	}
	return nil
}

func (a *Auth) CloseAuth() error {
	if _, err := a.etcd.client.AuthDisable(context.TODO()); err != nil {
		fmt.Println("etcd auth enable " + err.Error())
		return err
	}
	return nil
}

func (a *Auth) AuthSetupRoot(user User) error {
	user.Role = "root"
	authUser, err := a.etcd.client.UserGet(context.TODO(), user.Name)
	if err == nil {
		if authUser.Header.MemberId > 0 {
			return nil
		}
	}
	if _, err := a.etcd.client.UserAdd(context.TODO(), user.Name, user.Password); err != nil {
		fmt.Println("etcd user add " + err.Error())
		return err
	}
	if _, err := a.etcd.client.RoleAdd(context.TODO(), user.Role); err != nil {
		fmt.Println("etcd role add " + err.Error())
		return err
	}
	if _, err := a.etcd.client.UserGrantRole(context.TODO(), user.Name, user.Role); err != nil {
		fmt.Println("etcd user grant role  " + err.Error())
		return err
	}
	return nil
}

func (a *Auth) DelUserRole(username string) error {
	if _, err := a.etcd.client.UserDelete(context.TODO(), username); err != nil {
		return err
	}

	if _, err := a.etcd.client.RoleDelete(context.TODO(), username); err != nil {
		return err
	}
	return nil
}

func (a *Auth) CreateUserRole(username, password string) error {
	//创建用户
	if _, err := a.etcd.client.UserAdd(context.TODO(), username, password); err != nil {
		return err
	}

	if _, err := a.etcd.client.RoleAdd(context.TODO(), username); err != nil {
		return err
	}

	if _, err := a.etcd.client.UserGrantRole(context.TODO(), username, username); err != nil {
		return err
	}

	//角色目录授权
	auths := GetAuthPaths(username)
	for _, auth := range auths {
		if _, err := a.etcd.client.RoleGrantPermission(context.TODO(), username, auth.Path, auth.PathEnd, auth.Type); err != nil {
			return err
		}
	}
	return nil
}
