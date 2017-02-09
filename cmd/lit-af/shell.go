package main

import (
	"fmt"

	"github.com/mit-dci/lit/lnutil"
	"github.com/mit-dci/lit/litrpc"
)

// Shellparse parses user input and hands it to command functions if matching
func (lc *litAfClient) Shellparse(cmdslice []string) error {
	var err error
	var args []string
	cmd := cmdslice[0]
	if len(cmdslice) > 1 {
		args = cmdslice[1:]
	}
	if cmd == "exit" || cmd == "quit" {
		return lc.Exit(args)
	}
	// help gives you really terse help.  Just a list of commands.
	if cmd == "help" {
		err = lc.Help(args)
		if err != nil {
			fmt.Printf("help error: %s\n", err)
		}
		return nil
	}
	// adr generates a new address and displays it
	if cmd == "adr" {
		err = lc.Adr(args)
		if err != nil {
			fmt.Printf("adr error: %s\n", err)
		}
		return nil
	}

	// bal shows the current set of utxos, addresses and score
	if cmd == "ls" {
		err = lc.Ls(args)
		if err != nil {
			fmt.Printf("ls error: %s\n", err)
		}
		return nil
	}

	// send sends coins to the address specified
	if cmd == "send" {
		err = lc.Send(args)
		if err != nil {
			fmt.Printf("send error: %s\n", err)
		}
		return nil
	}

	if cmd == "lis" { // listen for lnd peers
		err = lc.Lis(args)
		if err != nil {
			fmt.Printf("lis error: %s\n", err)
		}
		return nil
	}

	if cmd == "stop" { // stop remote node
		// actually returns an error
		return lc.Stop(args)
	}

	if cmd == "sweep" { // make lots of 1-in 1-out txs
		err = lc.Sweep(args)
		if err != nil {
			fmt.Printf("sweep error: %s\n", err)
		}
		return nil
	}

	// push money in a channel away from you
	if cmd == "push" {
		err = lc.Push(args)
		if err != nil {
			fmt.Printf("push error: %s\n", err)
		}
		return nil
	}

	if cmd == "con" { // connect to lnd host
		err = lc.Connect(args)
		if err != nil {
			fmt.Printf("con error: %s\n", err)
		}
		return nil
	}
	// fund and create a new channel
	if cmd == "fund" {
		err = lc.FundChannel(args)
		if err != nil {
			fmt.Printf("fund error: %s\n", err)
		}
		return nil
	}

	// cooperateive close of a channel
	if cmd == "close" {
		err = lc.CloseChannel(args)
		if err != nil {
			fmt.Printf("close error: %s\n", err)
		}
		return nil
	}
	if cmd == "break" {
		err = lc.BreakChannel(args)
		if err != nil {
			fmt.Printf("break error: %s\n", err)
		}
		return nil
	}
	if cmd == "say" {
		err = lc.Say(args)
		if err != nil {
			fmt.Printf("say error: %s\n", err)
		}
		return nil
	}

	if cmd == "fan" { // fan-out tx
		err = lc.Fan(args)
		if err != nil {
			fmt.Printf("fan error: %s\n", err)
		}
		return nil
	}

	/*
		if cmd == "msend" {
			err = MSend(args)
			if err != nil {
				fmt.Printf("Msend error: %s\n", err)
			}
			return nil
		}
		if cmd == "rsend" {
			err = RSend(args)
			if err != nil {
				fmt.Printf("Rsend error: %s\n", err)
			}
			return nil
		}
		if cmd == "nsend" {
			err = NSend(args)
			if err != nil {
				fmt.Printf("Nsend error: %s\n", err)
			}
			return nil
		}


		if cmd == "txs" { // show all txs
			err = Txs(args)
			if err != nil {
				fmt.Printf("txs error: %s\n", err)
			}
			return nil
		}

		if cmd == "wcon" { // connect to watch tower
			err = WCon(args)
			if err != nil {
				fmt.Printf("wcon error: %s\n", err)
			}
			return nil
		}

		if cmd == "watch" { // connect to watch tower
			err = Watch(args)
			if err != nil {
				fmt.Printf("watch error: %s\n", err)
			}
			return nil
		}



		// Peer to peer actions
		// send text message
		if cmd == "say" {
			err = Say(args)
			if err != nil {
				fmt.Printf("say error: %s\n", err)
			}
			return nil
		}

		if cmd == "fix" {
			err = Resume(args)
			if err != nil {
				fmt.Printf("fix error: %s\n", err)
			}
			return nil
		}
	*/
	fmt.Printf("Command not recognized. type help for command list.\n")
	return nil
}

func (lc *litAfClient) Exit(textArgs []string) error {
	if len(textArgs) > 0 {
		if len(textArgs) == 1 && textArgs[0] == "-h" {
			fmt.Printf(lnutil.White("exit") + "\nAlias: quit\nExit the interactive shell.\n")
			return nil
		}
		fmt.Printf("Unexpected argument: " + textArgs[0])
		return nil
	}
	return fmt.Errorf("User exit")
}

func (lc *litAfClient) Ls(textArgs []string) error {
	if len(textArgs) > 0 && textArgs[0] == "-h" {
		fmt.Printf(lnutil.White("ls\n"))
		fmt.Printf("Show various information about our current state, such as connections,\n")
		fmt.Printf("addresses, UTXO's, balances, etc.\n")
		return nil
	}

	pReply := new(litrpc.ListConnectionsReply)
	cReply := new(litrpc.ChannelListReply)
	aReply := new(litrpc.AdrReply)
	tReply := new(litrpc.TxoListReply)
	bReply := new(litrpc.BalReply)
	sReply := new(litrpc.SyncHeightReply)

	err := lc.rpccon.Call("LitRPC.ListConnections", nil, pReply)
	if err != nil {
		return err
	}
	if len(pReply.Connections) > 0 {
		fmt.Printf("\t%s\n", lnutil.Header("Peers:"))
		for _, peer := range pReply.Connections {
			fmt.Printf("%s %s\n", lnutil.White(peer.PeerNumber), peer.RemoteHost)
		}
	}

	err = lc.rpccon.Call("LitRPC.ChannelList", nil, cReply)
	if err != nil {
		return err
	}
	if len(cReply.Channels) > 0 {
		fmt.Printf("\t%s\n", lnutil.Header("Channels:"))
	}

	for _, c := range cReply.Channels {
		if c.Closed {
			fmt.Printf(lnutil.Red("Closed  "))
		} else {
			fmt.Printf(lnutil.Green("Channel "))
		}
		fmt.Printf("%s (peer %d) %s\n\t cap: %s bal: %s h: %d state: %d\n",
			lnutil.White(c.CIdx), c.PeerIdx, lnutil.OutPoint(c.OutPoint),
			lnutil.SatoshiColor(c.Capacity), lnutil.SatoshiColor(c.MyBalance), c.Height, c.StateNum)
	}

	err = lc.rpccon.Call("LitRPC.TxoList", nil, tReply)
	if err != nil {
		return err
	}
	if len(tReply.Txos) > 0 {
		fmt.Printf(lnutil.Header("\tTxos:\n"))
	}
	for i, t := range tReply.Txos {
		fmt.Printf("%d %s h:%d amt:%s %s",
			i, lnutil.OutPoint(t.OutPoint), t.Height, lnutil.SatoshiColor(t.Amt), t.KeyPath)
		if t.Delay != 0 {
			fmt.Printf(" delay: %d", t.Delay)
		}
		if !t.Witty {
			fmt.Printf(" non-witness")
		}
		fmt.Printf("\n")
	}

	err = lc.rpccon.Call("LitRPC.Address", nil, aReply)
	if err != nil {
		return err
	}
	fmt.Printf(lnutil.Header("\tAddresses:\n"))
	for i, a := range aReply.WitAddresses {
		fmt.Printf("%d %s (%s)\n", i, lnutil.Address(a), lnutil.Address(aReply.LegacyAddresses[i]))
	}
	err = lc.rpccon.Call("LitRPC.Bal", nil, bReply)
	if err != nil {
		return err
	}

	fmt.Printf("\t%s %s %s %s %s %s\n",
		lnutil.Header("Utxo:"), lnutil.SatoshiColor(bReply.TxoTotal), lnutil.Header("Conf:"), lnutil.SatoshiColor(bReply.Mature), lnutil.Header("Channel:"), lnutil.SatoshiColor(bReply.ChanTotal))

	err = lc.rpccon.Call("LitRPC.SyncHeight", nil, sReply)
	if err != nil {
		return err
	}
	fmt.Printf("%s %d\n", lnutil.Header("Sync Height:"), sReply.SyncHeight)

	return nil
}

func (lc *litAfClient) Stop(textArgs []string) error {
	if len(textArgs) > 0 && textArgs[0] == "-h" {
		fmt.Printf(lnutil.White("stop\n"))
		fmt.Printf("Shut down the lit node.\n")
		return nil
	}

	reply := new(litrpc.StatusReply)

	err := lc.rpccon.Call("LitRPC.Stop", nil, reply)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", reply.Status)

	lc.rpccon.Close()
	return fmt.Errorf("stopped remote lit node")
}

func (lc *litAfClient) Help(textArgs []string) error {
	if len(textArgs) == 0 {
		fmt.Printf("commands:\n")
		fmt.Printf("help say ls adr send fan sweep lis con fund push close break stop exit\n")
		return nil
	}
	if textArgs[0] == "help" || textArgs[0] == "-h" {
		fmt.Printf("%s%s\n", lnutil.White("help"), lnutil.OptColor("command"))
		fmt.Printf("Show information about a given command\n")
		return nil
	}
	res := make([]string, 0)
	res = append(res, textArgs[0])
	res = append(res, "-h")
	return lc.Shellparse(res)
}
