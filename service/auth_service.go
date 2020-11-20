package service

import (
	"j2pay-server/model"
	"j2pay-server/model/response"
)

var (
	// 缓存权限树结构
	authTreeCache []response.Auth

	// 缓存权限map结构
	authMapCache map[int]model.Auth

	//错误
	errs error
)

// 返回无权限分类方式的权限
func AuthTreeCache() ([]response.Auth, error) {
	if len(authTreeCache) == 0 {
		authTreeCache,errs = authTree(0)
		if errs != nil {
			return nil, errs
		}
	}
	return authTreeCache,nil
}

// 返回有权限分类方式的权限
func  AuthListCache() ([]response.Auth,error){
	if len(authTreeCache) == 0 {
		authTreeCache ,errs= authList(0)
	}
	return authTreeCache,nil
}



func authTree(pid int) ([]response.Auth,error){
	res,err := model.GetAllAuth("pid = ?", pid)
	if err != nil {
		return nil,err
	}
	for i, v := range res {
		tree, err := authTree(v.Id)
		if err != nil {
			return nil, err
		}
		res[i].Children = tree
	}
	return res,nil
}

func authList(pid int) ([]response.Auth,error) {
	res,err := model.GetAllAuth("pid != ?", pid)
	if err != nil {
		return nil, err
	}
	for i, v := range res {
		tree, err := authTree(v.Id)
		if err != nil {
			return nil,err
		}
		res[i].Children = tree
	}
	return res,nil
}

// 缓存权限
func AuthMapCache() (map[int]model.Auth,error){
	if len(authMapCache) == 0 {
		authMapCache = make(map[int]model.Auth)
		base ,err:= model.GetAllBaseAuth()
		if err != nil {
			return nil, err
		}
		for _, v := range base {
			authMapCache[v.Id] = v
		}
	}
	return authMapCache,nil
}
