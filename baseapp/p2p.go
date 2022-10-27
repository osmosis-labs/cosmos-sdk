package baseapp

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func handleQueryP2P(app *BaseApp, path []string) abci.ResponseQuery {
	// "/p2p" prefix for p2p queries
	if len(path) < 4 {
		return sdkerrors.QueryResult(
			sdkerrors.Wrap(
				sdkerrors.ErrUnknownRequest, "path should be p2p filter <addr|id> <parameter>",
			),
		)
	}

	var resp abci.ResponseQuery

	cmd, typ, arg := path[1], path[2], path[3]
	switch cmd {
	case "filter":
		switch typ {
		case "addr":
			resp = app.FilterPeerByAddrPort(arg)

		case "id":
			resp = app.FilterPeerByID(arg)
		}

	default:
		resp = sdkerrors.QueryResult(sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "expected second parameter to be 'filter'"))
	}

	return resp
}

// FilterPeerByAddrPort filters peers by address/port.
func (app *BaseApp) FilterPeerByAddrPort(info string) abci.ResponseQuery {
	if app.addrPeerFilter != nil {
		return app.addrPeerFilter(info)
	}

	return abci.ResponseQuery{}
}

// FilterPeerByID filters peers by node ID.
func (app *BaseApp) FilterPeerByID(info string) abci.ResponseQuery {
	if app.idPeerFilter != nil {
		return app.idPeerFilter(info)
	}

	return abci.ResponseQuery{}
}
