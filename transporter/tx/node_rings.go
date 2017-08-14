package tx

import (
	"github.com/thewayma/suricataM/comm/utils"
	"github.com/thewayma/suricataM/transporter/g"
	"github.com/toolkits/consistent/rings"
)

// 建立一致性哈希环
func initNodeRings() {
	cfg := g.Config()

	CheckerNodeRing = rings.NewConsistentHashNodesRing(int32(cfg.Checker.Replicas), utils.KeysOfMap(cfg.Checker.Cluster))
}
