/// Abstractions for PKL records
@ModuleInfo { minPklVersion = "0.28.2" }

@go.Package { name = "github.com/kdeps/schema/gen/pkl_resource" }

open module org.kdeps.pkl.PklResource

import "package://pkg.pkl-lang.org/pkl-go/pkl.golang@0.10.0#/go.pkl"
import "package://pkg.pkl-lang.org/pkl-pantry/pkl.experimental.uri@1.0.3#/URI.pkl"

/// Retrieves a session record by its [id]
///
/// Returns the textual content of the session entry, or an empty string if not found.
///
/// [id]: The identifier of the session record.
function getPklRecord(id: String, typ: String): String = read("pklres:/\(id)?type=\(typ)")?.text ?? ""

/// Sets or updates a session record with a new [value]
///
/// Returns the set value as confirmation.
///
/// [id]: The identifier of the session record.
/// [value]: The value to store.
function setPklValue(id: String, typ: String, key: String, value: String): String = read("pklres:/\(id)?type=\(typ)?key=\(key)?op=set&value=\(URI.encodeComponent(value))")?.text ?? ""

/// Get a pkl key record
///
/// Returns the key value
///
/// [id]: The identifier of the session record.
/// [value]: The value to store.
function getPklValue(id: String, typ: String, key: String): String = read("pklres:/\(id)?type=\(typ)?key=\(key)?op=get")?.text ?? ""
