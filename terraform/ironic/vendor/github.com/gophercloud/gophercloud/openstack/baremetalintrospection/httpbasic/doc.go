/*
Package httpbasic provides support for http_basic bare metal introspection endpoints.

Example of obtaining and using a client:

	client, err := httpbasic.NewBareMetalIntrospectionHTTPBasic(httpbasic.EndpointOpts{
		IronicInspectorEndpoint:     "http://localhost:5050/v1/",
		IronicInspectorUser:         "myUser",
		IronicInspectorUserPassword: "myPassword",
	})
	if err != nil {
		panic(err)
	}

	introspection.GetIntrospectionStatus(client, "a62b8495-52e2-407b-b3cb-62775d04c2b8")
*/
package httpbasic
