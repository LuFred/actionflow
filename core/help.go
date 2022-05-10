package core

var (
	usageline = `Usage:
  main [flages]
  Start an Sever.
  main --conf-file ./config/development.yml
  `
	flagsline = `
Member:
  --env 'development'
    Service runtime environment
  --conf-file '/config/development.yml'
    Path to the config file
  --service-name 'demo'
    The service name
  --listen '0.0.0.0'
    Service binding address
  --port '4000'
    The port number for the service binding
`
)
