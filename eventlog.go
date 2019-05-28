package sdk

import (
	"github.com/dappledger/AnnChain-go-sdk/types"
	"github.com/dappledger/AnnChain-go-sdk/abi"
	"fmt"
	"bytes"
)

type EventLog struct {
	types.Receipt
	abi.ABI
}

func (l *EventLog) Print(names ...string) string {

	sb := new(bytes.Buffer)

	for _, name := range names {
		e, ok := l.Events[name];
		if !ok {
			continue
		}
		for _, r := range l.Receipt.Logs {
			eid := e.Id()
			if eid == r.Topics[0] {

				vals, err := e.Inputs.UnpackValues(r.Data)
				if err != nil {
					return err.Error()
				}
				nonIndexedArgs := e.Inputs.NonIndexed()
				if len(vals) != len(nonIndexedArgs) {
					return fmt.Sprintf("unexpected abi, event %s input args length should be %d, but got %d", e.Name, len(vals), len(nonIndexedArgs))
				}

				sb.WriteString(e.Name + "(")
				for i, args := range nonIndexedArgs {
					sb.WriteString(args.Name)
					sb.WriteString(fmt.Sprintf(`:"%v"`, vals[i]))
					if i != len(nonIndexedArgs)-1 {
						sb.WriteString(",")
					}
				}
				sb.WriteString(")\n")
				break
			}
		}
	}
	ret := sb.String()
	if ret == "" {
		return ""
	}
	return ret[:len(ret)-1]
}
