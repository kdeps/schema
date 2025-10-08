module github.com/kdeps/schema

go 1.23.7

toolchain go1.24.4

require github.com/apple/pkl-go v0.11.1

retract (
	v0.3.0
	v0.3.1
	v0.3.2
	v0.3.3
	v0.3.4
	v0.3.5
	v0.3.6
	v0.3.7
	v0.3.8
	v0.4.0
	v0.4.1
	v0.4.2
	v0.4.3
	v0.4.4
	v0.4.5
	v0.4.6
)

require (
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
)
