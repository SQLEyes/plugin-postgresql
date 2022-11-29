package plugin_postgresql

import "engine/util"

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
func DataRowHandle(buffer *util.ByteBuffer) (r string) {
	buffer.GetInt32()
	columnSize := buffer.GetInt16()
	i := 0
	for i < int(columnSize) {
		length := buffer.GetInt32()
		r += buffer.GetString(int64(length)) + "\t"
		i++
	}
	return
}
