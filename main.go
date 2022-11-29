package plugin_postgresql

import (
	. "engine"
	"engine/abstract"
	"engine/util"
)

type PostgreSQL struct {
	abstract.Plugin
	BPFFilter string
	Device    string
}

func (p *PostgreSQL) React(msg any) (command abstract.Command) {
	switch v := msg.(type) {
	case abstract.Installed:
		plugin.Infof("%s", v.Text)
		command = abstract.Start
	case abstract.Broken:
		plugin.Infof("%s:%d->%s:%d", v.SrcIP, v.SrcPort, v.DstIP, v.DstPort)
	case abstract.ERROR:
		plugin.Errorf("%s \t", v.Text)
	}
	return
}

var c = PostgreSQL{}

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
