package parsex_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bbfh-dev/parsex/parsex"
)

func Program(in parsex.Input, args ...string) error {
	return nil
}

func AnotherProgram(in parsex.Input, args ...string) error {
	return nil
}

var BuildCLI = parsex.New("example build", AnotherProgram, []parsex.Arg{
	{Name: "output", Match: "--AUTO,-o", Desc: "output directory", Check: parsex.ValidPath},
})

var CLI = parsex.New("example", Program, []parsex.Arg{
	{Name: "version", Match: "--AUTO,-v", Desc: "print version and exit"},
	{Name: "amount", Match: "--AUTO,-A", Desc: "amount of random stuff", Check: parsex.ValidInt},
	{Name: "build", Match: "build", Desc: "build subcommand", Branch: BuildCLI},
})

func TestHelp(test *testing.T) {
	fmt.Println("--- VERIFY HELP MESSAGE:")
	err := CLI.FromString("--help").Run()
	if err != nil {
		test.Fatal(err)
	}
}

func TestRun(test *testing.T) {
	if err := CLI.FromString("-v -A69 build --output=/tmp/ test").Run(); err != nil {
		test.Fatal(err)
	}

	if err := CLI.FromString("-v -A69 -- build --output=/tmp/ test").Run(); err != nil {
		test.Fatal(err)
	}

	if err := CLI.FromSlice([]string{"-v", "--amount", "69", "build", "-o", "/tmp/", "test"}).Run(); err != nil {
		test.Fatal(err)
	}

	if err := CLI.Run(); err != nil {
		test.Fatal(err)
	}
}

func TestValid(test *testing.T) {
	str, ok := parsex.ValidString("example")
	if str != "example" || !ok {
		test.Fatal("Valid string must always be true and not modify the output")
	}

	str, ok = parsex.ValidInt("15")
	if str != 15 || !ok {
		test.Fatal("This int should be valid")
	}

	str, ok = parsex.ValidInt("15a")
	if str != 0 || ok {
		test.Fatal("This int should not be valid")
	}

	str, ok = parsex.ValidUint(10, 8)("15")
	if str != uint64(15) || !ok {
		test.Fatalf("This uint (%s<%d>) should be valid", reflect.TypeOf(str), str)
	}

	str, ok = parsex.ValidUint(10, 8)("258")
	if str != uint64(255) || ok {
		test.Fatalf("This uint (%s<%d>) should be valid", reflect.TypeOf(str), str)
	}

	str, ok = parsex.ValidUint(10, 8)("15a")
	if str != uint64(0) || ok {
		test.Fatal("This uint should not be valid")
	}

	str, ok = parsex.ValidPath(".")
	if len(str.(string)) == 1 || !ok {
		test.Fatal("This path should be valid")
	}
	fmt.Printf("VERIFY: \".\" is interpreted as %q\n", str)
}
