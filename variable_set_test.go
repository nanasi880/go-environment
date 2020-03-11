package environment_test

import (
	"os"
	"testing"

	"go.nanasi880.dev/environment"
)

func TestVariableSet(t *testing.T) {

	_ = os.Setenv("ENV_TEST_STRING", "this is string")
	_ = os.Setenv("ENV_TEST_INT", "100")

	set := environment.NewVariableSet("test", environment.ContinueOnError)

	sp := set.String("ENV_TEST_STRING", "this is default", "usage")
	ip := set.Int("ENV_TEST_INT", -100, "usage")

	if err := set.Parse(os.Environ()); err != nil {
		t.Fatal(err)
	}

	if *sp != "this is string" {
		t.Fatal(*sp)
	}
	if *ip != 100 {
		t.Fatal(*ip)
	}
}
