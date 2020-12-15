package shared_test

import (
	"testing"

	"github.com/dave/courtney/shared"
	"github.com/dave/patsy"
	"github.com/dave/patsy/builder"
	"github.com/dave/patsy/vos"
)

func TestParseArgs(t *testing.T) {
	env := vos.Mock()
	b, err := builder.New(env, "ns")
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup()
	apath, adir, err := b.Package("a", map[string]string{
		"a.go": `package a`,
	})
	bpath, bdir, err := b.Package("a/b", map[string]string{
		"b.go": `package b`,
	})
	cpath, cdir, err := b.Package("a/c", map[string]string{
		"c.go": `package c`,
	})
	if err != nil {
		t.Fatal(err)
	}

	paths := patsy.NewCache(env)

	if err := env.Setwd(adir); err != nil {
		t.Fatal(err)
	}

	setup := shared.Setup{
		Env:   env,
		Paths: paths,
	}

	expectedA := shared.PackageSpec{
		Dir:  adir,
		Path: apath,
	}
	expectedB := shared.PackageSpec{
		Dir:  bdir,
		Path: bpath,
	}
	expectedC := shared.PackageSpec{
		Dir:  cdir,
		Path: cpath,
	}

	if err := setup.Parse([]string{"."}); err != nil {
		t.Fatal(err)
	}
	if len(setup.Packages) != 1 {
		t.Fatalf("Error in ParseArgs - wrong number of packages. Expected 1, got %d", len(setup.Packages))
	}
	if setup.Packages[0] != expectedA {
		t.Fatalf("Error in ParseArgs - wrong package. Expected %#v. Got %#v.", expectedA, setup.Packages[0])
	}

	setup = shared.Setup{
		Env:   env,
		Paths: paths,
	}
	if err := setup.Parse(nil); err != nil {
		t.Fatal(err)
	}
	if len(setup.Packages) != 3 {
		t.Fatalf("Error in ParseArgs - wrong number of packages. Expected 2, got %d", len(setup.Packages))
	}
	if setup.Packages[0] != expectedA && setup.Packages[0] != expectedB && setup.Packages[0] != expectedC {
		t.Fatal("Error in ParseArgs - wrong package.")
	}
	if setup.Packages[1] != expectedA && setup.Packages[1] != expectedB && setup.Packages[1] != expectedC {
		t.Fatal("Error in ParseArgs - wrong package.")
	}
	if setup.Packages[2] != expectedA && setup.Packages[2] != expectedB && setup.Packages[2] != expectedC {
		t.Fatal("Error in ParseArgs - wrong package.")
	}

	if err := env.Setwd(bdir); err != nil {
		t.Fatal(err)
	}

	setup = shared.Setup{
		Env:   env,
		Paths: paths,
	}
	if err := setup.Parse([]string{"."}); err != nil {
		t.Fatal(err)
	}
	if len(setup.Packages) != 1 {
		t.Fatalf("Error in ParseArgs - wrong number of packages. Expected 1, got %d", len(setup.Packages))
	}
	if setup.Packages[0] != expectedB {
		t.Fatalf("Error in ParseArgs - wrong package. Expected %#v. Got %#v.", expectedB, setup.Packages[0])
	}

	setup = shared.Setup{
		Env:   env,
		Paths: paths,
	}
	// should correctly strip "/" suffix
	if err := setup.Parse([]string{"ns/a/b/"}); err != nil {
		t.Fatal(err)
	}
	if len(setup.Packages) != 1 {
		t.Fatalf("Error in ParseArgs - wrong number of packages. Expected 1, got %d", len(setup.Packages))
	}
	if setup.Packages[0] != expectedB {
		t.Fatalf("Error in ParseArgs - wrong package. Expected %#v. Got %#v.", expectedB, setup.Packages[0])
	}
}
