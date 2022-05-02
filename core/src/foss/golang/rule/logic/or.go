package logic

import (
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/rule/common"
)

type OR struct {
	*common.Base
	rules   []C.Rule
	payload string
	adapter string
	needIP  bool
}

func (or *OR) ShouldFindProcess() bool {
	return false
}

func (or *OR) RuleType() C.RuleType {
	return C.OR
}

func (or *OR) Match(metadata *C.Metadata) bool {
	for _, rule := range or.rules {
		if rule.Match(metadata) {
			return true
		}
	}

	return false
}

func (or *OR) Adapter() string {
	return or.adapter
}

func (or *OR) Payload() string {
	return or.payload
}

func (or *OR) ShouldResolveIP() bool {
	return or.needIP
}

func NewOR(payload string, adapter string) (*OR, error) {
	or := &OR{Base: &common.Base{}, payload: payload, adapter: adapter}
	rules, err := parseRuleByPayload(payload)
	if err != nil {
		return nil, err
	}

	or.rules = rules
	for _, rule := range rules {
		if rule.ShouldResolveIP() {
			or.needIP = true
			break
		}
	}

	return or, nil
}
