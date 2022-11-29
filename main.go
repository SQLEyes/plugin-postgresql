package plugin_postgresql

import (
	. "engine"
	"engine/config"
)

type PostgreSQL struct {
	config.Plugin
	BPFFilter string
	Device    string
}

func (p *PostgreSQL) React(msg any) {
	switch v := msg.(type) {
	case config.Installed:
		plugin.Infof("%s", v.Text)
	case config.Broken:
		plugin.Infof("%s:%d->%s:%d", v.SrcIP, v.SrcPort, v.DstIP, v.DstPort)
	case config.ERROR:
		plugin.Errorf("%s", v.Text)
	}

}

var c = PostgreSQL{BPFFilter: "(tcp and port 80) or (tcp and port 443)", Device: "\\Device\\NPF_{8AAAA995-A1E9-4493-B984-2E6D03F06143}"}

var plugin = InstallPlugin(&c)

func (p *PostgreSQL) OnEvent() {
	plugin.Info("123")
}
