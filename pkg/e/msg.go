package e

var MsgFlags = map[int]string{
	SUCCESS:                        "ok",
	ERROR:                          "fail",
	INVALID_PARAMS:                 "请求参数错误",
	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",
	ERROR_AUTH_TOKEN:               "Token获取失败",
	ERROR_AUTH:                     "Token错误",
	ERROR_UPLOAD_ATTACHMENT:        "附件上传失败",
	ERROR_NOT_FOUND_FLAG:           "未找到指定立志",
	ERROR_NOT_FOUND_USER:           "未找到当前用户",
	ERROR_FLAGER_NOT_CURRENT_USER:  "当前用户非当前立志者",
	ERROR_CLOSED_FLAG:              "当前flag已经关闭",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
