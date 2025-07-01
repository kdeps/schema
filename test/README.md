# PKL Function Test Suite

This directory contains comprehensive tests for the KDEPS PKL codebase to validate all camelCase function implementations and ensure proper functionality.

## 📁 Test Files

| File | Purpose | Coverage |
|------|---------|----------|
| `test_functions.pkl` | Comprehensive integration tests | All modules and camelCase functions |
| `test_utils.pkl` | Focused unit tests | Utils.pkl module (example) |

## 🎯 Test Coverage

### ✅ Functions Tested

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

## 🚀 Running Tests

### Using Make (Recommended)

```bash
# Run comprehensive test suite
make test

# Run individual module tests
make test-utils

# Run all tests
make test-all

# Run tests and regenerate Go code
make test-and-generate

# Show help
make help
```

### Direct PKL Evaluation

```bash
# Run comprehensive tests
pkl eval test/test_functions.pkl

# Run Utils unit tests
pkl eval test/test_utils.pkl
```

## 📋 Test Structure

### Comprehensive Tests (`test_functions.pkl`)
- **Integration testing approach**
- **Function existence validation**
- **Default value verification**
- **Cross-module compatibility**
- **Configuration validation**

### Unit Tests (`test_utils.pkl`)
- **Focused on individual modules**
- **Comprehensive edge case testing**
- **Detailed function behavior validation**
- **Example of best practices for unit testing**

## ✅ Test Results

All tests validate:
1. **Function Existence**: All camelCase functions are accessible
2. **Default Values**: Resource functions return correct defaults
3. **Base64 Handling**: Proper encoding/decoding behavior
4. **File Operations**: File existence and folder validation
5. **Configuration**: Default settings and validation functions
6. **Data Processing**: JSON/YAML/XML parsing and rendering

## 🔧 Extending Tests

To add new tests:

1. **For comprehensive testing**: Add facts to `test_functions.pkl`
2. **For focused testing**: Create new files like `test_module.pkl`
3. **Update Makefile**: Add new test targets as needed

### Example Test Pattern

```pkl
facts {
    ["Test description"] {
        // Test assertion using pkl:test utilities
        test.catchOrNull(() -> YourModule.yourFunction("param")) != null
    }
    
    ["Function returns expected value"] {
        YourModule.yourFunction("input") == "expected_output"
    }
}
```

## 🎯 Validation Strategy

The test suite follows a **validation-first approach**:

1. **Syntax Validation**: All PKL files parse correctly
2. **Function Accessibility**: camelCase functions are callable
3. **Type Safety**: Functions return expected types
4. **Error Handling**: Graceful handling of edge cases
5. **Integration**: Cross-module compatibility

## 📈 Coverage Goals

- ✅ **100% Function Coverage**: All camelCase functions tested
- ✅ **Default Value Testing**: Resource defaults validated
- ✅ **Edge Case Handling**: Error conditions covered
- ✅ **Integration Testing**: Module interactions verified
- ✅ **Configuration Testing**: Settings and validation confirmed

---

**🎉 All camelCase function names are working correctly!**

This test suite ensures the KDEPS PKL codebase maintains consistent functionality after the PascalCase → camelCase function naming conversion. 