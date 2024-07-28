package conf

type Format struct {
	Reconfigure bool `yaml:"reconfigure"`
	Enabled     bool `yaml:"enabled"`
	Discord     struct {
		ServerID              string   `yaml:"serverID"`
		ServerName            string   `yaml:"serverName"`
		Token                 string   `yaml:"token"`
		WhitelistRemoveRoleID []string `yaml:"adminWhitelistRemoveRoleID"`
		WhitelistWhoisRoleID  []string `yaml:"adminWhitelistWhoisRoleID"`
		WhitelistBanRoleID    []string `yaml:"adminWhitelistBanRoleID"`
		WhitelistServerRoleID []string `yaml:"whitelistServerRoleID"`
		EmbedErrorIcon        string   `yaml:"embedErrorIcon"`
		EmbedErrorAuthorURL   string   `yaml:"embedErrorAuthorURL"`
		FooterIcon            bool     `yaml:"footerIcon"`
	} `yaml:"discord"`
	Pterodactyl struct {
		Enabled                bool     `yaml:"enabled"`
		RegexRemoveAnsi        string   `yaml:"regexRemoveAnsi"`
		APIKey                 string   `yaml:"APIKey"`
		PanelURL               string   `yaml:"panelURL"`
		WhitelistAddCommand    string   `yaml:"whitelistAddCommand"`
		WhitelistRemoveCommand string   `yaml:"whitelistRemoveCommand"`
		ChatCommand            string   `yaml:"chatCommand"`
		Servers                []Server `yaml:"servers"`
	} `yaml:"pterodactyl"`
	SRV struct {
		API struct {
			Port     int    `yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
		} `yaml:"api"`
		States struct {
			Online   string `yaml:"online"`
			Offline  string `yaml:"offline"`
			Starting string `yaml:"starting"`
			Stopping string `yaml:"stopping"`
		}
	} `yaml:"srv"`
	Whitelist struct {
		Enabled     bool `yaml:"enabled"`
		MaxAccounts []struct {
			RoleID string `yaml:"roleID"`
			Max    int    `yaml:"max"`
		} `yaml:"maxAccounts"`
		BannedUsersToMaxAccounts bool     `yaml:"bannedUsersToMaxAccounts"`
		KickUnWhitelisted        bool     `yaml:"kickUnWhitelisted"`
		KickCommand              string   `yaml:"kickCommand"`
		Servers                  []string `yaml:"serverNames"`
		Database                 struct {
			Provider         string `yaml:"provider"`
			Host             string `yaml:"host"`
			Port             uint16 `yaml:"port"`
			User             string `yaml:"user"`
			Pass             string `yaml:"pass"`
			DatabaseName     string `yaml:"databaseName"`
			WhitelistTable   string `yaml:"whitelistCollectionName"`
			BanTable         string `yaml:"banCollectionName"`
			ReWhitelistTable string `yaml:"reWhitelistCollectionName"`
			ReportTable      string `yaml:"reportCollectionName"`
			WaitListTable    string `yaml:"waitListTable"`
			RolesTable       string `yaml:"rolesTable"`
		} `yaml:"db"`

		Roles struct {
			ServerRoleID      []string `yaml:"serverRoleID"`
			RemoveUserWithout bool     `yaml:"removeUserWithout"`
			ReWhitelistWith   bool     `yaml:"reWhitelistWith"`
		} `yaml:"roles"`
		Report struct {
			Enabled               bool     `yaml:"enabled"`
			ChannelID             []string `yaml:"channelID"`
			Ping                  bool     `yaml:"ping"`
			PingRoleID            []string `yaml:"pingRoleID"`
			Roles                 []string `yaml:"roles"`
			PlayerNotifyDM        bool     `yaml:"playerNotifyDM"`
			PlayerNotifyChannelID []string `yaml:"playerNotifyChannelID"`
		} `yaml:"report"`
		Luckperms struct {
			Enabled      bool   `yaml:"enabled"`
			SetRole      string `yaml:"setRole"`
			DatabaseHost string `yaml:"databaseHost"`
			DatabasePort uint16 `yaml:"databasePort"`
			DatabaseUser string `yaml:"databaseUser"`
			DatabasePass string `yaml:"databasePass"`
		} `yaml:"luckperms"`
	} `yaml:"whitelist"`
}

type Server struct {
	ServerID      string `yaml:"serverID"`
	ServerName    string `yaml:"serverName"`
	StateMessages struct {
		Enabled        bool     `yaml:"enabled"`
		Start          string   `yaml:"start"`
		Stop           string   `yaml:"stop"`
		Online         string   `yaml:"online"`
		Offline        string   `yaml:"offline"`
		StartEnabled   bool     `yaml:"startEnabled"`
		StopEnabled    bool     `yaml:"stopEnabled"`
		OfflineEnabled bool     `yaml:"offlineEnabled"`
		OnlineEnabled  bool     `yaml:"onlineEnabled"`
		ChannelID      []string `yaml:"channelID"`
	} `yaml:"stateMessages"`
	ChannelInfo struct {
		Enabled       bool     `yaml:"enabled"`
		ChannelID     []string `yaml:"channelID"`
		Format        string   `yaml:"format"`
		OfflineFormat string   `yaml:"offlineFormat"`
		InfoState     string   `yaml:"-"`
		Interval      string   `yaml:"interval"`
	} `yaml:"channelInfo"`
	PowerActionsRoleIDs struct {
		Start   []string `yaml:"start"`
		Stop    []string `yaml:"stop"`
		Restart []string `yaml:"restart"`
	} `yaml:"powerActionsRoleIDs"`
	Console struct {
		Enabled       bool     `yaml:"enabled"`
		MessageLines  int      `yaml:"messageLines"`
		MaxTime       string   `yaml:"maxTime"`
		ChannelID     []string `yaml:"channelID"`
		Reverse       bool     `yaml:"reverse"`
		ReverseRoleID []string `yaml:"reverseRoleID"`
		ReversePrefix string   `yaml:"reversePrefix"`
	} `yaml:"console"`
	Chat struct {
		Enabled      bool     `yaml:"enabled"`
		Embed        bool     `yaml:"embed"`
		Prefix       string   `yaml:"prefix"`
		ChannelID    []string `yaml:"channelID"`
		EmbedFooter  bool     `yaml:"embedFooter"`
		EmbedOneLine bool     `yaml:"oneLine"`
		FooterIcon   bool     `yaml:"footerIcon"`
		Reverse      bool     `yaml:"reverse"`
		RoleID       []string `yaml:"roleID"`
	} `yaml:"chat"`
	SRV struct {
		Enabled    bool     `yaml:"enabled"`
		OneLine    bool     `yaml:"oneLine"`
		Footer     bool     `yaml:"footer"`
		ChannelID  []string `yaml:"channelID"`
		FooterIcon bool     `yaml:"footerIcon"`
		Events     struct {
			PlayerJoinLeft   bool `yaml:"playerJoinLeft"`
			Advancement      bool `yaml:"advancement"`
			RootAdvancements bool `yaml:"rootAdvancements"`
			Death            bool `yaml:"death"`
		} `yaml:"events"`
	} `yaml:"srv"`
}
