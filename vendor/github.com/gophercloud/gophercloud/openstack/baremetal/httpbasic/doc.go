/*
Package httpbasic provides support for http_basic bare metal endpoints.

Example of obtaining and using a client:

      client, err := httpbasic.NewBareMetalHTTPBasic(httpbasic.Endpoints{
		IronicEndpoing:     "http://localhost:6385/v1/",
		IronicUser:         "myUser",
		IronicUserPassword: "myPassword",
	})
	if err != nil {
		panic(err)
	}

	client.Microversion = "1.50"
	nodes.ListDetail(client, nodes.listOpts{})
*/
package httpbasic
