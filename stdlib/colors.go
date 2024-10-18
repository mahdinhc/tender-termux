package stdlib

import (
	"github.com/2dprototype/tender/v/colorable"
	"github.com/2dprototype/tender"
)

var colorsModule = map[string]tender.Object{
	"stdout": &IOWriter{Value: colorable.NewColorableStdout()},
	"stderr": &IOWriter{Value: colorable.NewColorableStderr()},
	// //Color Values
	// "reset" : &tender.String{Value: "\033[0m"},           // Text Reset
	// // Regular Colors
	// "black" : &tender.String{Value: "\033[0;30m"},        // black
	// "red" : &tender.String{Value: "\033[0;31m"},          // red
	// "green" : &tender.String{Value: "\033[0;32m"},        // green
	// "yellow" : &tender.String{Value: "\033[0;33m"},       // yellow
	// "blue" : &tender.String{Value: "\033[0;34m"},         // blue
	// "purple" : &tender.String{Value: "\033[0;35m"},       // purple
	// "cyan" : &tender.String{Value: "\033[0;36m"},         // cyan
	// "white" : &tender.String{Value: "\033[0;37m"},        // white
	// // Bold
	// "bblack" : &tender.String{Value: "\033[1;30m"},       // black
	// "bred" : &tender.String{Value: "\033[1;31m"},         // red
	// "bgreen" : &tender.String{Value: "\033[1;32m"},       // green
	// "byellow" : &tender.String{Value: "\033[1;33m"},      // yellow
	// "bblue" : &tender.String{Value: "\033[1;34m"},        // blue
	// "bpurple" : &tender.String{Value: "\033[1;35m"},      // purple
	// "bcyan" : &tender.String{Value: "\033[1;36m"},        // cyan
	// "bwhite" : &tender.String{Value: "\033[1;37m"},       // white
	// // Underline
	// "ublack" : &tender.String{Value: "\033[4;30m"},       // black
	// "ured" : &tender.String{Value: "\033[4;31m"},         // red
	// "ugreen" : &tender.String{Value: "\033[4;32m"},       // green
	// "uyellow" : &tender.String{Value: "\033[4;33m"},      // yellow
	// "ublue" : &tender.String{Value: "\033[4;34m"},        // blue
	// "upurple" : &tender.String{Value: "\033[4;35m"},      // purple
	// "ucyan" : &tender.String{Value: "\033[4;36m"},        // cyan
	// "uwhite" : &tender.String{Value: "\033[4;37m"},       // white
	// // Background
	// "on_black" : &tender.String{Value: "\033[40m"},       // black
	// "on_red" : &tender.String{Value: "\033[41m"},         // red
	// "on_green" : &tender.String{Value: "\033[42m"},       // green
	// "on_yellow" : &tender.String{Value: "\033[43m"},      // yellow
	// "on_blue" : &tender.String{Value: "\033[44m"},        // blue
	// "on_purple" : &tender.String{Value: "\033[45m"},      // purple
	// "on_cyan" : &tender.String{Value: "\033[46m"},        // cyan
	// "on_white" : &tender.String{Value: "\033[47m"},       // white
	// // High Intensty
	// "iblack" : &tender.String{Value: "\033[0;90m"},       // black
	// "ired" : &tender.String{Value: "\033[0;91m"},         // red
	// "igreen" : &tender.String{Value: "\033[0;92m"},       // green
	// "iyellow" : &tender.String{Value: "\033[0;93m"},      // yellow
	// "iblue" : &tender.String{Value: "\033[0;94m"},        // blue
	// "ipurple" : &tender.String{Value: "\033[0;95m"},      // purple
	// "icyan" : &tender.String{Value: "\033[0;96m"},        // cyan
	// "iwhite" : &tender.String{Value: "\033[0;97m"},       // white
	// // Bold High Intensty
	// "biblack" : &tender.String{Value: "\033[1;90m"},      // black
	// "bired" : &tender.String{Value: "\033[1;91m"},        // red
	// "bigreen" : &tender.String{Value: "\033[1;92m"},      // green
	// "biyellow" : &tender.String{Value: "\033[1;93m"},     // yellow
	// "biblue" : &tender.String{Value: "\033[1;94m"},       // blue
	// "bipurple" : &tender.String{Value: "\033[1;95m"},     // purple
	// "bicyan" : &tender.String{Value: "\033[1;96m"},       // cyan
	// "biwhite" : &tender.String{Value: "\033[1;97m"},      // white
	// // High Intensty backgrounds
	// "on_iblack" : &tender.String{Value: "\033[0;100m"},   // black
	// "on_ired" : &tender.String{Value: "\033[0;101m"},     // red
	// "on_igreen" : &tender.String{Value: "\033[0;102m"},   // green
	// "on_iyellow" : &tender.String{Value: "\033[0;103m"},  // yellow
	// "on_iblue" : &tender.String{Value: "\033[0;104m"},    // blue
	// "on_ipurple" : &tender.String{Value: "\033[10;95m"},  // purple
	// "on_icyan" : &tender.String{Value: "\033[0;106m"},    // cyan
	// "on_iwhite" : &tender.String{Value: "\033[0;107m"},   // white
}

// func colorablePrint(args ...tender.Object) (ret tender.Object, err error) {
	// str := ""
	// for i, arg := range args {
		// s, _ := tender.ToString(arg)
		// str += s
		// if i < len(args) - 1 {
			// str += " "
		// }
	// }
	// fmt.Fprint(colorable.NewColorableStdout(), str)
	// return nil, nil
// }

// func colorablePrintln(args ...tender.Object) (ret tender.Object, err error) {
	// str := ""
	// for i, arg := range args {
		// s, _ := tender.ToString(arg)
		// str += s
		// if i < len(args) - 1 {
			// str += " "
		// }
	// }
	// fmt.Fprintln(colorable.NewColorableStdout(), str)
	// return nil, nil
// }


