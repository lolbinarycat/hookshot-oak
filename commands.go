package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/lolbinarycat/utils"
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
		if player.Mods["fly"].Active() {
			if len(args) == 1 &&
				utils.EqualsAny(args[0], "stop", "end", "halt") {

				player.SetState(AirState)
			} else {
				player.SetState(FlyState)
			}
		}
	})
	oak.AddCommand("mods",ModCommand)
	oak.AddCommand("mod",ModCommand)
	oak.AddCommand("exit", func(_ []string) {os.Exit(0)})
	oak.AddCommand("playerInfo",func(_ []string) {fmt.Println(player)})
}

func ModCommand(args []string) {
	if len(args) == 0 {
		player.Mods.ListModules()
	} else {
		switch args[0] {
		case "list":
			player.Mods.ListModules()
		case "equip":
			if len(args) < 2 {
				goto NeedMoreArgs
			} else {
				player.Mods[args[1]].Equip()
				fmt.Println("equipped",args[1])
			}
		case "unequip":
			if len(args) < 2 {
				goto NeedMoreArgs
			} else if args[1] == "all" {
				for _, m := range player.Mods {
					m.Unequip()
				}
			} else {
				player.Mods[args[1]].Unequip()
				fmt.Println("unequiped",args[1])
			}
		case "bind":
			if args[1] == "input" {
				goto BindInput
			} else if len(args) == 3 {
				oldArgs := args
				args = make([]string,4)
				args[3] = oldArgs[2]
				args[2] = oldArgs[1]
				goto BindInput
			} else {
				fmt.Println("malformed command")
			}
		case "inputs":
		case "input" :
			if len(args) < 2 || args[1] == "list" {
				for _, m := range player.Ctrls.Mod {
					fmt.Println(m)
				}
			} else if args[1] == "bind" {
				goto BindInput
			}
		default:
			fmt.Println("unknown subcommand",args[0])
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
	mod, ok := player.Mods[args[2]].(*CtrldPlayerModule)
	if !ok {
		fmt.Println("module",args[2],"cannot be bound")
	} else {
		mod.Bind(&player,inpNum)
		fmt.Println("module",args[2],"bound to input",args[3])
	}
}
}
