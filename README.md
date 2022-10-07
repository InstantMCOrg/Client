# InstantMinecraft Server
A standalone container running a minecraft server which is controllable through an http endpoint

## Setup

````bash
docker pull ghcr.io/instantminecraft/client:latest
docker run -d --name mcclient -p 25585:25585 ghcr.io/instantminecraft/client:latest
curl localhost:25585/server/start
````

## HTTP Endpoints
Port: 25585

### ``GET /``
Returns the server status

Response:
````json
{
    "server": {
        "running": false
    }
}
````

### ``GET /server/start?blocking=true``
Starts the Minecraft Server
- ``blocking`` is optional. If true the response is sent when the minecraft server has fully booted up

Response:
````json
{
  "message": "Minecraft Server has been started"
}
````

### ``GET /server/stop?blocking=true``
Stops the Minecraft Server
- ``blocking`` is optional. If true the response is sent when the minecraft server has fully stopped

Response:
````json
{
  "message": "Minecraft Server has stopped"
}
````

### ``GET /server/player/op/{playername}``
Grants operator permission to the target player

