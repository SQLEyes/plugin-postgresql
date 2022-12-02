package plugin_postgresql

import "fmt"

type ClientMsg byte
type ServerMsg byte

const (
	Bind                     ClientMsg = 'B'
	Close                    ClientMsg = 'C'
	CopyDataClient           ClientMsg = 'd'
	CopyDoneClient           ClientMsg = 'c'
	CopyFail                 ClientMsg = 'f'
	Describe                 ClientMsg = 'D'
	Execute                  ClientMsg = 'E'
	Flush                    ClientMsg = 'H'
	FunctionCallClient       ClientMsg = 'F'
	Parse                    ClientMsg = 'P'
	PasswordMessage          ClientMsg = 'p'
	Query                    ClientMsg = 'Q'
	Sync                     ClientMsg = 'S'
	Terminate                ClientMsg = 'X'
	Authentication           ServerMsg = 'R'
	BackendKeyData           ServerMsg = 'K'
	BindComplete             ServerMsg = '2'
	CloseComplete            ServerMsg = '3'
	CommandComplete          ServerMsg = 'C'
	CopyDataServer           ServerMsg = 'd'
	CopyDoneServer           ServerMsg = 'c'
	CopyIn                   ServerMsg = 'G'
	CopyOut                  ServerMsg = 'H'
	CopyBoth                 ServerMsg = 'W'
	DataRow                  ServerMsg = 'D'
	EmptyQuery               ServerMsg = 'I'
	Error                    ServerMsg = 'E'
	FunctionCallServer       ServerMsg = 'V'
	NegotiateProtocolVersion ServerMsg = 'v'
	NoData                   ServerMsg = 'n'
	Notice                   ServerMsg = 'N'
	Notification             ServerMsg = 'A'
	ParameterDescription     ServerMsg = 't'
	ParameterStatus          ServerMsg = 'S'
	ParseComplete            ServerMsg = '1'
	PortalSuspended          ServerMsg = 's'
	ReadyForQuery            ServerMsg = 'Z'
	RowDescription           ServerMsg = 'T'
)

func (c ClientMsg) String() string {
	switch c {
	case Bind:
		return "Bind"
	case Close:
		return "Close"
	case CopyDataClient:
		return "CopyData"
	case CopyDoneClient:
		return "CopyDone"
	case CopyFail:
		return "CopyFail"
	case Describe:
		return "Describe"
	case Execute:
		return "Error or Execute"
	case Flush:
		return "Flush"
	case FunctionCallClient:
		return "FunctionCall"
	case Parse:
		return "Parse"
	case PasswordMessage:
		return "PasswordMessage"
	case Query:
		return "Query"
	case Sync:
		return "Sync"
	case Terminate:
		return "Terminate"
	default:
		return fmt.Sprintf("Unknown_%d", c)
	}
}

func (s ServerMsg) String() string {
	switch s {
	case Authentication:
		return "authentication..."
	case BackendKeyData:
		return "BackendKeyData"
	case BindComplete:
		return "BindComplete"
	case CommandComplete:
		return "Command Complete"
	case CloseComplete:
		return "CloseComplete"
	case CopyDataServer:
		return "CopyData"
	case CopyDoneServer:
		return "CopyDone"
	case CopyIn:
		return "CopyIn"
	case CopyOut:
		return "CopyOut or Flush"
	case CopyBoth:
		return "CopyBoth"
	case DataRow:
		return "DataRow or Describe"
	case EmptyQuery:
		return "EmptyQuery"
	case FunctionCallServer:
		return "FunctionCall"
	case Error:
		return "Error"
	case NegotiateProtocolVersion:
		return "NegotiateProtocolVersion"
	case NoData:
		return "NoData"
	case Notice:
		return "Notice"
	case Notification:
		return "NotificationResponse"
	case ParameterDescription:
		return "ParameterDescription"
	case ParameterStatus:
		return "ParameterStatus"
	case ParseComplete:
		return "ParseComplete"
	case PortalSuspended:
		return "PortalSuspended"
	case ReadyForQuery:
		return "ReadyForQuery"
	case RowDescription:
		return "RowDescription"
	default:
		return fmt.Sprintf("Unknown_%d", s)
	}
}
