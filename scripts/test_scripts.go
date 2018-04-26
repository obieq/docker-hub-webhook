package scripts

// ScriptTest => stringified JSON script for testing
var ScriptTest = `
	{
		"scripts": [
			{
				"command": "echo",
				"args": ["hi world"]
			},
			{
				"command": "touch",
				"args": ["/does_not_exist/tree.txt"]
			}
		]
    }
`

var ScriptTestDockerServiceLS = `ID                  NAME                    MODE                REPLICAS            IMAGE                              PORTS
	usp94evkbq8o        CI_cache                replicated          1/1                 redis:3.2.6-alpine                 *:30110->6379/tcp
	hgjl4fbhfplw        CI_haproxy              replicated          1/1                 dockercloud/haproxy:latest         *:81->80/tcp, *:442->442/tcp, *:1937->1936/tcp
	5qn01tbi7uke        CI_jsreport             replicated          2/2                 jsreport/jsreport:1.3.3            *:30119->5488/tcp
	rniflno9a4hy        QA_cache                replicated          1/1                 redis:3.2.6-alpine                 *:30072->6379/tcp
	s1hdschpvigj        QA_haproxy              replicated          1/1                 dockercloud/haproxy:latest         *:82->80/tcp, *:445->445/tcp, *:1938->1936/tcp
	pp9jlle452pn        QA_jsreport             replicated          2/2                 jsreport/jsreport:1.3.3            *:30069->5488/tcp
`
