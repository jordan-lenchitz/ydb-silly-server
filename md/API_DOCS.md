Go
Skip to Main Content

    Why Go submenu dropdown icon
    Learn
    Docs submenu dropdown icon
    Packages
    Community submenu dropdown icon

API Documentation

Pkgsite provides a JSON API for accessing information about Go packages and modules.
Requests

All requests are GET requests whose first URL segment selects the request type. Except for search, the remaining segments specify a module or package path. Additional arguments are specified as query parameters.

Package paths are ambiguous: the same path a/b/c could be package c in module a/b or package b/c in module a. If both exist, the API returns an error with a list of possible modules in its candidates field. Make a follow-up request with the module query parameter to disambiguate. (Note: this behavior differs from the pkgsite UI, which selects the longest matching module path.)

Some routes accept an optional version query parameter. If missing, the latest version is chosen.

Some routes accept optional GOOS and GOARCH query parameters to select the desired build context of the documentation. If omitted, the default documentation build context is selected, typically linux/amd64.

Some routes accept a filter. Only items that match the filter are returned. The filter must be a Go expression that returns a boolean value. Filters support the following subset of Go:

    The values true, false and nil.
    == and != on any value.
    +, -,*, / and % on integers.
    + on strings.
    <, <=, > and >= on strings and integers.
    The following functions:

    contains(s, sub)
        String s contains substring sub.
    hasPrefix(s, pre)
        String s has the prefix pre.
    hasSuffix(s, suf)
        String s has the suffix suf.
    matches(s, re)
        String s matches the regular expression re.

In addition, each route adds its own variables for the list items being filtered. In general, each JSON field is a variable.

Most filters will need to be percent-encoded. Go programs should call url.QueryEscape to encode filters and all other query parameter values. Or better, put all values into a url.Values and then call Encode on that.
Responses

All routes return a JSON response. You can find a link to the corresponding Go type for a route in the documentation below. Errors are also returned as a JSON object, whose message field contains the error message.

You can access the machine-readable OpenAPI Specification.
Pagination

Some responses are paginated. Each response page has a nextPageToken, which is non-empty if there is another page of items. To retrieve the next page, repeat the original request verbatim, adding a token query parameter whose value is the next-page token from the response. Changing the request in any way other than providing a token may result in an error.

As long as nextPageToken is non-empty, there is another page, even if the current page has no items.

Routes

    /v1beta/package/{path}

    Information about the package at {path}.
    Parameters:
    Name 	Type 	Description
    module 	string 	Module path.
    version 	string 	Module version: semantic version, 'latest', or default branches 'master' or 'main'. (Latest if empty).
    goos 	string 	GOOS of documentation build context.
    goarch 	string 	GOARCH of documentation build context.
    doc 	string 	Documentation format: text, html, md or markdown. If omitted, documentation is not returned.
    examples 	bool 	Whether to include examples with the returned documentation.
    imports 	bool 	Whether to include the packages that this one imports.
    licenses 	bool 	Whether to include licenses in the result.
    Response:
    Package
    Example Request

    curl -L https://pkg.go.dev/v1beta/package/golang.org/x/time/rate

    Example Response

    {
      "modulePath": "golang.org/x/time",
      "version": "v0.15.0",
      "isLatest": true,
      "isStandardLibrary": false,
      "goos": "all",
      "goarch": "all",
      "path": "golang.org/x/time/rate",
      "name": "rate",
      "synopsis": "Package rate provides a rate limiter.",
      "isRedistributable": true
    }

    /v1beta/module/{path}

    Information about the module at {path}.
    Parameters:
    Name 	Type 	Description
    version 	string 	Module version: semantic version, 'latest', or default branches 'master' or 'main'. (Latest if empty).
    licenses 	bool 	Whether to include licenses in the result.
    readme 	bool 	Whether to include the README in the result.
    Response:
    Module
    Example Request

    curl -L https://pkg.go.dev/v1beta/module/golang.org/x/time

    Example Response

    {
      "path": "golang.org/x/time",
      "version": "v0.15.0",
      "commitTime": "2026-02-11T19:14:29Z",
      "isLatest": true,
      "isRedistributable": true,
      "isStandardLibrary": false,
      "hasGoMod": true,
      "repoUrl": "https://cs.opensource.google/go/x/time"
    }

    /v1beta/versions/{path}

    Versions of the module at {path}. If there are tagged versions, they are returned. Otherwise, the 10 most recent pseudo-versions are returned. The versions are in descending order. Only results that match the filter query parameter are returned.
    Parameters:
    Name 	Type 	Description
    limit 	int 	Max number of items to return.
    token 	string 	Where to resume listing.
    filter 	string 	Include only items matching the regular expression filter.
    Response:
    PaginatedResponse[ModuleVersion]
    Example Request

    curl -L https://pkg.go.dev/v1beta/versions/golang.org/x/time?limit=3

    Example Response

    {
      "items": [
        {
          "modulePath": "golang.org/x/time",
          "version": "v0.15.0",
          "commitTime": "2026-02-11T19:14:29Z",
          "isRedistributable": true,
          "hasGoMod": true,
          "latestVersion": "v0.15.0",
          "deprecated": false,
          "deprecationReason": "",
          "retracted": false,
          "retractionReason": ""
        },
        {
          "modulePath": "golang.org/x/time",
          "version": "v0.14.0",
          "commitTime": "2025-09-16T23:29:52Z",
          "isRedistributable": true,
          "hasGoMod": true,
          "latestVersion": "v0.15.0",
          "deprecated": false,
          "deprecationReason": "",
          "retracted": false,
          "retractionReason": ""
        },
        {
          "modulePath": "golang.org/x/time",
          "version": "v0.13.0",
          "commitTime": "2025-08-13T14:44:00Z",
          "isRedistributable": true,
          "hasGoMod": true,
          "latestVersion": "v0.15.0",
          "deprecated": false,
          "deprecationReason": "",
          "retracted": false,
          "retractionReason": ""
        }
      ],
      "total": 15,
      "nextPageToken": "185e216cd0076de0df84b79cf2e3210d"
    }

    /v1beta/packages/{path}

    Information about packages of the module at {path}. Filtering is applied to the list of packages in the response. Only packages that match the filter query parameter are returned.
    Parameters:
    Name 	Type 	Description
    version 	string 	Module version: semantic version, 'latest', or default branches 'master' or 'main'. (Latest if empty).
    limit 	int 	Max number of items to return.
    token 	string 	Where to resume listing.
    filter 	string 	Include only items matching the regular expression filter.
    Response:
    PackagesResponse
    Example Request

    curl -L https://pkg.go.dev/v1beta/packages/golang.org/x/time/rate

    Example Response

    {
      "code": 400,
      "message": "golang.org/x/time/rate is a package, not a module",
      "fixes": [
        "retry the call with the containing module: \"golang.org/x/time\""
      ]
    }

    /v1beta/search

    Search results. Only results that match the filter query parameter are returned. Results are sorted by how well the match the query, with the best match first.
    Parameters:
    Name 	Type 	Description
    q 	string 	Find packages matching this query.
    symbol 	string 	If non-empty, find symbols matching this string. The query further restricts the search to matching packages.
    limit 	int 	Max number of items to return.
    token 	string 	Where to resume listing.
    filter 	string 	Include only items matching the regular expression filter.
    Response:
    PaginatedResponse[SearchResult]
    Example Request

    curl -L https://pkg.go.dev/v1beta/search?q=xyzzy

    Example Response

    {
      "items": [
        {
          "packagePath": "github.com/docktermj/go-xyzzy-helpers/logger",
          "modulePath": "github.com/docktermj/go-xyzzy-helpers",
          "version": "v0.2.2",
          "synopsis": "Package helper ..."
        },
        {
          "packagePath": "github.com/docktermj/go-xyzzy-helpers/g2configuration",
          "modulePath": "github.com/docktermj/go-xyzzy-helpers",
          "version": "v0.2.2",
          "synopsis": "Package helper ..."
        },
        {
          "packagePath": "github.com/docktermj/go-xyzzy-grpc/g2diagnostic",
          "modulePath": "github.com/docktermj/go-xyzzy-grpc",
          "version": "v0.0.0-20221028160434-9c34bd3a1481",
          "synopsis": ""
        },
        {
          "packagePath": "github.com/docktermj/go-xyzzy-helpers/logmessage",
          "modulePath": "github.com/docktermj/go-xyzzy-helpers",
          "version": "v0.2.2",
          "synopsis": "Package message ..."
        },
        {
          "packagePath": "github.com/wade-rees-me/striker-go/cmd/striker/xyzzy",
          "modulePath": "github.com/wade-rees-me/striker-go",
          "version": "v1.0.3",
          "synopsis": ""
        },
        {
          "packagePath": "github.com/docktermj/go-xyzzy-grpc/g2diagnosticclient",
          "modulePath": "github.com/docktermj/go-xyzzy-grpc",
          "version": "v0.0.0-20221028160434-9c34bd3a1481",
          "synopsis": "Package main implements a client for the service."
        },
        {
          "packagePath": "github.com/docktermj/go-xyzzy-grpc/g2diagnosticclientcli",
          "modulePath": "github.com/docktermj/go-xyzzy-grpc",
          "version": "v0.0.0-20221028160434-9c34bd3a1481",
          "synopsis": "Package main implements a client for the service."
        },
        {
          "packagePath": "github.com/docktermj/go-xyzzy-grpc/g2diagnosticserver",
          "modulePath": "github.com/docktermj/go-xyzzy-grpc",
          "version": "v0.0.0-20221028160434-9c34bd3a1481",
          "synopsis": ""
        },
        {
          "packagePath": "github.com/docktermj/go-xyzzy-grpc/g2diagnosticservercli",
          "modulePath": "github.com/docktermj/go-xyzzy-grpc",
          "version": "v0.0.0-20221028160434-9c34bd3a1481",
          "synopsis": ""
        },
        {
          "packagePath": "github.com/docktermj/go-xyzzy-helpers",
          "modulePath": "github.com/docktermj/go-xyzzy-helpers",
          "version": "v0.2.2",
          "synopsis": ""
        },
        {
          "packagePath": "github.com/docktermj/go-xyzzy-grpc/v2/g2diagnostic",
          "modulePath": "github.com/docktermj/go-xyzzy-grpc/v2",
          "version": "v2.0.0-20220622205531-7c693c6f8eb1",
          "synopsis": ""
        },
        {
          "packagePath": "github.com/docktermj/go-xyzzy-grpc/v2/server",
          "modulePath": "github.com/docktermj/go-xyzzy-grpc/v2",
          "version": "v2.0.0-20220622205531-7c693c6f8eb1",
          "synopsis": ""
        },
        {
          "packagePath": "github.com/blasphemy/xyzzy2/lib/xyzzy",
          "modulePath": "github.com/blasphemy/xyzzy2",
          "version": "v0.0.0-20141223201256-7f5f89b3fd2e",
          "synopsis": ""
        }
      ],
      "total": 13
    }

    /v1beta/symbols/{path}

    List of symbols for the package at {path}. Filtering is applied to the list of symbols in the response. Only symbols that match the filter query parameter are returned.
    Parameters:
    Name 	Type 	Description
    module 	string 	Module path.
    version 	string 	Module version: semantic version, 'latest', or default branches 'master' or 'main'. (Latest if empty).
    goos 	string 	GOOS of documentation build context.
    goarch 	string 	GOARCH of documentation build context.
    limit 	int 	Max number of items to return.
    token 	string 	Where to resume listing.
    filter 	string 	Include only items matching the regular expression filter.
    Response:
    PackageSymbols
    Example Request

    curl -L https://pkg.go.dev/v1beta/symbols/golang.org/x/time/rate

    Example Response

    {
      "modulePath": "golang.org/x/time",
      "version": "v0.15.0",
      "symbols": {
        "items": [
          {
            "name": "Limit",
            "kind": "Type",
            "synopsis": "type Limit float64",
            "parent": "Limit"
          },
          {
            "name": "Every",
            "kind": "Function",
            "synopsis": "func Every(interval time.Duration) Limit",
            "parent": "Limit"
          },
          {
            "name": "Limiter",
            "kind": "Type",
            "synopsis": "type Limiter struct{}",
            "parent": "Limiter"
          },
          {
            "name": "Limiter.Allow",
            "kind": "Method",
            "synopsis": "func (lim *Limiter) Allow() bool",
            "parent": "Limiter"
          },
          {
            "name": "Limiter.AllowN",
            "kind": "Method",
            "synopsis": "func (lim *Limiter) AllowN(t time.Time, n int) bool",
            "parent": "Limiter"
          },
          {
            "name": "Limiter.Burst",
            "kind": "Method",
            "synopsis": "func (lim *Limiter) Burst() int",
            "parent": "Limiter"
          },
          {
            "name": "Limiter.Limit",
            "kind": "Method",
            "synopsis": "func (lim *Limiter) Limit() Limit",
            "parent": "Limiter"
          },
          {
            "name": "Limiter.Reserve",
            "kind": "Method",
            "synopsis": "func (lim *Limiter) Reserve() *Reservation",
            "parent": "Limiter"
          },
          {
            "name": "Limiter.ReserveN",
            "kind": "Method",
            "synopsis": "func (lim *Limiter) ReserveN(t time.Time, n int) *Reservation",
            "parent": "Limiter"
          },
          {
            "name": "Limiter.SetBurst",
            "kind": "Method",
            "synopsis": "func (lim *Limiter) SetBurst(newBurst int)",
            "parent": "Limiter"
          },
          {
            "name": "Limiter.SetBurstAt",
            "kind": "Method",
            "synopsis": "func (lim *Limiter) SetBurstAt(t time.Time, newBurst int)",
            "parent": "Limiter"
          },
          {
            "name": "Limiter.SetLimit",
            "kind": "Method",
            "synopsis": "func (lim *Limiter) SetLimit(newLimit Limit)",
            "parent": "Limiter"
          },
          {
            "name": "Limiter.SetLimitAt",
            "kind": "Method",
            "synopsis": "func (lim *Limiter) SetLimitAt(t time.Time, newLimit Limit)",
            "parent": "Limiter"
          },
          {
            "name": "Limiter.Tokens",
            "kind": "Method",
            "synopsis": "func (lim *Limiter) Tokens() float64",
            "parent": "Limiter"
          },
          {
            "name": "Limiter.TokensAt",
            "kind": "Method",
            "synopsis": "func (lim *Limiter) TokensAt(t time.Time) float64",
            "parent": "Limiter"
          },
          {
            "name": "Limiter.Wait",
            "kind": "Method",
            "synopsis": "func (lim *Limiter) Wait(ctx context.Context) (err error)",
            "parent": "Limiter"
          },
          {
            "name": "Limiter.WaitN",
            "kind": "Method",
            "synopsis": "func (lim *Limiter) WaitN(ctx context.Context, n int) (err error)",
            "parent": "Limiter"
          },
          {
            "name": "NewLimiter",
            "kind": "Function",
            "synopsis": "func NewLimiter(r Limit, b int) *Limiter",
            "parent": "Limiter"
          },
          {
            "name": "Reservation",
            "kind": "Type",
            "synopsis": "type Reservation struct{}",
            "parent": "Reservation"
          },
          {
            "name": "Reservation.Cancel",
            "kind": "Method",
            "synopsis": "func (r *Reservation) Cancel()",
            "parent": "Reservation"
          },
          {
            "name": "Reservation.CancelAt",
            "kind": "Method",
            "synopsis": "func (r *Reservation) CancelAt(t time.Time)",
            "parent": "Reservation"
          },
          {
            "name": "Reservation.Delay",
            "kind": "Method",
            "synopsis": "func (r *Reservation) Delay() time.Duration",
            "parent": "Reservation"
          },
          {
            "name": "Reservation.DelayFrom",
            "kind": "Method",
            "synopsis": "func (r *Reservation) DelayFrom(t time.Time) time.Duration",
            "parent": "Reservation"
          },
          {
            "name": "Reservation.OK",
            "kind": "Method",
            "synopsis": "func (r *Reservation) OK() bool",
            "parent": "Reservation"
          },
          {
            "name": "Sometimes",
            "kind": "Type",
            "synopsis": "type Sometimes struct{ ... }",
            "parent": "Sometimes"
          },
          {
            "name": "Sometimes.Do",
            "kind": "Method",
            "synopsis": "func (s *Sometimes) Do(f func())",
            "parent": "Sometimes"
          },
          {
            "name": "Sometimes.Every",
            "kind": "Field",
            "synopsis": "Every int",
            "parent": "Sometimes"
          },
          {
            "name": "Sometimes.First",
            "kind": "Field",
            "synopsis": "First int",
            "parent": "Sometimes"
          },
          {
            "name": "Sometimes.Interval",
            "kind": "Field",
            "synopsis": "Interval time.Duration",
            "parent": "Sometimes"
          },
          {
            "name": "Inf",
            "kind": "Constant",
            "synopsis": "const Inf",
            "parent": "Inf"
          },
          {
            "name": "InfDuration",
            "kind": "Constant",
            "synopsis": "const InfDuration",
            "parent": "InfDuration"
          }
        ],
        "total": 31
      }
    }

    /v1beta/imported-by/{path}

    Paths of packages importing the package at {path}, not including packages in the same module. Filtering is applied to the list of paths in the response. Only paths that match the filter query parameter are returned. Within a filter, the variable `path` is set to the import path.
    Parameters:
    Name 	Type 	Description
    module 	string 	Module path.
    version 	string 	Module version: semantic version, 'latest', or default branches 'master' or 'main'. (Latest if empty).
    limit 	int 	Max number of items to return.
    token 	string 	Where to resume listing.
    filter 	string 	Include only items matching the regular expression filter.
    Response:
    PackageImportedBy
    Example Request

    curl -L https://pkg.go.dev/v1beta/imported-by/golang.org/x/time/rate?limit=10&filter=%5E.%2A%5C.io%2F

    Example Response

    {
      "code": 400,
      "message": "parsing filter \"^.*\\.io/\": 1:2: expected operand, found '.' (and 1 more errors)",
      "fixes": [
        "the 'filter' query parameter must be a valid Go expression; see the documentation at /v1beta/api"
      ]
    }

    /v1beta/vulns/{path}

    Vulnerabilities of the module or package at {path}, from the Go vulnerability database (https://vuln.go.dev). Only results that match the filter query parameter are returned.
    Parameters:
    Name 	Type 	Description
    module 	string 	Module path.
    version 	string 	Module version: semantic version, 'latest', or default branches 'master' or 'main'. (Latest if empty).
    limit 	int 	Max number of items to return.
    token 	string 	Where to resume listing.
    filter 	string 	Include only items matching the regular expression filter.
    Response:
    PaginatedResponse[Vulnerability]
    Example Request

    curl -L https://pkg.go.dev/v1beta/vulns/golang.org/x/image

    Example Response

    {
      "items": null,
      "total": 0
    }

Why Go
Use Cases
Case Studies
Get Started
Playground
Tour
Stack Overflow
Help
Packages
Standard Library
Sub-repositories
About Go Packages
About
Download
Blog
Issue Tracker
Release Notes
Brand Guidelines
Code of Conduct
Connect
Twitter
GitHub
Slack
r/golang
Meetup
Golang Weekly
Gopher in flight goggles

    Copyright
    Terms of Service
    Privacy Policy
    Report an Issue

Google logo
go.dev uses cookies from Google to deliver and enhance the quality of its services and to analyze traffic. Learn more.
 :)