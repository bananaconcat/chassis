# Chassis

Simple Unity Networking Solution using WebSockets

# TODO ⌛

NetworkTransform
Relay Messages? function?
.RegisterCallback(key string, method ?)
.Invoke/Call(key string)
Player Prefab
Leaving
Save
Load

Params in Unity API

`spawn <lobbyId> <prefabName> <playerId>` - Spawn a new NetworkObject Prefab with Owner for everyone

`destroy <lobbyId> <objectId>` - Destroy a NetworkObject's GameObject for everyone

`swapscene <lobbyId> <sceneName>` - Change Scene for everyone

# API 📜

`host <maxPlayers> <playerId>` - Host Lobby - join_ret lobbyId playerId

`join <lobbyId> <playerId>` Join Lobby - join_ret lobbyId playerId

`globbs` - List of Lobbies - globbs_ret lobbies

`invoke <lobbyId> <callbackKey> <params>` - Invoke Callback for every Player

`updnetvar <lobbyId> <objectid> <key> <value>` - Update NetworkVariable for NetworkObject

`leave <lobbyId> <playerId>` - Change Scene for everyone
