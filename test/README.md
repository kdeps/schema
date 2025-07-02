# PKL Function Test Suite

This directory contains comprehensive tests for the KDEPS PKL codebase to validate all camelCase function implementations, embedded PKL assets functionality, and ensure proper functionality across both PKL and Go environments.

## 📁 Test Files

| File | Purpose | Coverage |
|------|---------|----------|
| `test_functions.pkl` | Comprehensive PKL integration tests | All modules and camelCase functions |
| `test_utils.pkl` | Focused PKL unit tests | Utils.pkl module (example) |
| `assets_test.go` | **NEW** Go-based assets tests | PKL embedding, workspace setup, tmpdir handling |
| `test_new_attributes.pkl` | Focused new attributes testing | ValidationCheck, DockerSettings, Project rate limiting |

## 🎯 Test Coverage

### ✅ PKL Functions Tested

The test suite validates all camelCase functions across the PKL modules:

#### **Utils.pkl**
- `isBase64(str: String) -> Boolean`

#### **Memory.pkl & Session.pkl**
- `getRecord(id: String) -> String`
- `setRecord(id: String, value: String) -> String`
- `deleteRecord(id: String) -> String`
- `clear() -> String`

#### **Tool.pkl**
- `getOutput(id: String) -> String`
- `runScript(id: String, script: String, params: String) -> String`
- `history(id: String) -> String`

#### **Data.pkl**
- `filepath(agentName: String, fileName: String) -> String`

#### **Item.pkl**
- `current() -> String`
- `prev() -> String`
- `next() -> String`
- `values(id: String) -> Listing<String>`

#### **Skip.pkl**
- `ifFileExists(it: String) -> Boolean`
- `ifFolderExists(it: String) -> Boolean`
- `ifFileIsEmpty(it: String) -> Boolean`

#### **Document.pkl**
- `jsonParser(data: String) -> Any`
- `jsonParserMapping(data: String) -> Any`
- `jsonRenderDocument(value: Any) -> String`
- `jsonRenderValue(value: Any) -> String`
- `yamlRenderDocument(value: Any) -> String`
- `yamlRenderValue(value: Any) -> String`
- `xmlRenderDocument(value: Any) -> String`
- `xmlRenderValue(value: Any) -> String`

#### **LLM.pkl**
- `resource(actionID: String) -> ResourceChat`
- `response(actionID: String) -> String`
- `prompt(actionID: String) -> String`
- `jsonResponse(actionID: String) -> Boolean`
- `jsonResponseKeys(actionID: String) -> Listing<String>`
- `itemValues(actionID: String) -> Listing<String>`
- `file(actionID: String) -> String`

#### **Exec.pkl & Python.pkl**
- `resource(actionID: String) -> ResourceExec/ResourcePython`
- `stderr(actionID: String) -> String`
- `stdout(actionID: String) -> String`
- `exitCode(actionID: String) -> Int`
- `file(actionID: String) -> String`
- `itemValues(actionID: String) -> Listing<String>`
- `env(actionID: String, envName: String) -> String`

#### **HTTP.pkl**
- `resource(actionID: String) -> ResourceHTTPClient`
- `responseBody(actionID: String) -> String`
- `file(actionID: String) -> String`
- `itemValues(actionID: String) -> Listing<String>`
- `responseHeader(actionID: String, headerActionID: String) -> String`

#### **APIServerRequest.pkl**
- `data() -> String`
- `params(name: String) -> String`
- `header(name: String) -> String`
- `file(name: String) -> APIServerRequestUploads`
- `filetype(name: String) -> String`
- `filepath(name: String) -> String`
- `filecount() -> String`
- `fileList() -> Listing`
- `filetypes() -> Listing`
- `filesByType(mimeType: String) -> Listing`
- `path() -> String`
- `method() -> String`
- `ip() -> String`
- `id() -> String`

### ✅ **NEW** Assets Package Functionality

The Go test suite (`assets_test.go`) validates the embedded PKL assets system:

#### **PKL File Embedding**
- `PKLFileExists(filename) -> bool` - Verify embedded files
- `GetPKLFile(filename) -> []byte` - Read embedded PKL files
- `GetPKLFileAsString(filename) -> string` - Get PKL content as string
- `ListPKLFiles() -> []string` - List all embedded files
- `ValidatePKLFiles() -> error` - Ensure all expected files present

#### **TmpDir Operations**
- `GetTmpDir() -> string` - Get system tmpdir path
- `ExtractPKLFileToTemp(filename) -> string` - Single file extraction
- `ExtractPKLFileWithName(filename, dir) -> string` - **Preserve original names**
- `ExtractAllPKLFilesToDir(dir) -> string` - Extract all with names preserved

#### **PKL Workspace for Testing**
- `SetupPKLWorkspace(dir) -> *PKLWorkspace` - General workspace setup
- `SetupPKLWorkspaceInTmpDir() -> *PKLWorkspace` - **Explicit tmpdir workspace**
- `workspace.GetImportPath(filename) -> string` - **Absolute paths for PKL imports**
- `workspace.IsInTmpDir() -> bool` - Verify tmpdir location
- `workspace.IsTemporary() -> bool` - Check if auto-cleanup enabled
- `workspace.Cleanup() -> error` - **Clean up temp files**

## 🚀 Running Tests

### Using Make (Recommended)

```bash
# Run PKL function tests
make test

# Run Utils unit tests  
make test-utils

# Run all PKL tests
make test-all

# NEW: Run Go assets tests
make test-assets

# NEW: Run all tests (PKL + Go)
make test-all-comprehensive

# Run tests and regenerate Go code
make test-and-generate

# Show help
make help
```

### Direct Test Execution

```bash
# PKL function tests
pkl eval test/test_functions.pkl
pkl eval test/test_utils.pkl

# Go assets tests
cd test && go test -v .

# Go tests with benchmarks  
cd test && go test -bench=. -v .
```

## 📋 Test Structure

### PKL Tests (`test_functions.pkl`, `test_utils.pkl`)
- **Integration testing approach**
- **Function existence validation**
- **Default value verification**
- **Cross-module compatibility**
- **Configuration validation**

### **NEW** Go Assets Tests (`assets_test.go`)
- **PKL file embedding validation**
- **Workspace setup for external testing**
- **TmpDir operations and cleanup**
- **File name preservation verification**
- **Performance benchmarking**

## 🧪 **NEW** External Testing with Assets

For tests in external repositories that need to import PKL schemas:

```go
func TestMyPKLWorkflow(t *testing.T) {
    // Setup PKL workspace with all schema files
    workspace, err := assets.SetupPKLWorkspaceInTmpDir()
    if err != nil {
        t.Fatalf("Failed to setup PKL workspace: %v", err)
    }
    defer workspace.Cleanup() // Important: clean up temp files

    // Get absolute paths for PKL imports
    workflowPath := workspace.GetImportPath("Workflow.pkl")
    resourcePath := workspace.GetImportPath("Resource.pkl")
    
    // Create test PKL file with absolute imports
    testContent := fmt.Sprintf(`
        import "%s" as Workflow
        import "%s" as Resource
        
        myWorkflow = new Workflow {
            AgentID = "test-agent"
            Description = "Test workflow"
            Version = "1.0.0"
            TargetActionID = "test-action"
            Workflows {}
            Settings = new Workflow.Project.Settings {
                AgentSettings = new Workflow.Docker.DockerSettings {}
                APISettings = new Workflow.APIServer.APIServerSettings {}
            }
        }
    `, workflowPath, resourcePath)
    
    // Write and evaluate test PKL file
    testFile := filepath.Join(workspace.Directory, "test.pkl")
    err = os.WriteFile(testFile, []byte(testContent), 0644)
    // ... rest of test
}
```

## ✅ Test Results

All tests validate:

### PKL Function Tests
1. **Function Existence**: All camelCase functions are accessible
2. **Default Values**: Resource functions return correct defaults
3. **Base64 Handling**: Proper encoding/decoding behavior
4. **File Operations**: File existence and folder validation
5. **Configuration**: Default settings and validation functions
6. **Data Processing**: JSON/YAML/XML parsing and rendering

### **NEW** Assets Tests
1. **Embedding Validation**: All 21 PKL files properly embedded
2. **File Name Preservation**: Original names maintained (e.g., "Workflow.pkl")
3. **TmpDir Operations**: Proper system tmpdir usage
4. **Workspace Setup**: Functional PKL import paths
5. **Cleanup Operations**: No temp file leaks
6. **Cross-Import Compatibility**: All schema cross-references work

## 🔧 Extending Tests

### For PKL Function Testing
Add facts to `test_functions.pkl` or create new module-specific files like `test_module.pkl`

### **NEW** For Assets Testing
Add test functions to `assets_test.go`:

```go
func TestNewAssetFeature(t *testing.T) {
    // Test new assets functionality
    workspace, err := assets.SetupPKLWorkspaceInTmpDir()
    if err != nil {
        t.Fatalf("Setup failed: %v", err)
    }
    defer workspace.Cleanup()
    
    // Your test assertions here
}
```

## 🎯 Validation Strategy

The test suite follows a **comprehensive validation approach**:

### PKL Level
1. **Syntax Validation**: All PKL files parse correctly
2. **Function Accessibility**: camelCase functions are callable
3. **Type Safety**: Functions return expected types
4. **Error Handling**: Graceful handling of edge cases
5. **Integration**: Cross-module compatibility

### **NEW** Go Level
1. **Embedding Verification**: PKL files accessible in binaries
2. **Workspace Functionality**: External testing capability
3. **File Operations**: Proper tmpdir usage and cleanup
4. **Name Preservation**: Original filenames maintained
5. **Performance**: Benchmark key operations

## 📈 Coverage Goals

- ✅ **100% PKL Function Coverage**: All camelCase functions tested
- ✅ **Default Value Testing**: Resource defaults validated
- ✅ **Edge Case Handling**: Error conditions covered
- ✅ **Integration Testing**: Module interactions verified
- ✅ **Configuration Testing**: Settings and validation confirmed
- ✅ **NEW: 100% Assets Coverage**: All embedding functionality tested
- ✅ **NEW: External Testing**: Workspace setup for other repos validated
- ✅ **NEW: TmpDir Compliance**: Proper system tmpdir usage verified

---

**🎉 All camelCase function names and embedded assets are working correctly!**

This comprehensive test suite ensures the KDEPS PKL codebase maintains consistent functionality after the PascalCase → camelCase function naming conversion, while providing robust embedded assets for external testing and integration.

### **🆕 Key Benefits of New Assets Testing:**

1. **External Integration**: Other repos can easily test with PKL schemas
2. **Predictable Paths**: TmpDir usage ensures consistent test environments  
3. **Name Preservation**: Original PKL filenames maintained for easy reference
4. **Cross-Import Support**: All schema cross-references work in extracted workspaces
5. **Clean Testing**: Automatic cleanup prevents temp file accumulation 

### 🆕 **New Attributes Tested** (8+ new features)
- **ValidationCheck**: `Retry` (Boolean), `RetryTimes` (Int) - Retry functionality for validation checks
- **ResourceAction**: `PostflightCheck` (ValidationCheck?) - Post-execution validation
- **DockerSettings**: `ExposedPorts` (Listing<String>?) - Docker port configuration
- **DockerSettings**: Updated `OllamaVersion` to "0.9.2"
- **Project Settings**: `RateLimitMax` (Int = 100) - Rate limiting configuration
- **Project Settings**: `Environment` (BuildEnv = "dev") - Environment setting
- **BuildEnv**: Type alias supporting "dev" | "prod" values
- **Version Updates**: PKL 0.28.2 and pkl-go 0.10.0 compatibility

## 🚨 Important Notes

- **PKL Version**: All tests use PKL 0.28.2 and pkl-go 0.10.0
- **Naming Convention**: Functions are camelCase, attributes are PascalCase
- **Assets**: PKL files are embedded even with `assets/pkl/` in `.gitignore`
- **External Testing**: Use assets package for external repository integration
- **Backwards Compatibility**: All existing functionality preserved

## 📈 Future Enhancements

The test suite is designed to easily accommodate:
- Additional PKL modules and functions
- New attribute testing patterns
- Enhanced Go integration scenarios
- Performance optimization validation
- Cross-platform compatibility testing

For detailed examples and advanced usage patterns, see the individual test files and the comprehensive Go test suite in `assets_test.go`. 