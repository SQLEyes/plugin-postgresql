package plugin_postgresql

import (
	. "engine"
	"engine/config"
	"engine/util"
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
		plugin.Errorf("%s \t", v.Text)
	}

}

var c = PostgreSQL{BPFFilter: "tcp and port 5432", Device: "\\Device\\NPF_{8AAAA995-A1E9-4493-B984-2E6D03F06143}"}

var plugin = InstallPlugin(&c)

type Sql struct {
	Type string
	Text string
}

func (p *PostgreSQL) Handle(pkt []byte) {
	buffer := util.NewByteBuffer(pkt)
	for buffer.HasNext() {
		r := Sql{}
		head := buffer.ReadShort()
		switch head {
		case 0x54:
			r.Type = "Row description"
			r.Text = RowDescriptionHandle(buffer)
		case 0x51:
			r.Type = "Simple query"
			r.Text = BaseHandle(buffer)
		case 0x43:
			r.Type = "Command completion"
			r.Text = BaseHandle(buffer)
		case 0x44:
			r.Type = "Data row"
			r.Text = DataRowHandle(buffer)
		case 0x4b:
			r.Type = "Backend key data"
			r.Text = BaseHandle(buffer)
		case 0x5a:
			r.Type = "Ready for query"
			r.Text = BaseHandle(buffer)
		case 0x52:
			r.Type = "Authentication request"
			//r.Text = BaseHandle(stream)
			return
		case 0x70:
			r.Type = "Password message"
			//r.Text = BaseHandle(stream)
			return
		case 0x53:
			r.Type = "Parameter status"
			r.Text = BaseHandle(buffer)
		case 0x0:
			r.Type = "Startup message"
			return
		default:
			r.Type = "Unknown"
		}
		plugin.Infof("%s->%s", r.Type, r.Text)
	}
}
func BaseHandle(buffer *util.ByteBuffer) string {
	length := int64(buffer.GetInt32())
	return buffer.GetString(length - 4)
}
