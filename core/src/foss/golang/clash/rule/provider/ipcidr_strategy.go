package provider

import (
	"github.com/Dreamacro/clash/component/trie"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
)

type ipcidrStrategy struct {
	count           int
	shouldResolveIP bool
	trie            *trie.IpCidrTrie
}

func (i *ipcidrStrategy) Match(metadata *C.Metadata) bool {
	return i.trie != nil && i.trie.IsContain(metadata.DstIP)
}

func (i *ipcidrStrategy) Count() int {
	return i.count
}

func (i *ipcidrStrategy) ShouldResolveIP() bool {
	return i.shouldResolveIP
}

func (i *ipcidrStrategy) OnUpdate(rules []string) {
	ipCidrTrie := trie.NewIpCidrTrie()
	for _, rule := range rules {
		err := ipCidrTrie.AddIpCidrForString(rule)
		if err != nil {
			log.Warnln("invalid Ipcidr:[%s]", rule)
		} else {
			i.count++
		}
	}

	i.trie = ipCidrTrie
	i.shouldResolveIP = i.count > 0
}

func NewIPCidrStrategy() *ipcidrStrategy {
	return &ipcidrStrategy{}
}
