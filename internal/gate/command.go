package gate

import (
	"encoding/json"
	"github.com/wyw14/cry-153/internal/model"
)

func EncodeCommand(cmd model.GateCommand) ([]byte, error) { return json.Marshal(cmd) }
func DecodeCommand(data []byte) (model.GateCommand, error) {
	var cmd model.GateCommand
	err := json.Unmarshal(data, &cmd)
	return cmd, err
}
func IsClose(cmd model.GateCommand) bool {
	return cmd.Close && cmd.IntakeID != "" && cmd.IncidentID != ""
}
