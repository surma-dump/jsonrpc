JsonRPC v1.1
============

JsonRPC is a package for Go which enables RPCs via a generic JSON formatted protocol.
JsonRPC is medium-oblivious - it doesn’t matter how you obtain the JSON strings (uuuhh).

Details
-------
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

The return value(s) will be an marshalled into an array.

Pitfalls
--------
Due to limitations of Go’s reflect package (or rather the language design), all method
parameters have to be pointer types (or interfaces).

Credits
-------
Alexander Surma <alexander.surma@gmail.com>
