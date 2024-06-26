package generator

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/zarldev/ffakes/pkg/app/info"
)

// ParseAndGenerate parses the given filename and generates the fakes
func ParseAndGenerate(input string, interfaceNames []string, output string) error {
	slog.Debug("parsing file", slog.String("filename", input))
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, input, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file while generating fakes: %w", err)
	}
	packageName := getPackageName(node)
	slog.Debug("package name", slog.String("name", packageName))
	interfaces := getInterfaceInfo(node, interfaceNames)
	igi := InterfaceGenInfo{
		filename:   input,
		pkg:        packageName,
		interfaces: interfaces,
	}
	slog.Debug("generating fakes", slog.String("filename", igi.Filename()))
	err = generateFakes(igi)
	if err != nil {
		return fmt.Errorf("failed to generate fakes: %w", err)
	}
	return nil
}

type InterfaceGenInfo struct {
	filename   string
	pkg        string
	output     string
	interfaces []InterfaceInfo
}

func (i InterfaceGenInfo) Filename() string {
	lowerName := strings.ToLower(i.filename)
	remove := strings.NewReplacer(".go", "")
	lowerName = remove.Replace(lowerName)
	fname := fmt.Sprintf("%s_fakes.go", lowerName)
	if i.output != "" {
		fname = i.output
	}
	return fname
}

func generateFakes(i InterfaceGenInfo) error {
	if len(i.interfaces) == 0 {
		slog.Debug("no interfaces to generate fakes for")
		return nil
	}

	f, err := os.Create(i.Filename())
	if err != nil {
		fmt.Println("Failed to create file", err)
		return err
	}
	slog.Debug("file created", slog.String("filename", i.Filename()))
	w := io.StringWriter(f)
	defer f.Close()
	writeAll(w, i)
	slog.Debug("file written", slog.String("filename", i.Filename()))
	err = formatFile(i.Filename())
	if err != nil {
		fmt.Println("Failed to format file", err)
	}
	slog.Debug("file formatted", slog.String("filename", i.Filename()))
	slog.Debug("total fakes generated", slog.Int("count", len(i.interfaces)))
	return nil
}

func writeAll(w io.StringWriter, igi InterfaceGenInfo) {
	writeGeneratedComment(w)
	writeHeader(w, igi)
	for _, i := range igi.interfaces {
		writeInterface(w, i)
	}
}

func writeGeneratedComment(w io.StringWriter) {
	write(w, fmt.Sprintf("// Code generated by ffakes %s DO NOT EDIT.\n\n", info.AppInfo.Version()))
}

func writeInterface(w io.StringWriter, i InterfaceInfo) {
	writeFakeStruct(w, i)
	writeFuncDefinitions(w, i)
	writeOption(w, i)
	writeOptionFuncs(w, i)
	writeStructOptionFuncs(w, i)
	writeNewFunc(w, i)
	writeFakeMethods(w, i)
}

func write(w io.StringWriter, strs ...string) {
	for i, s := range strs {
		_, err := w.WriteString(s)
		if i != len(strs)-1 {
			_, err = w.WriteString(", ")
		}
		if err != nil {
			fmt.Println("Failed to write string", err)
		}
	}
}

func writeFuncDefinitions(w io.StringWriter, i InterfaceInfo) {
	for _, m := range i.Methods {
		if m.Composite {
			continue
		}
		write(w, fmt.Sprintf("type %sFunc = func(", m.Name))
		for i, p := range m.Params {
			write(w, fmt.Sprintf("%s %s", paramString(p), p.TypeString()))
			if i != len(m.Params)-1 {
				write(w, ", ")
			}
		}
		write(w, ") (")
		for i, r := range m.Results {
			write(w, fmt.Sprintf("%s %s", paramString(r), r.TypeString()))
			if i != len(m.Results)-1 {
				write(w, ", ")
			}
		}
		write(w, ")\n")
	}
}

func writeNewFunc(w io.StringWriter, i InterfaceInfo) {
	write(w, fmt.Sprintf("func NewFake%s(t *testing.T, opts ...%sOption) *Fake%s {\n", i.Name, i.Name, i.Name))
	write(w, fmt.Sprintf("\tf := &Fake%s{ t: t }\n", i.Name))
	write(w, "\tfor _, opt := range opts {\n")
	write(w, "\t\topt(f)\n")
	write(w, "\t}\n")
	writeAssertExpectations(w, i)
	write(w, "\treturn f\n")
	write(w, "}\n\n")
}

func writeAssertExpectations(w io.StringWriter, i InterfaceInfo) {
	write(w, "t.Cleanup(func() {\n")
	for _, m := range i.Methods {
		write(w, fmt.Sprintf("\tif f.%sCount != len(f.F%s) {\n", m.Name, m.Name))
		write(w, fmt.Sprintf("\t\tt.Fatalf(\"expected %s to be called %%d times but got %%d\", len(f.F%s), f.%sCount)\n", m.Name, m.Name, m.Name))
		write(w, "\t}\n")
	}
	write(w, "})\n")
}

func writeOption(w io.StringWriter, i InterfaceInfo) {
	write(w, fmt.Sprintf("type %sOption func(f *", i.Name))
	write(w, fmt.Sprintf("Fake%s)\n\n", i.Name))
}

func writeOptionFuncs(w io.StringWriter, i InterfaceInfo) {
	iname := i.Name
	prefix := ""
	for _, m := range i.Methods {
		if i.IsComposite {
			prefix = iname
		}
		write(w, fmt.Sprintf("func %sOn%s(fn ...%sFunc) %sOption {\n", prefix, m.Name, m.Name, iname))
		write(w, fmt.Sprintf("\treturn func(f *Fake%s) {\n", i.Name))
		write(w, fmt.Sprintf("\t\tf.F%s = append(f.F%s, fn...)\n", m.Name, m.Name))
		write(w, "\t}\n")
		write(w, "}\n\n")
	}
}

func writeStructOptionFuncs(w io.StringWriter, i InterfaceInfo) {
	iname := i.Name
	for _, m := range i.Methods {
		write(w, fmt.Sprintf("func (f *Fake%s) On%s(fns ...%sFunc) {\n", iname, m.Name, m.Name))
		write(w, "for _, fn := range fns {\n")
		write(w, fmt.Sprintf("\tf.F%s = append(f.F%s, fn)\n", m.Name, m.Name))
		write(w, "}\n")
		write(w, "}\n\n")
	}
}

func writeFakeStruct(w io.StringWriter, i InterfaceInfo) {
	write(w, fmt.Sprintf("type Fake%s struct {\n", i.Name))
	write(w, "\t t *testing.T\n")
	for _, m := range i.Methods {
		write(w, fmt.Sprintf("\t%sCount int\n", m.Name))
	}
	for _, m := range i.Methods {
		write(w, fmt.Sprintf("\tF%s []func(", m.Name))
		for i, p := range m.Params {
			write(w, fmt.Sprintf("%s %s", paramString(p), p.TypeString()))
			if i != len(m.Params)-1 {
				write(w, ", ")
			}
		}
		write(w, ")")
		if len(m.Results) > 1 {
			write(w, " (")
		}
		for i, r := range m.Results {
			write(w, fmt.Sprintf("%s %s", paramString(r), r.TypeString()))
			if i != len(m.Results)-1 {
				write(w, ", ")
			}
		}
		if len(m.Results) > 1 {
			write(w, ")")
		}
		write(w, "\n")
	}
	write(w, "}\n\n")
}

func writeFakeMethods(w io.StringWriter, i InterfaceInfo) {
	for _, m := range i.Methods {
		writeFakeMethod(w, i.Name, m)
	}
}

func writeFakeMethod(w io.StringWriter, iName string, m FunctionInfo) {
	write(w, fmt.Sprintf("func (fake *Fake%s) %s(", iName, m.Name))
	for i, p := range m.Params {
		if len(p.Name) != 0 {
			write(w, fmt.Sprintf("%s %s", paramString(p), p.TypeString()))
			if i != len(m.Params)-1 {
				write(w, ", ")
			}
		} else {
			in := "i" + fmt.Sprint(i+1)
			write(w, fmt.Sprintf("%s %s", in, p.TypeString()))
			if i != len(m.Params)-1 {
				write(w, ", ")
			}
		}
		if i != len(m.Params)-1 {
			write(w, ", ")
		}
	}
	write(w, ") (")
	for i, r := range m.Results {
		if i != 0 {
			write(w, ", ")
		}
		if len(r.Name) == 0 {
			write(w, r.TypeString())
			continue
		}
		write(w, fmt.Sprintf("%s %s", paramString(r), r.TypeString()))
	}
	write(w, ") {\n")
	write(w, fmt.Sprintf("\t var idx = fake.%sCount\n", m.Name))
	write(w, fmt.Sprintf("\tif fake.%sCount >= len(fake.F%s) {\n", m.Name, m.Name))
	write(w, fmt.Sprintf("\t idx = len(fake.F%s) - 1\n", m.Name))
	write(w, "\t}\n")
	// write(w, fmt.Sprintf("\tif f.%sCount >= len(f.F%s) {\n", m.Name, m.Name))
	// write(w, "\t f.t.Fatalf(\"too many calls to "+m.Name+" expected %d calls got %d \", "+fmt.Sprintf("f.%sCount, f.%sCount + 1)", m.Name, m.Name)+"\n")
	// write(w, "\t}\n")
	write(w, "\tif len(fake.F")
	write(w, m.Name)
	write(w, ") != 0 {\n")
	for i, p := range m.Results {
		if i != 0 {
			write(w, ", ")
		}
		if len(p.Name) == 0 {
			write(w, fmt.Sprintf("o%d", i+1))
			continue
		}
		write(w, p.Name...)
	}
	write(w, " := fake.F")
	write(w, m.Name)
	write(w, "[idx](")
	for i, p := range m.Params {
		if i != 0 {
			write(w, ", ")
		}
		for i, name := range p.Name {
			if name == "" {
				write(w, "i"+fmt.Sprint(i+1))
				continue
			}
		}
		write(w, p.Name...)
	}
	write(w, ")\n")
	write(w, fmt.Sprintf("\tfake.%sCount++\n", m.Name))
	write(w, "return ")
	for i, p := range m.Results {
		if i != 0 {
			write(w, ", ")
		}
		if len(p.Name) == 0 {
			write(w, fmt.Sprintf("o%d", i+1))
			continue
		}
		write(w, p.Name...)
	}
	write(w, "}\n")
	write(w, "return ")
	for i, r := range m.Results {
		write(w, r.ZeroValue())
		if i != len(m.Results)-1 {
			write(w, ", ")
		}
	}
	write(w, "}\n\n")
}

func paramString(p ParamInfo) string {
	n := strings.Builder{}
	for i, name := range p.Name {
		n.WriteString(name)
		if i != len(p.Name)-1 {
			n.WriteString(", ")
		}
	}
	return n.String()
}

func writeHeader(w io.StringWriter, igi InterfaceGenInfo) {
	write(w, fmt.Sprintf("package %s\n\n", igi.pkg))
	write(w, "import (\n")
	write(w, "\t\"testing\"\n")
	for _, i := range igi.interfaces {
		for _, imp := range i.Imports {
			write(w, fmt.Sprintf("\t\"%s\"\n", imp))
		}
	}

	write(w, ")\n\n")
}

func formatFile(filename string) error {
	f, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	b, err := format.Source(f)
	if err != nil {
		return fmt.Errorf("failed to format file: %w", err)
	}
	err = os.WriteFile(filename, b, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

type InterfaceInfo struct {
	Name        string
	Methods     []FunctionInfo
	IsComposite bool
	Composites  []string
	Imports     []string
}

type FunctionInfo struct {
	Name      string
	Composite bool
	Params    []ParamInfo
	Results   []ParamInfo
}

type ParamInfo struct {
	Name   []string
	Type   string
	Zero   string
	IsPtr  bool
	IsChan bool
	IsList bool
}

func (p ParamInfo) TypeString() string {
	str := strings.Builder{}
	if p.IsChan {
		str.WriteString("chan ")
	}
	if p.IsList {
		str.WriteString("[]")
	}
	str.WriteString(" ")
	if p.IsPtr {
		str.WriteString("*")
	}
	str.WriteString(p.Type)
	return str.String()
}

func (p ParamInfo) ZeroValue() string {
	if p.IsPtr || p.IsList || p.IsChan {
		return "nil"
	}
	switch p.Type {
	case "string":
		return `""`
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "0"
	case "float32", "float64":
		return "0.0"
	case "bool":
		return "false"
	case "error":
		return "nil"
	default:
		return p.Type + "{}"
	}
}

func getInterfaceInfo(node *ast.File, ifnames []string) []InterfaceInfo {
	var (
		interfaces = make(map[string]InterfaceInfo)
	)
	slog.Debug("getting interface info")
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if tSpec, ok := spec.(*ast.TypeSpec); ok {
					if iType, ok := tSpec.Type.(*ast.InterfaceType); ok {
						if !slices.Contains(ifnames, tSpec.Name.Name) {
							continue
						}
						iInfo := InterfaceInfo{
							Name:       tSpec.Name.Name,
							Composites: []string{},
						}
						for _, field := range iType.Methods.List {
							if len(field.Names) == 0 {
								val, importStr, _ := fieldValueAndIsPointer(field)
								if importStr != "" {
									iInfo.Imports = append(iInfo.Imports, importStr)
								}
								iInfo.Composites = append(iInfo.Composites, val)
								iInfo.IsComposite = true
								continue
							}
							var info FunctionInfo
							if len(field.Names) > 0 {
								info = FunctionInfo{
									Name: field.Names[0].Name,
								}
							}
							if f, ok := field.Type.(*ast.FuncType); ok {
								for _, param := range f.Params.List {
									params := make([]ParamInfo, 0, len(param.Names))
									if param, importStr, ok := paramInfo(param); ok {
										params = append(params, param)
										info.Params = params
										if importStr != "" {
											iInfo.Imports = append(iInfo.Imports, importStr)
										}
									}
								}
								for _, result := range f.Results.List {
									var (
										isChan, isPtr             bool
										name, importStr, typeName string
									)
									paramType, ok := result.Type.(*ast.Ident)
									if !ok {
										if chanResult, importStr, ok := paramInfo(result); ok {
											info.Results = append(info.Results, chanResult)
											if importStr != "" {
												iInfo.Imports = append(iInfo.Imports, importStr)
											}

											continue
										}
										if typeName, importStr, isPtr = fieldValueAndIsPointer(result); isPtr || importStr != "" {
											iInfo.Imports = append(iInfo.Imports, importStr)
											if len(result.Names) == 0 {
												info.Results = append(info.Results, ParamInfo{
													Name:   []string{},
													Type:   typeName,
													IsPtr:  isPtr,
													IsChan: isChan,
												})
											}
										}
										continue
									}
									buildResults(result, &info, ParamInfo{
										Name:   []string{name},
										Type:   paramType.Name,
										IsPtr:  isPtr,
										IsChan: isChan,
									})
								}
							}
							iInfo.Methods = append(iInfo.Methods, info)
						}
						slog.Debug("adding interface info")
						slog.Debug("interface info", slog.String("name", iInfo.Name))
						slog.Debug("interface methods", slog.String("methods", fmt.Sprint(iInfo.Methods)))
						slog.Debug("composite", slog.Bool("composite", iInfo.IsComposite))
						interfaces[iInfo.Name] = iInfo
					}
				}
			}
		}
	}
	for _, i := range interfaces {
		if i.IsComposite {
			for _, c := range i.Composites {
				methods := interfaces[c].Methods
				meths := slices.Clone(methods)
				for _, m := range meths {
					m.Composite = true
					i.Methods = append(i.Methods, m)
				}
			}
			interfaces[i.Name] = i
		}
	}
	info := make([]InterfaceInfo, 0, len(interfaces))
	for _, i := range interfaces {
		info = append(info, i)
	}
	slog.Debug("total interfaces", slog.Int("count", len(info)))
	return info
}

func buildList(field *ast.Field, list []ParamInfo, param ParamInfo) []ParamInfo {
	r := slices.Clone(list)
	if len(field.Names) == 0 {
		r = append(r, ParamInfo{
			Name:   param.Name,
			Type:   param.Type,
			IsPtr:  param.IsPtr,
			IsChan: param.IsChan,
		})
	}
	for _, paramName := range field.Names {
		r = append(r, ParamInfo{
			Name:   []string{paramName.Name},
			Type:   param.Type,
			IsPtr:  param.IsPtr,
			IsChan: param.IsChan,
		})
	}
	return r
}

func buildResults(field *ast.Field, info *FunctionInfo, param ParamInfo) {
	param.Name = []string{}
	info.Results = buildList(field, info.Results, param)
}

func paramInfo(param *ast.Field) (ParamInfo, string, bool) {
	var (
		typeName, importStr   string
		isPtr, isChan, isList bool
		names                 []string
	)
	if len(param.Names) > 0 {
		for _, n := range param.Names {
			names = append(names, n.Name)
		}
	}
	isChan = isChannel(param)
	isList = isListType(param)
	typeName, importStr, isPtr = valAndIsPointer(param)
	p := ParamInfo{
		Name:   names,
		Type:   typeName,
		IsPtr:  isPtr,
		IsChan: isChan,
		IsList: isList,
	}
	return p, importStr, true
}

func isListType(param *ast.Field) bool {
	if _, ok := param.Type.(*ast.ArrayType); ok {
		return true
	}
	return false
}

func isChannel(param *ast.Field) bool {
	if _, ok := param.Type.(*ast.ChanType); ok {
		return true
	}
	return false
}

func valAndIsPointerExpr(expr ast.Expr) (string, string, bool) {
	var (
		isPtr    bool
		typeName string
		importS  string
	)
	if val, ok := expr.(*ast.Ident); ok {
		isPtr = false
		typeName = val.Name
	}
	if ptr, ok := expr.(*ast.StarExpr); ok {
		isPtr = true
		if valX, ok := ptr.X.(*ast.Ident); ok {
			typeName = valX.Name
		}
		if valX, ok := ptr.X.(*ast.SelectorExpr); ok {
			typeName = valX.Sel.Name
		}
	}
	if sel, ok := expr.(*ast.SelectorExpr); ok {
		isPtr = false
		switch x := sel.X.(type) {
		case *ast.Ident:
			importS = x.Name
			typeName = x.Name + "." + sel.Sel.Name
		case *ast.SelectorExpr:
			typeName = x.Sel.Name + "." + sel.Sel.Name
		}
	}
	return typeName, importS, isPtr
}
func valAndIsPointer(field *ast.Field) (string, string, bool) {
	switch t := field.Type.(type) {
	case *ast.Ident, *ast.SelectorExpr, *ast.StarExpr, *ast.FuncType, *ast.InterfaceType, *ast.StructType:
		return valAndIsPointerExpr(t)
	case *ast.ChanType:
		return valAndIsPointerExpr(t.Value)
	case *ast.ArrayType:
		return valAndIsPointerExpr(t.Elt)
	case *ast.MapType:
		return valAndIsPointerExpr(t.Value)
	default:
		return "", "", false
	}
}

func fieldValueAndIsPointer(field *ast.Field) (string, string, bool) {
	return valAndIsPointer(field)
}

func getPackageName(node *ast.File) string {
	var packageName string
	if node.Name != nil {
		packageName = node.Name.Name
	}
	return packageName
}
