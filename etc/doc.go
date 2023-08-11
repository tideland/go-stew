// Tideland Go Stew - Etc
//
// Copyright (C) 2019-2023 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package etc provides the reading, parsing, and accessing of configuration
// data. Different readers can be passed as sources for the JSON formatted
// input. It is based on the dynaj package and therefore allows to access
// the data by paths. But the access to data can contain a default value and
// if the read data contains macros they will be resolved.
//
//	    {
//	       "global": {
//	           "baseDirectory": "[[$MYAPP_BASEDIR||/var/lib/my-server]]"",
//	           "hostAddress": "localhost:1234",
//	           "maxUsers": 50
//	       },
//	       "services": [
//	           {
//	               "id": "service-a",
//	               "url": "http://[[global::hostAddress]]/service-a",
//	               "directory": "[[global::baseDirectory||.]]/service-a"
//		       }, {
//	               "id": "service-b",
//	               "url": "http://[[global::hostAddress]]/service-b",
//	               "directory": "[[global::baseDirectory||.]]/service-b"
//	           }
//	       ]
//	   }
//
// The macros here are [[<env-or-path>||<default>]]. The first part is
// the name of an environment variable or a path inside the configuration.
// The second part is optional and contains a default value. If the
// environment variable or the path cannot be found the default value
// will be used. If there's no default value the macro will be the programmed
// default value.
//
// The configuration also can be updated and values can be added. Here on
// the writing copy the macros are not resolved. These values are only changed
// in case they are directly overwritten.
package etc // import "tideland.dev/go/stew/etc"

// EOF
