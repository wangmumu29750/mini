package mock

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	phonePattern    = regexp.MustCompile(`^\d{11}$`)
	bankCardPattern = regexp.MustCompile(`^\d{12,24}$`)
	idCardPattern   = regexp.MustCompile(`^(\d{15}|\d{17}[\dXx])$`)
)

func VerifyIdentity(realName, idCardNo, phone, bankCardNo string) error {
	if strings.TrimSpace(realName) == "" {
		return fmt.Errorf("姓名不能为空")
	}
	if !idCardPattern.MatchString(strings.TrimSpace(idCardNo)) {
		return fmt.Errorf("身份证号格式不正确")
	}
	if !phonePattern.MatchString(strings.TrimSpace(phone)) {
		return fmt.Errorf("手机号格式不正确")
	}
	if !bankCardPattern.MatchString(strings.TrimSpace(bankCardNo)) {
		return fmt.Errorf("银行卡号格式不正确")
	}
	return nil
}
