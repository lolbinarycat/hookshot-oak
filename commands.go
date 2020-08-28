package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/lolbinarycat/hookshot-oak/player"
	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/dlog"
)

func BindCommands() {
	oak.AddCommand("setAspectRatio", func(args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: setAspectRatio .75")
		}
		xToY, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			fmt.Println("failed to parse floating point ratio")
		}
		oak.SetAspectRatio(xToY)
		oak.ChangeWindow(oak.ScreenWidth, oak.ScreenHeight)
	})
	oak.AddCommand("fly", func(args []string) {
		if player.GetPlayer(0).Mods["fly"].Active() {
			if len(args) == 1 &&
				args[0] == "stop" ||
				args[0] == "end" ||
				args[0] == "halt" {

				player.GetPlayer(0).SetState(player.AirState)
			} else {
				player.GetPlayer(0).SetState(player.FlyState)
			}
		}
	})
	oak.AddCommand("mods", ModCommand)
	oak.AddCommand("mod", ModCommand)
	oak.AddCommand("exit", func(_ []string) { os.Exit(0) })
	oak.AddCommand("playerInfo", func(_ []string) { fmt.Println(player.GetPlayer(0)) })
	oak.AddCommand("kill", func(args []string) {
		player.GetPlayer(0).Die()
	})
	oak.AddCommand("dbglvl",func(args []string) {
		switch len(args) {
		case 0:
			goto ShowLevel
		case 1:
			switch args[0] {
			case "get":
				goto ShowLevel
			case "set":
				args[0] = args[1]
				goto SetLevel
			default:
				goto SetLevel
			}
		}
		return
	ShowLevel:
		fmt.Println("log level:",dlog.GetLogLevel().String())
		return
	SetLevel: // sets level to args[0]
		if i, err := strconv.Atoi(args[0]); err == nil {
			// if arg[0] is int
			dlog.SetDebugLevel(dlog.Level(i))
		} else if lvl, err := dlog.ParseDebugLevel(args[0]); err == nil {
			dlog.SetDebugLevel(lvl)
		} else {
			fmt.Println("invalid debug level", args[0])
		}
		return
	})
}

func ModCommand(args []string) {
	if len(args) == 0 {
		player.GetPlayer(0).Mods.ListModules()
	} else {
		switch args[0] {
		case "list":
			player.GetPlayer(0).Mods.ListModules()
		case "equip":
			if len(args) < 2 {
				goto NeedMoreArgs
			} else {
				player.GetPlayer(0).Mods[args[1]].Equip()
				fmt.Println("equipped", args[1])
			}
		case "grant":
			oak.RunCommand("mod","give",args[1])
			oak.RunCommand("mod","equip",args[1])
		case "give":
			if len(args) < 2 {
				goto NeedMoreArgs
			} else {
				player.GetPlayer(0).Mods[args[1]].Obtain()
				fmt.Println("gave", args[1])
			}

		case "unequip", "remove":
			if len(args) < 2 {
				goto NeedMoreArgs
			} else if args[1] == "all" {
				for _, m := range player.GetPlayer(0).Mods {
					m.Unequip()
				}
			} else {
				player.GetPlayer(0).Mods[args[1]].Unequip()
				fmt.Println("unequiped", args[1])
			}
		case "bind":
			if args[1] == "input" {
				goto BindInput
			} else if len(args) == 3 {
				oldArgs := args
				args = make([]string, 4)
				args[3] = oldArgs[2]
				args[2] = oldArgs[1]
				goto BindInput
			} else {
				fmt.Println("malformed command")
			}
		case "inputs":
			fallthrough
		case "input":
			if len(args) < 2 || args[1] == "list" {
				for _, m := range player.GetPlayer(0).Ctrls.Mod {
					fmt.Println(m)
				}
			} else if args[1] == "bind" {
				goto BindInput
			}
		default:
			fmt.Println("unknown subcommand", args[0])
		}
	}
	return
NeedMoreArgs:
	fmt.Println("not enough args")
	return
BindInput:
	if len(args) < 4 {
		fmt.Println("not enough args")
	} else {
		inpNum, err := strconv.Atoi(args[3])
		if err != nil {
			dlog.Error(err)
			return
		}
		mod, ok := player.GetPlayer(0).Mods[args[2]].(*player.CtrldPlayerModule)
		if !ok {
			fmt.Println("module", args[2], "cannot be bound")
		} else {
			pl := player.GetPlayer(0)
			mod.Bind(pl.Ctrls.Mod, inpNum)
			fmt.Println("module", args[2], "bound to input", args[3])
		}
	}
}
