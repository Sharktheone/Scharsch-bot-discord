package listeners

import (
	"context"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/pterodactyl/types"
	"github.com/Sharktheone/ScharschBot/srv/serversrv"
)

func StatusListener(ctx context.Context, server *conf.Server, data chan types.ChanData) {
	for {
		select {
		case d := <-data:
			if d.Event == types.WebsocketStatus {
				serversrv.HandlePower(d.Data.State, server)
			} else {
				continue
			}
		case <-ctx.Done():
			return
		}
	}
}
