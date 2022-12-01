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
	server    map[string][]byte
	client    map[string][]byte
}

func (p *PostgreSQL) React(msg any) (command abstract.Command) {
	switch v := msg.(type) {
	case abstract.Installed:
		p.server = make(map[string][]byte)
		p.client = make(map[string][]byte)
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
	defer func() {
		//捕获异常
		err := recover()
		if err != nil { //条件判断，是否存在异常
			//存在异常,抛出异常
			plugin.Infof("system error %s", err)
		}
	}()

	Source := "Client"
	if strings.Index(p.BPFFilter, fmt.Sprintf("%d", broken.SrcPort)) != -1 {
		Source = "Server"
	}
	serverKey := fmt.Sprintf("%s_%s_%d", broken.SrcIP, broken.DstIP, broken.DstPort)
	clientKey := fmt.Sprintf("%s_%d", broken.DstIP, broken.DstPort)
	if len(p.server[serverKey]) != 0 && Source == "Server" {
		broken.Payload = append(p.server[serverKey], broken.Payload...)
		p.server[serverKey] = []byte{}
	}
	if len(p.client[clientKey]) != 0 && Source == "Client" {
		broken.Payload = append(p.client[clientKey], broken.Payload...)
		p.client[clientKey] = []byte{}
	}
	buffer := util.NewByteBuffer(broken.Payload)

	for buffer.HasNext() {
		r := Sql{}
		head := buffer.ReadShort()
		if head == 0x00 {
			return
		}
		if buffer.Position()+4 > buffer.Len() {
			if Source == "Server" {
				p.server[serverKey] = buffer.ReadEnd()
			} else {
				p.client[clientKey] = buffer.ReadEnd()
			}
			return
		}
		length := int64(buffer.GetInt32())
		if length < 0 {
			return
		}
		if buffer.Position(-4)+length > buffer.Len() {
			if Source == "Server" {
				buffer.Position(-1)
				p.server[serverKey] = buffer.ReadEnd()
			} else {
				buffer.Position(-1)
				p.client[clientKey] = buffer.ReadEnd()
			}
			return

		}
		switch head {
		case 0x51, 0x50:
			r.Type = "Simple query"
			r.Text = p.BaseHandle(buffer)
		case 0x43:
			r.Type = "Command completion"
			r.Text = p.BaseHandle(buffer)
			if length == 10 {
				continue
			}
		case 0x45:
			r.Type = "Error"
			r.Text = p.ErrorHandle(buffer)
		case 0x4e:
			r.Type = "Notice"
			r.Text = p.BaseHandle(buffer)
		case 0x00:
			return
		case 0x53:
			return
		default:
			r.Text = ""
			p.BaseHandle(buffer)
		}
		if r.Text != "" {
			if strings.Index(r.Text, "insert into") != -1 {
				t1 := r.Text[0:strings.Index(r.Text, "values")]
				r.Text = r.Text[strings.Index(r.Text, "values"):]
				t2 := r.Text[0:strings.Index(r.Text, ")")]
				r.Text = t1 + t2 + ")..."
			}
			plugin.Infof("[serv->%s:%d]->[client->%s:%d]: %s", broken.SrcIP, broken.SrcPort, broken.DstIP, broken.DstPort, r.Text)
		}
	}
}
func (p *PostgreSQL) BaseHandle(buffer *util.ByteBuffer) string {
	length := int64(buffer.GetInt32())
	if buffer.Position()+length > buffer.Len() {
		return string(buffer.ReadEnd())
	}
	return buffer.GetString(length - 4)
}
func (p *PostgreSQL) ErrorHandle(buffer *util.ByteBuffer) string {
	length := int64(buffer.GetInt32())
	if length < 20 {
		buffer.Read(length - 4)
		return ""
	}
	buffer.ReadShort()
	severity := buffer.GetString(7 - 1)
	buffer.ReadShort()
	text := buffer.GetString(7 - 1)
	buffer.ReadShort()
	code := buffer.GetString(7 - 1)
	buffer.ReadShort()
	msg := buffer.GetString(length - 4 - 7 - 7 - 7 - 4 - 13 - 5 - 27 - 1)
	buffer.ReadShort()
	postition := buffer.GetString(4 - 1)
	buffer.ReadShort()
	file := buffer.GetString(13 - 1)
	buffer.ReadShort()
	line := buffer.GetString(5 - 1)
	buffer.ReadShort()
	routine := buffer.GetString(27 - 1)

	return fmt.Sprintln(severity, text, code, msg, postition, file, line, routine)
}
