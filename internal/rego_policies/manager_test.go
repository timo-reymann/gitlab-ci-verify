package rego_policies

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/open-policy-agent/opa/v1/ast"
	"github.com/open-policy-agent/opa/v1/rego"
	"github.com/open-policy-agent/opa/v1/types"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestRegoPolicyManager_NewRegoCtx(t *testing.T) {
	rpm := NewRegoPolicyManager()
	rpm.NewRegoCtx()
}

func TestRegoPolicyManager_LoadBundle(t *testing.T) {
	rpm := NewRegoPolicyManager()

	t.Run("valid bundle", func(t *testing.T) {
		err := rpm.LoadBundle("test_data/valid-bundle")
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("invalid bundle", func(t *testing.T) {
		err := rpm.LoadBundle("test_data/invalid-bundle")
		if err == nil {
			t.Fatal(err)
		}
	})
}

func TestRegoPolicyManager_CompilerErrors(t *testing.T) {
	rpm := NewRegoPolicyManager()
	err := rpm.LoadBundle("test_data/compiler-error-bundle")
	if err != nil {
		t.Fatal(err)
	}
	rpm.NewRegoCtx()
}

func TestRegoPolicyManager_AddBuiltinFunc(t *testing.T) {
	rpm := NewRegoPolicyManager()
	rpm.AddBuiltinFunc1("say_hello", types.NewFunction(types.Args(types.S), types.S), func(context rego.BuiltinContext, param1 *ast.Term) (*ast.Term, error) {
		var param1Val string
		if err := ast.As(param1.Value, &param1Val); err != nil {
			return nil, err
		}

		val, err := ast.InterfaceToValue(fmt.Sprintf("Hello %s", param1Val))
		if err != nil {
			return nil, err
		}

		return ast.NewTerm(val), nil
	})
	if err := rpm.LoadBundle("test_data/func-call-bundle"); err != nil {
		t.Fatal(err)
	}

	regoCtx := rpm.NewRegoCtx()

	ctx := context.TODO()
	query, err := regoCtx.PrepareForEval(ctx)
	if err != nil {
		t.Fatal(err)
	}

	results, err := query.Eval(ctx, rego.EvalInput(map[string]any{
		"foo": "bar",
	}))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%v\n", results)
}

func TestRegoPolicyManager_LoadBundleFromRemote(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		file, _ := os.Open("test_data/bundle.tar.gz")

		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/gzip")
		writer.Header().Set("Content-Disposition", "attachment; filename=bundle.tar.gz")
		_, _ = io.Copy(writer, file)
	}))
	defer mock.Close()

	rpm := NewRegoPolicyManager()
	err := rpm.LoadBundleFromRemote(mock.URL)
	if err != nil {
		t.Fatal(err)
	}
	regoCtx := rpm.NewRegoCtx()
	ctx := context.Background()
	query, err := regoCtx.PrepareForEval(ctx)
	if err != nil {
		t.Fatal(err)
	}

	result, err := query.Eval(ctx, rego.EvalInput(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatal("expected 1 result, got ", len(result))
	}
}
