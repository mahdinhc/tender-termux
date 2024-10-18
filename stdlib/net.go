package stdlib

import (
	"net"
	"github.com/2dprototype/tender"
)

var netModule = map[string]tender.Object{
	"dnslookup": &tender.UserFunction{Value: netDnsLookup},
	"resolve_tcp_addr": &tender.UserFunction{Value: netResolveTCPAddr},
	"resolve_udp_addr": &tender.UserFunction{Value: netResolveUDPAddr},
	"dial": &tender.UserFunction{Value: netDial},
	"dialtcp": &tender.UserFunction{Value: netDialTCP},
}


func netDialTCP(args ...tender.Object) (tender.Object, error) {
	if len(args) != 2 {
		return nil, tender.ErrWrongNumArguments
	}
	network, _ := args[0].(*tender.String)
	address, _ := args[1].(*tender.String)
	tcpAddr, err := net.ResolveTCPAddr(network.Value, address.Value)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP(network.Value, nil, tcpAddr)
	if err != nil {
		return wrapError(err), nil
	}
	return makeNetConn(conn), nil
}

func netDial(args ...tender.Object) (tender.Object, error) {
	if len(args) != 2 {
		return nil, tender.ErrWrongNumArguments
	}
	network, _ := args[0].(*tender.String)
	address, _ := args[1].(*tender.String)
	conn, err := net.Dial(network.Value, address.Value)
	if err != nil {
		return wrapError(err), nil
	}
	return makeNetConn(conn), nil
}

func makeNetConn(conn net.Conn) *tender.ImmutableMap {
	return &tender.ImmutableMap{
		Value: map[string]tender.Object{
			"close": &tender.UserFunction{
				Value: FuncARE(conn.Close),
			},
			"read": &tender.UserFunction{
				Value: FuncAYRIE(conn.Read),
			},
			"write": &tender.UserFunction{
				Value: FuncAYRIE(conn.Write),
			},
			"local_addr": &tender.UserFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					addr := conn.LocalAddr()
					return &tender.ImmutableMap{
						Value: map[string]tender.Object{
							"string" :  &tender.UserFunction{Value: FuncARS(addr.String)},
							"network" :  &tender.UserFunction{Value: FuncARS(addr.Network)},
						},
					}, nil
				},
			},
			"remote_addr": &tender.UserFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 {
						return nil, tender.ErrWrongNumArguments
					}
					addr := conn.RemoteAddr()
					return &tender.ImmutableMap{
						Value: map[string]tender.Object{
							"string" :  &tender.UserFunction{Value: FuncARS(addr.String)},
							"network" :  &tender.UserFunction{Value: FuncARS(addr.Network)},
						},
					}, nil
				},
			},
			"set_deadline": &tender.UserFunction{
				Value: FuncATRE(conn.SetDeadline),
			},
			"set_readdeadline": &tender.UserFunction{
				Value: FuncATRE(conn.SetReadDeadline),
			},
			"set_writedeadline": &tender.UserFunction{
				Value: FuncATRE(conn.SetWriteDeadline),
			},
		},
	}
}

func netDnsLookup(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 1 {
		return nil, tender.ErrWrongNumArguments
	}
	host, _ := args[0].(*tender.String)
	addresses, err := net.LookupHost(host.Value)
	if err != nil {
		return wrapError(err), nil
	}
	results := make([]tender.Object, len(addresses))
	for i, addr := range addresses {
		results[i] = &tender.String{Value: addr}
	}
	return &tender.Array{Value: results}, nil
}


func netResolveTCPAddr(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 2 {
		return nil, tender.ErrWrongNumArguments
	}
	network, _ := args[0].(*tender.String)
	address, _ := args[1].(*tender.String)
	tcpAddr, err := net.ResolveTCPAddr(network.Value, address.Value)
	if err != nil {
		return wrapError(err), nil
	}
	return &tender.String{Value: tcpAddr.String()}, nil
}

func netResolveUDPAddr(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 2 {
		return nil, tender.ErrWrongNumArguments
	}
	network, _ := args[0].(*tender.String)
	address, _ := args[1].(*tender.String)
	udpAddr, err := net.ResolveUDPAddr(network.Value, address.Value)
	if err != nil {
		return wrapError(err), nil
	}
	return &tender.String{Value: udpAddr.String()}, nil
}