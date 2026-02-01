# Slack-Monday Bot

This is a project loosely modelled around the concept of ChatOps that i learned from the Go For DevOps book by John Doak. 
Basically a ChatOps system is composed of 2 different services:

1. **The OPS service**, responsible for connecting to a given knowledge source an retrieving information from it. A proxy for outer systems to have insights into the said knowledge source. In the book, this was a jaeger client that offered 3 basic functionalities to outward systems: showing traces, finding a trace with spans by id, and finding logs. 

2. **The Chat serivce**, acting as an integration between a communication application and the OPS service. Basically the proxy between the users and the knowledge source. In the book, this was a slack bot.


Building on this structure, this project is a bridging between a Chat Service via Slack, just as the book, and a Monday Service. I want the users to mention this slack bot in their channels with appropriate commands that will trigger the monday.com functionalities.

1. Adding a contact to a specific contact board:
>  @contact add to [BOARD_NAME] [NAME] [SURNAME] [EMAIL] [PHONE]

This command should create a new board Item under the specified BOARD_NAME with the specified details

2. Searching a contact
> @contact find [NAME] [SURNAME]
> contact: [NAME] [SURNAME] found! 


## Project structure
```
slack-bot
    bot
        --
        --
    ops
        main.go //entrypoint
        internal/
            server/
                server.go //server that exposes API
            monday/
                client.go //monday.com client
        proto/
            ops.proto //protobuf description of server

```