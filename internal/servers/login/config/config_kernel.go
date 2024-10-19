// Package config @Author evan_yxw
// @Date 2024/8/27 18:52:00
// @Desc
package config

import (
	"github.com/evanyxw/monster-go/internal/common"
	"github.com/evanyxw/monster-go/pkg/module/module_def"
	"github.com/evanyxw/monster-go/pkg/network"
)

var (
	kernel *SchemaKernel
)

type SchemaKernel struct {
	//msgHandler module.Handler
	DBMongo *common.CfgDBMongo
	DBRedis *common.CfgDBRedis
	//LuBanTables         *Tables
	//MaxActorLevel       int32
	//Global              *DataGlobalConfig
	//CharacterMaxLevel   int32
	//CharacterBreak      *CharacterCfgCharacterBreakExt
	//CharacterAwaking    *CharacterCfgCharacterAwakingExt
	//CharacterLikability *CharacterCfgCharacterLikabilityExt
	//CharacterSkill      *CharacterCfgCharacterSkillLevelExt
	//ActorTaskMain       *ActorCfgTaskMainExt
	//ActorTaskCount      *ActorCfgTaskCountExt
	//Rewards             *ActorCfgRewardsExt
	//PlayerSkillUpgrade  *ActorCfgPlayerSkillUpgradeExt
	//KeyMask             *CfgKeyMask
}

func XSFSchema() *SchemaKernel {
	return kernel
}

func New() *SchemaKernel {
	s := &SchemaKernel{
		//msgHandler: msgHandler,
	}
	kernel = s
	return s
}

func (s *SchemaKernel) Init(baseModule module_def.IBaseModule) bool {
	//s.msgHandler.OnInit(baseModule)
	return true
}

func (s *SchemaKernel) DoRegister() {

}

func (s *SchemaKernel) GetNPManager() network.INPManager {
	return nil
}

func (s *SchemaKernel) GetStatus() int {
	return 0
}

func (s *SchemaKernel) DoRun() {
	//s.msgHandler.Start()
	s.DBMongo = new(common.CfgDBMongo)
	if err := s.DBMongo.Load(); err != nil {
		//xsf_log.Panicf("SchemaKernel Start load CfgDBMongo error, error=%v", err)
		return
	}

	s.DBRedis = new(common.CfgDBRedis)
	if err := s.DBRedis.Load(); err != nil {
		//xsf_log.Panicf("SchemaKernel Start load CfgDBRedis error, error=%v", err)
		return
	}
}

func (s *SchemaKernel) DoWaitStart() {

}

func (s *SchemaKernel) DoRelease() {

}

func (s *SchemaKernel) Update() {

}

func (s *SchemaKernel) OnOk() {
	//s.msgHandler.OnOk()
}

func (s *SchemaKernel) OnStartClose() {

}

func (s *SchemaKernel) DoClose() {

}

func (s *SchemaKernel) OnStartCheck() int {
	return 0
}

func (s *SchemaKernel) OnCloseCheck() int {
	return 0
}

func (s *SchemaKernel) GetNoWaitStart() bool {
	return true
}

func (s *SchemaKernel) MessageHandler(packet *network.Packet) {

}

func (s *SchemaKernel) OnRpcNetAccept(args []interface{}) {

}

func (s *SchemaKernel) OnRpcNetConnected(args []interface{}) {

}

func (s *SchemaKernel) OnRpcNetError(args []interface{}) {

}

func (s *SchemaKernel) OnRpcNetData(args []interface{}) {

}

func (s *SchemaKernel) OnRpcNetMessage(args []interface{}) {

}
