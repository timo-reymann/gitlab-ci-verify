package rego_policies

import (
	_ "embed"
	"github.com/open-policy-agent/opa/v1/ast"
	"github.com/open-policy-agent/opa/v1/bundle"
	"github.com/open-policy-agent/opa/v1/loader"
	"github.com/open-policy-agent/opa/v1/rego"
	"github.com/open-policy-agent/opa/v1/types"
	"github.com/timo-reymann/gitlab-ci-verify/internal/httputils"
	"io"
	"os"
	"strings"
)

//go:embed lib.rego
var regoLib string

// NewRegoPolicyManager creates a new rego instance to manage policies
func NewRegoPolicyManager() *RegoPolicyManager {
	fileLoader := loader.NewFileLoader().WithProcessAnnotation(true)

	mgr := &RegoPolicyManager{
		fileLoader: fileLoader,
		options: []func(r *rego.Rego){
			rego.EnablePrintStatements(true),
			rego.PrintHook(logPrinter{}),
			rego.Query("data[_].findings"),
		},
	}

	_ = mgr.LoadModuleFromString("__gitlab_ci_verify_lib.rego", regoLib)

	return mgr
}

// RegoPolicyManager contains everything required to create and manage rego contexts
type RegoPolicyManager struct {
	fileLoader loader.FileLoader
	options    []func(r *rego.Rego)
}

// LoadBundle from a folder. It searches the entire folder and includes
// all .rego files that are in the top level
func (rpm *RegoPolicyManager) LoadBundle(path string) error {
	loadedModules, err := loader.NewFileLoader().WithProcessAnnotation(true).Filtered([]string{path}, func(_ string, info os.FileInfo, _ int) bool {
		return !info.IsDir() && !strings.HasSuffix(info.Name(), bundle.RegoExt)
	})
	if err != nil {
		return err
	}

	modules := loadedModules.ParsedModules()

	for _, v := range modules {
		rpm.options = append(rpm.options, rego.ParsedModule(v))
	}

	return nil
}

func (rpm *RegoPolicyManager) LoadBundleFromRemote(url string) error {
	hc := httputils.NewRfc7232HttpClient()
	bundleReader, err := hc.ReadRemoteOrCached(url)
	if err != nil {
		return err
	}
	tl := bundle.NewTarballLoaderWithBaseURL(bundleReader, url)
	br := bundle.NewCustomReader(tl)
	parsedBundle, err := br.Read()
	if err != nil {
		return err
	}

	for _, module := range parsedBundle.Modules {
		rpm.options = append(rpm.options, rego.ParsedModule(module.Parsed))
	}

	return nil
}

// LoadModuleFromString given path and content
func (rpm *RegoPolicyManager) LoadModuleFromString(path string, content string) error {
	module, err := ast.ParseModule(path, content)
	if err != nil {
		return err
	}

	rpm.options = append(rpm.options, rego.ParsedModule(module))

	return nil
}

// LoadModuleFromFile given path and read contents
func (rpm *RegoPolicyManager) LoadModuleFromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	code, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	return rpm.LoadModuleFromString(path, string(code))
}

// NewRegoCtx creates a new rego instanced based on manager configuration
func (rpm *RegoPolicyManager) NewRegoCtx() *rego.Rego {
	return rego.New(rpm.options...)
}

// AddBuiltinFunc1 for rego, callable by policies
func (rpm *RegoPolicyManager) AddBuiltinFunc1(
	name string,
	signature *types.Function,
	impl rego.Builtin1,
) {
	f := rego.Function1(
		&rego.Function{
			Name:             name,
			Decl:             signature,
			Memoize:          true,
			Nondeterministic: true,
		},
		impl,
	)
	rpm.options = append(rpm.options, f)
}

// AddBuiltinFunc2 for rego, callable by policies
func (rpm *RegoPolicyManager) AddBuiltinFunc2(
	name string,
	signature *types.Function,
	impl rego.Builtin2,
) {
	f := rego.Function2(
		&rego.Function{
			Name:             name,
			Decl:             signature,
			Memoize:          true,
			Nondeterministic: true,
		},
		impl,
	)
	rpm.options = append(rpm.options, f)
}
