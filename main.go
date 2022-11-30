package plugin_postgresql

import (
	"fmt"
	. "github.com/sqleyes/engine"
	"github.com/sqleyes/engine/abstract"
	"github.com/sqleyes/engine/util"
	"strings"
)

type PostgreSQL struct {
	abstract.Plugin
	BPFFilter string
	Device    string
	packet    []byte
}

func (p *PostgreSQL) React(msg any) (command abstract.Command) {
	switch v := msg.(type) {
	case abstract.Installed:
		plugin.Infof("%s", v.Text)
		command = abstract.Start
	case abstract.Broken:
		p.Handle(v)
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

func (p *PostgreSQL) Handle(broken abstract.Broken) {
	Source := "Client"
	if strings.Index(p.BPFFilter, fmt.Sprintf("%d", broken.SrcPort)) != -1 {
		Source = "Server"
	}
	if len(p.packet) != 0 {
		broken.Payload = append(p.packet, broken.Payload...)
		p.packet = p.packet[:0]
	}
	buffer := util.NewByteBuffer(broken.Payload)
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
			r.Text = p.DataRowHandle(buffer)
		case 0x45:
			r.Type = "Error"
			r.Text = ErrorHandle(buffer)
		case 0x4b:
			r.Type = "Backend key data"
			r.Text = BaseHandle(buffer)
		case 0x4e:
			r.Type = "Notice"
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
			r.Text = "Unknown Data"
		}
		if r.Text != "" {
			plugin.Infof("%s->%s", Source, r.Text)

		}
	}
}
func BaseHandle(buffer *util.ByteBuffer) string {
	length := int64(buffer.GetInt32())
	return buffer.GetString(length - 4)
}
func ErrorHandle(buffer *util.ByteBuffer) string {
	length := int64(buffer.GetInt32())
	severity := buffer.GetString(7)
	text := buffer.GetString(7)
	code := buffer.GetString(7)
	msg := buffer.GetString(length - 4 - 7 - 7 - 7 - 4 - 13 - 5 - 27)

	postition := buffer.GetString(4)
	file := buffer.GetString(13)
	line := buffer.GetString(5)
	routine := buffer.GetString(27)

	return fmt.Sprintln(severity, text, code, msg, postition, file, line, routine)
}
