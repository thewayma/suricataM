package tx

import (
	"github.com/thewayma/suricataM/transporter/g"
	"github.com/thewayma/suricataM/comm/utils"
	"github.com/toolkits/consistent/rings"
)

// 建立一致性哈希环
func initNodeRings() {
	cfg := g.Config()

	JudgeNodeRing = rings.NewConsistentHashNodesRing(int32(cfg.Judge.Replicas), utils.KeysOfMap(cfg.Judge.Cluster))
	//GraphNodeRing = rings.NewConsistentHashNodesRing(int32(cfg.Graph.Replicas), KeysOfMap(cfg.Graph.Cluster))
}
