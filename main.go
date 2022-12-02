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
	c         int
	id        string
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

func (p *PostgreSQL) Handle(broken abstract.Broken) {
	p.c++

	if p.IsServer(broken) {
		p.id = fmt.Sprintf("%s:%d", broken.DstIP, broken.DstPort)
		buffer := util.NewByteBuffer(broken.Payload)
		p.HandleServer(buffer, broken)
	} else {
		p.id = fmt.Sprintf("%s:%d", broken.SrcIP, broken.SrcPort)
		if len(p.client[p.id]) > 0 {
			broken.Payload = append(p.client[p.id], broken.Payload...)
			p.client[p.id] = p.client[p.id][0:0]
			delete(p.client, p.id)
		}
		buffer := util.NewByteBuffer(broken.Payload)
		p.HandleClient(buffer, broken)
	}
}
func (p *PostgreSQL) HandleServer(buffer *util.ByteBuffer, broken abstract.Broken) {

	for buffer.HasNext() {
		msg := ServerMsg(buffer.ReadShort())

		packet := p.NeedPacketLen(buffer, 4)
		if packet != nil {
			p.client[p.id] = append(p.client[p.id], packet...)
			return
		}
		length := buffer.GetInt()
		packet = p.NeedPacketData(buffer, length)
		if packet != nil {
			p.client[p.id] = append(p.client[p.id], packet...)
			return
		}
		if strings.Index(msg.String(), "Unknown") != -1 || buffer.Position() > buffer.Position()+length-4 {
			return
		}
		payload := buffer.Read(length - 4)
		text := "Unknown"
		switch msg {
		case ReadyForQuery, BindComplete, EmptyQuery, NoData, ParseComplete, RowDescription, DataRow, Authentication, ParameterStatus, BackendKeyData:

		default:
			text = fmt.Sprintf("%s", payload)

			plugin.Green("%s [Srv->%s:%d]->[Cli->%s:%d]: %s", msg, broken.SrcIP, broken.SrcPort, broken.DstIP, broken.DstPort, text)
		}
	}

}

func (p *PostgreSQL) HandleClient(buffer *util.ByteBuffer, broken abstract.Broken) {
	for buffer.HasNext() {
		msg := ClientMsg(buffer.ReadShort())
		if msg == 0 {
			break
		}
		packet := p.NeedPacketLen(buffer, 4)
		if packet != nil {
			p.client[p.id] = append(p.client[p.id], packet...)
			return
		}
		length := buffer.GetInt()
		packet = p.NeedPacketData(buffer, length)
		if packet != nil {
			p.client[p.id] = append(p.client[p.id], packet...)
			return
		}
		payload := buffer.Read(length - 4)
		text := "Unknown"
		switch msg {
		case Parse:
			text = fmt.Sprintf("%s", payload)
			if strings.Index(text, "insert into") != -1 {
				t1 := text[0:strings.Index(text, "values")]
				text = text[strings.Index(text, "values"):]
				t2 := text[0:strings.Index(text, ")")]
				text = t1 + t2 + ")..."
			}
			plugin.Blue("%s [Cli->%s:%d]->[SrV->%s:%d]: %s", msg, broken.SrcIP, broken.SrcPort, broken.DstIP, broken.DstPort, text)
		default:
			//text = fmt.Sprintf("%s", payload)
		}
	}
}
func (p *PostgreSQL) IsServer(broken abstract.Broken) bool {
	return strings.Index(p.BPFFilter, fmt.Sprintf("%d", broken.SrcPort)) != -1
}
func (p *PostgreSQL) NeedPacketLen(buffer *util.ByteBuffer, length int) []byte {
	if buffer.Position()+length > buffer.Len() {
		buffer.Position(-1)
		return buffer.ReadEnd()
	} else {
		return nil
	}
}
func (p *PostgreSQL) NeedPacketData(buffer *util.ByteBuffer, length int) []byte {
	if buffer.Position()+length-4 > buffer.Len() {
		buffer.Position(-5)
		return buffer.ReadEnd()
	} else {
		return nil
	}
}

var c = PostgreSQL{}

var plugin = InstallPlugin(&c)

type Sql struct {
	Type string
	Text string
}
