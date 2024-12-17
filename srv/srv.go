package srv

import (
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/srv/api"
)

func Start() {

	if conf.Config.SRV.Enabled {
		go api.Start()
	}

}
