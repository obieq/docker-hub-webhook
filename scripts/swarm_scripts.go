package scripts

// ScriptListSwarmServices => stringified JSON script for listing all services running on a swarm cluster
// var ScriptListSwarmServices = `
// 	{
// 		"scripts": [
// 			{
// 				"command": "docker",
// 				"args": ["container", "ls"]
// 			}
// 		]
//     }
// `

func ScriptListSwarmServices() RunScript {
	return RunScript{
		Scripts: []script{
			script{Command: "docker",
				Args: []string{
					"service",
					// "container",
					"ls",
				},
			},
		},
	}
}

func ScriptUpdateSwarmService(serviceName string, image string) RunScript {
	// docker service update "$service" $detach_option $registry_auth --image="$image"
	return RunScript{
		Scripts: []script{
			script{Command: "docker",
				Args: []string{
					"service",
					"update",
					serviceName,
					"--with-registry-auth",
					`--image=` + image,
				},
			},
		},
	}
}
