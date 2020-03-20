package configs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/addrs"
)

// TestNewModule_provider_fqns exercises module.gatherProviderLocalNames()
func TestNewModule_provider_local_name(t *testing.T) {
	mod, diags := testModuleFromDir("testdata/providers-explicit-fqn")
	if diags.HasErrors() {
		t.Fatal(diags.Error())
	}

	p := addrs.NewLegacyProvider("foo")
	if name, exists := mod.ProviderLocalNames[p]; !exists {
		fmt.Printf("%#v\n", mod.ProviderLocalNames)
		t.Fatal("provider FQN foo/test not found")
	} else {
		if name != "foo-test" {
			t.Fatalf("provider localname mismatch: got %s, want foo-test", name)
		}
	}

	// ensure the reverse lookup (fqn to local name) works as well
	localName := mod.LocalNameForProvider(p)
	if localName != "foo-test" {
		t.Fatal("provider local name not found")
	}
}

// This test validates the provider FQNs set in each Resource
func TestNewModule_resource_providers(t *testing.T) {
	cfg, diags := testNestedModuleConfigFromDir(t, "testdata/valid-modules/nested-providers-fqns")
	if diags.HasErrors() {
		t.Fatal(diags.Error())
	}

	// both the root and child module have two resources, one which should use
	// the default implied provider and one explicitly using a provider set in
	// required_providers
	wantImplicit := addrs.NewLegacyProvider("test")
	wantFoo := addrs.NewLegacyProvider("foo")
	wantBar := addrs.NewLegacyProvider("bar")

	// root module
	if !cfg.Module.ManagedResources["test_instance.explicit"].Provider.Equals(wantFoo) {
		t.Fatalf("wrong provider for \"test_instance.explicit\"\ngot:  %s\nwant: %s",
			cfg.Module.ManagedResources["test_instance.explicit"].Provider,
			wantFoo,
		)
	}
	if !cfg.Module.ManagedResources["test_instance.implicit"].Provider.Equals(wantImplicit) {
		t.Fatalf("wrong provider for \"test_instance.implicit\"\ngot:  %s\nwant: %s",
			cfg.Module.ManagedResources["test_instance.implicit"].Provider,
			wantImplicit,
		)
	}

	// a data source
	if !cfg.Module.DataResources["data.test_resource.explicit"].Provider.Equals(wantFoo) {
		t.Fatalf("wrong provider for \"module.child.test_instance.explicit\"\ngot:  %s\nwant: %s",
			cfg.Module.ManagedResources["test_instance.explicit"].Provider,
			wantBar,
		)
	}

	// child module
	cm := cfg.Children["child"].Module
	if !cm.ManagedResources["test_instance.explicit"].Provider.Equals(wantBar) {
		t.Fatalf("wrong provider for \"module.child.test_instance.explicit\"\ngot:  %s\nwant: %s",
			cfg.Module.ManagedResources["test_instance.explicit"].Provider,
			wantBar,
		)
	}
	if !cm.ManagedResources["test_instance.implicit"].Provider.Equals(wantImplicit) {
		t.Fatalf("wrong provider for \"module.child.test_instance.implicit\"\ngot:  %s\nwant: %s",
			cfg.Module.ManagedResources["test_instance.implicit"].Provider,
			wantImplicit,
		)
	}
}

func TestProviderForLocalConfig(t *testing.T) {
	mod, diags := testModuleFromDir("testdata/providers-explicit-fqn")
	if diags.HasErrors() {
		t.Fatal(diags.Error())
	}
	lc := addrs.LocalProviderConfig{LocalName: "foo-test"}
	got := mod.ProviderForLocalConfig(lc)
	want := addrs.NewLegacyProvider("foo")
	if !got.Equals(want) {
		t.Fatalf("wrong result! got %#v, want %#v\n", got, want)
	}
}
