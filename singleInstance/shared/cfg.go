package shared

import (
	"encoding/json"
	"github.com/KevinZonda/GoX/pkg/iox"
)

type CfgModel struct {
	ListenAddr string `json:"listen_addr"`
	Debug      bool   `json:"debug"`
	SQL        string `json:"sql_addr"`
}

var _cfgModel *CfgModel

func GetCfg() *CfgModel {
	return _cfgModel
}

func InitCfgModel(cfgPath string) error {
	bs, err := iox.ReadAllByte(cfgPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bs, &_cfgModel)
	if err != nil {
		return err
	}
	return nil
}
