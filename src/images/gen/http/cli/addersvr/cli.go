// Code generated by goa v3.12.4, DO NOT EDIT.
//
// addersvr HTTP client CLI support package
//
// Command:
// $ goa gen github.com/dipankardas011/crop-yield-monitor/src/images/design

package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	serversc "github.com/dipankardas011/crop-yield-monitor/src/images/gen/http/servers/client"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//	command (subcommand1|subcommand2|...)
func UsageCommands() string {
	return `servers (upload|fetch|get-health)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` servers upload --body '{
      "image": "Q29uc2VxdWF0dXIgaGFydW0gZXQgZG9sb3IgcmVwZWxsZW5kdXMgZXhwbGljYWJvIGVpdXMu",
      "uuid": "1"
   }'` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(
	scheme, host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restore bool,
) (goa.Endpoint, any, error) {
	var (
		serversFlags = flag.NewFlagSet("servers", flag.ContinueOnError)

		serversUploadFlags    = flag.NewFlagSet("upload", flag.ExitOnError)
		serversUploadBodyFlag = serversUploadFlags.String("body", "REQUIRED", "")

		serversFetchFlags    = flag.NewFlagSet("fetch", flag.ExitOnError)
		serversFetchBodyFlag = serversFetchFlags.String("body", "REQUIRED", "")

		serversGetHealthFlags = flag.NewFlagSet("get-health", flag.ExitOnError)
	)
	serversFlags.Usage = serversUsage
	serversUploadFlags.Usage = serversUploadUsage
	serversFetchFlags.Usage = serversFetchUsage
	serversGetHealthFlags.Usage = serversGetHealthUsage

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	if flag.NArg() < 2 { // two non flag args are required: SERVICE and ENDPOINT (aka COMMAND)
		return nil, nil, fmt.Errorf("not enough arguments")
	}

	var (
		svcn string
		svcf *flag.FlagSet
	)
	{
		svcn = flag.Arg(0)
		switch svcn {
		case "servers":
			svcf = serversFlags
		default:
			return nil, nil, fmt.Errorf("unknown service %q", svcn)
		}
	}
	if err := svcf.Parse(flag.Args()[1:]); err != nil {
		return nil, nil, err
	}

	var (
		epn string
		epf *flag.FlagSet
	)
	{
		epn = svcf.Arg(0)
		switch svcn {
		case "servers":
			switch epn {
			case "upload":
				epf = serversUploadFlags

			case "fetch":
				epf = serversFetchFlags

			case "get-health":
				epf = serversGetHealthFlags

			}

		}
	}
	if epf == nil {
		return nil, nil, fmt.Errorf("unknown %q endpoint %q", svcn, epn)
	}

	// Parse endpoint flags if any
	if svcf.NArg() > 1 {
		if err := epf.Parse(svcf.Args()[1:]); err != nil {
			return nil, nil, err
		}
	}

	var (
		data     any
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "servers":
			c := serversc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "upload":
				endpoint = c.Upload()
				data, err = serversc.BuildUploadPayload(*serversUploadBodyFlag)
			case "fetch":
				endpoint = c.Fetch()
				data, err = serversc.BuildFetchPayload(*serversFetchBodyFlag)
			case "get-health":
				endpoint = c.GetHealth()
				data = nil
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// serversUsage displays the usage of the servers command and its subcommands.
func serversUsage() {
	fmt.Fprintf(os.Stderr, `server handlers
Usage:
    %[1]s [globalflags] servers COMMAND [flags]

COMMAND:
    upload: Upload implements upload.
    fetch: Fetch implements fetch.
    get-health: GetHealth implements get health.

Additional help:
    %[1]s servers COMMAND --help
`, os.Args[0])
}
func serversUploadUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] servers upload -body JSON

Upload implements upload.
    -body JSON: 

Example:
    %[1]s servers upload --body '{
      "image": "Q29uc2VxdWF0dXIgaGFydW0gZXQgZG9sb3IgcmVwZWxsZW5kdXMgZXhwbGljYWJvIGVpdXMu",
      "uuid": "1"
   }'
`, os.Args[0])
}

func serversFetchUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] servers fetch -body JSON

Fetch implements fetch.
    -body JSON: 

Example:
    %[1]s servers fetch --body '{
      "uuid": "1"
   }'
`, os.Args[0])
}

func serversGetHealthUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] servers get-health

GetHealth implements get health.

Example:
    %[1]s servers get-health
`, os.Args[0])
}
