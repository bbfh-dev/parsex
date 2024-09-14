// A GNU-/POSIX-compiant CLI argument parsing and validation library.
//
// Example usage:
//
//	func Program(in parsex.Input, args ...string) error {
//		fmt.Printf("Running main program with options: %+v, arguments: %s", in, args)
//		return nil
//	}
//
//	func AnotherProgram(in parsex.Input, args ...string) error {
//		fmt.Printf("Running build with options: %+v, arguments: %s", in, args)
//		return nil
//	}
//
//	var BuildCLI = parsex.New(AnotherProgram, []parsex.Arg{
//		{Name: "output", Match: "--AUTO,-o", Desc: "output directory", Check: parsex.ValidPath},
//		{Name: "amount", Match: "--AUTO,-a", Desc: "amount of random stuff", Check: parsex.ValidInt},
//	})
//
//	var CLI = parsex.New(Program, []parsex.Arg{
//		{Name: "help", Match: "--AUTO", Desc: "print help message and exit"},
//		{Name: "version", Match: "--AUTO,-v", Desc: "print version and exit"},
//		{Name: "build", Match: "build", Desc: "build subcommand", Branch: BuildCLI},
//	})
//
//	func main() {
//		if err := CLI.FromString("--help -v build --amount=69 -o /tmp/ test").Run(); err != nil {
//			fmt.Println(err.Error())
//			os.Exit(1)
//		}
//	}
package parsex
