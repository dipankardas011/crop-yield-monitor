// Code generated by goa v3.12.4, DO NOT EDIT.
//
// addersvr HTTP client CLI support package
//
// Command:
// $ goa gen
// github.com/dipankardas011/crop-yield-monitor/src/authentication/design

package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	serversc "github.com/dipankardas011/crop-yield-monitor/src/authentication/gen/http/servers/client"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//	command (subcommand1|subcommand2|...)
func UsageCommands() string {
	return `servers (login|signup|get-health)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` servers login --body '{
      "password": "77777",
      "username": "demo"
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

		serversLoginFlags    = flag.NewFlagSet("login", flag.ExitOnError)
		serversLoginBodyFlag = serversLoginFlags.String("body", "REQUIRED", "")

		serversSignupFlags    = flag.NewFlagSet("signup", flag.ExitOnError)
		serversSignupBodyFlag = serversSignupFlags.String("body", "REQUIRED", "")

		serversGetHealthFlags = flag.NewFlagSet("get-health", flag.ExitOnError)
	)
	serversFlags.Usage = serversUsage
	serversLoginFlags.Usage = serversLoginUsage
	serversSignupFlags.Usage = serversSignupUsage
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
			case "login":
				epf = serversLoginFlags

			case "signup":
				epf = serversSignupFlags

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
			case "login":
				endpoint = c.Login()
				data, err = serversc.BuildLoginPayload(*serversLoginBodyFlag)
			case "signup":
				endpoint = c.Signup()
				data, err = serversc.BuildSignupPayload(*serversSignupBodyFlag)
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
	fmt.Fprintf(os.Stderr, `ksctl server handlers
Usage:
    %[1]s [globalflags] servers COMMAND [flags]

COMMAND:
    login: Login implements login.
    signup: Signup implements signup.
    get-health: GetHealth implements get health.

Additional help:
    %[1]s servers COMMAND --help
`, os.Args[0])
}
func serversLoginUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] servers login -body JSON

Login implements login.
    -body JSON: 

Example:
    %[1]s servers login --body '{
      "password": "77777",
      "username": "demo"
   }'
`, os.Args[0])
}

func serversSignupUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] servers signup -body JSON

Signup implements signup.
    -body JSON: 

Example:
    %[1]s servers signup --body '{
      "emailid": "demo@xyz.com",
      "first": "hello",
      "last": "world",
      "password": "77777"
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