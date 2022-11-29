package plugin_postgresql

import "github.com/sqleyes/engine/util"

type ZB string

func RowDescriptionHandle(buffer *util.ByteBuffer) (r string) {
	//length :=
	buffer.GetInt32()
	//fmt.Println(length)
	columnSize := buffer.GetInt16()
	var column []byte
	i := 0
	for i < int(columnSize) {
		s := buffer.ReadShort()
		if s != 0x00 {
			column = append(column, s)
		} else {
			r += string(column) + "\t"
			buffer.Read(18)
			column = column[0:0]
			i++
		}
	}
	return
}
func (p *PostgreSQL) DataRowHandle(buffer *util.ByteBuffer) (r string) {
	rowLen := buffer.GetInt32()
	if buffer.Position()+int64(rowLen) > buffer.Len() {
		buffer.Position(-5)
		p.packet = buffer.ReadEnd()
		return ""
	}
	columnSize := buffer.GetInt16()
	i := 0
	for i < int(columnSize) {
		i++
		length := buffer.GetInt32()
		if length == -1 {
			r += "NULL\t"
			continue
		}
		r += buffer.GetString(int64(length)) + "\t"

	}
	return
}
