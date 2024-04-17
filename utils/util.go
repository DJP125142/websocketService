package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	SSO_HOST          = "https://sso.djp.org.cn/"
	GET_USER_INFO_API = "v1/account/user"
)

type HTTP_RES struct {
	CODE int                    `json:"code"`
	DATA map[string]interface{} `json:"data"`
	MSG  string                 `json:"msg"`
}

func GetNowFormatTodayTime() string {
	now := time.Now()
	dateStr := fmt.Sprintf("%02d-%02d-%02d", now.Year(), int(now.Month()), now.Day())

	return dateStr
}

func GetUserInfoByToken(token string) (map[string]interface{}, error) {
	params := make(map[string]string)
	res, err := SendGetRequest(SSO_HOST+GET_USER_INFO_API, params, token)
	if err != nil {
		return nil, err
	}
	var response HTTP_RES
	json.Unmarshal([]byte(res), &response)
	if response.CODE == 200 {
		return response.DATA, nil
	}
	return nil, errors.New(response.MSG)
}
