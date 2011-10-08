About
=====

JsonRPC is a package for Go which enables RPCs via a generic JSON formatted protocol.
JsonRPC is medium-oblivious - it doesn’t matter how you obtain the JSON strings.

Details
=======
All the public methods of an arbitrary object can be called with a simple
JSON-protocol which is as follows:

	{
		"MethodName": "Name of the Method",
		"Parameters": [
			"Parameter1",
			2,
			false,
			[
				"Parameter4.1",
				4.2,
			]
		],
	}

Pitfalls
========
Due to limitations of Go’s reflect package, all method parameters have to be pointer
types (or interfaces).
