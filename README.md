# Why is this fork a thing?

- Starting the server in 1.18 works (Mojang added a random message in 1.18 that basically just say "watch out, servers starting now" which the original repo can't yet handle.)
- Auto-Stop Server (Allows the wrapper to automatically stop the server, if the server is empty for a given duration)
- Whitelist Stuff (print current whitelist, add to- and remove from whitelist)
- Playerlist ( List() ) will now not only return a list of currently connected player, but also their position (updated twice a minute)
- I added lots of missing death-messages (haven't gotten around to testing on this though)

Thanks to wlwanpan for the heavy lifting, which allowed me to add the stuff I personally needed within a couple of hours. (original readme starts here)	
	
# minecraft-wrapper

<p align="center">
  <img src="https://github.com/ByteSizedMarius/go-minecraft-wrapper/blob/master/assets/minecraft-gopher.png?raw=true" alt="Minecraft Gopher"/>
</p>

## What is minecraft-wrapper?

Wrapper is a Go package that wraps a Minecraft Server (JE) and interacts with it by pushing in commands and reading the server logs. This package is meant to provide nicer APIs for your Go program to manage your minecraft server.

## Installation

```bash
go get github.com/ByteSizedMarius/go-minecraft-wrapper
```

## Usage

- Starting the server and listening to game events:
```go
wpr := wrapper.NewDefaultWrapper("server.jar", 1024, 1024)
wpr.Start()
defer wpr.Stop()

// Listening to game events...
for {
  select {
  case ev := <-wpr.GameEvents():
    log.Println(ev.String())
  }
}
```

- Broadcast a `"Hello"` message once the server is loaded:
```go
wpr := wrapper.NewDefaultWrapper("server.jar", 1024, 1024)
wpr.Start()
defer wpr.Stop()

<-wpr.Loaded()
wpr.Say("Hello")
```

- Retrieving a player position from the [`/data get`](https://minecraft.gamepedia.com/Commands/data#get) command:
```go
out, err := wpr.DataGet("entity", PLAYER_NAME|PLAYER_UUID)
if err != nil {
	...
}
fmt.Println(out.Pos) // [POS_X, POS_Y, POS_Z]
```

- Save the game and `Tell` a game admin `"admin-player"`, when the server is overloading.
```go
wpr := wrapper.NewDefaultWrapper("server.jar", 1024, 1024)
wpr.Start()
defer wpr.Stop()
<-wpr.Loaded()

for {
  select {
  case ev := <-wpr.GameEvents():
    if ev.String() == events.ServerOverloaded {
      if err := wpr.SaveAll(true); err != nil {
        ...
      }
      broadcastMsg := fmt.Sprintf("Server is overloaded and lagging by %sms", ev.Data["lag_time"])
      err := wpr.Tell("admin-player", broadcastMsg)
      ...
    }
  }
}
```

For more example, go to the examples dir from this repo (more will be added soon).

Note: This package is developed and tested on Minecraft 1.16, though most functionalities (`Start`, `Stop`, `Seed`, ...) works across all versions. Commands like `/data get` was introduced in version 1.13 and might not work for earlier versions. :warning: 

## Overview

<p align="center">
  <img src="https://github.com/ByteSizedMarius/go-minecraft-wrapper/blob/master/assets/architecture.png?raw=true" alt="Minecraft Wrapper Overview"/>
</p>

If you are interested, you can visit this [Medium article](https://levelup.gitconnected.com/lets-build-a-minecraft-server-wrapper-in-go-122c087e0023) to learn some of the basic inner working of the wrapper.

## Commands :construction:

The following apis/commands are from the official minecraft gamepedia [list of commands](https://minecraft.gamepedia.com/Commands#List_and_summary_of_commands) unless otherwise specified.

- [ ] [Attributes](https://minecraft.gamepedia.com/Commands/attribute)
- [ ] [Advancement](https://minecraft.gamepedia.com/Commands/advancement)
- [x] [Ban](https://minecraft.gamepedia.com/Commands/ban)
- [x] [BanIp](https://minecraft.gamepedia.com/Commands/ban#ban-ip)
- [x] [BanList](https://minecraft.gamepedia.com/Commands/ban#banlist)
- [ ] [Bossbar](https://minecraft.gamepedia.com/Commands/bossbar)
- [x] [DataGet](https://minecraft.gamepedia.com/Commands/data#get)
- [ ] [DataMerge](https://minecraft.gamepedia.com/Commands/data#merge)
- [ ] [DataModify](https://minecraft.gamepedia.com/Commands/data#modify)
- [ ] [DataRemove](https://minecraft.gamepedia.com/Commands/data#remove)
- [x] [DefaultGameMode](https://minecraft.gamepedia.com/Commands/defaultgamemode)
- [x] [DeOp](https://minecraft.gamepedia.com/Commands/deop)
- [x] [Difficulty](https://minecraft.gamepedia.com/Commands/difficulty)
- [ ] [Effect](https://minecraft.gamepedia.com/Commands/effect)
- [ ] [Enchant](https://minecraft.gamepedia.com/Commands/enchant)
- [x] [ExperienceAdd](https://godoc.org/github.com/wlwanpan/minecraft-wrapper#Wrapper.ExperienceAdd)
- [x] [ExperienceQuery](https://godoc.org/github.com/wlwanpan/minecraft-wrapper#Wrapper.ExperienceQuery)
- [ ] [Fill](https://minecraft.gamepedia.com/Commands/fill)
- [ ] [ForceLoad](https://minecraft.gamepedia.com/Commands/forceload)
- [x] [ForceLoadRemoveAll](https://minecraft.gamepedia.com/Commands/forceload)
- [ ] [Function](https://minecraft.gamepedia.com/Commands/function)
- [x] [GameEvents](https://pkg.go.dev/github.com/wlwanpan/minecraft-wrapper#Wrapper.GameEvents) - A GameEvent channel of events happening during in-game (Unofficial)
- [ ] [GameMode](https://minecraft.gamepedia.com/Commands/gamemode)
- [ ] [GameRule](https://minecraft.gamepedia.com/Commands/gamerule)
- [x] [Give](https://minecraft.gamepedia.com/Commands/give)
- [x] [Kick](https://minecraft.gamepedia.com/Commands/kick)
- [x] [Kill](https://godoc.org/github.com/wlwanpan/minecraft-wrapper#Wrapper.Kill) - Terminates the Java Process (Unofficial)
- [x] [List](https://godoc.org/github.com/wlwanpan/minecraft-wrapper#Wrapper.List) - Returns an arr of connected player struct
- [x] [Loaded](https://godoc.org/github.com/wlwanpan/minecraft-wrapper#Wrapper.Loaded) - Returns bool from a read-only channel once the server is loaded (Unofficial)
- [x] [Reload](https://minecraft.gamepedia.com/Commands/reload)
- [x] [SaveAll](https://minecraft.gamepedia.com/Commands/save#save-all)
- [x] [SaveOff](https://minecraft.gamepedia.com/Commands/save#save-off)
- [x] [SaveOn](https://minecraft.gamepedia.com/Commands/save#save-on)
- [x] [Say](https://minecraft.gamepedia.com/Commands/say)
- [ ] [Schedule](https://minecraft.gamepedia.com/Commands/scoreboard)
- [ ] [Scoreboard](https://minecraft.gamepedia.com/Commands/scoreboard)
- [x] [Seed](https://minecraft.gamepedia.com/Commands/seed)
- [ ] [SetBlock](https://minecraft.gamepedia.com/Commands/setblock)
- [x] [SetIdleTime](https://minecraft.gamepedia.com/Commands/setidletimeout)
- [ ] [SetWorldSpawn](https://minecraft.gamepedia.com/Commands/setworldspawn)
- [ ] [SpawnPoint](https://minecraft.gamepedia.com/Commands/spawnpoint)
- [ ] [Spectate](https://minecraft.gamepedia.com/Commands/spectate)
- [ ] [SpreadPlayers](https://minecraft.gamepedia.com/Commands/spreadplayers)
- [x] [Start](https://godoc.org/github.com/wlwanpan/minecraft-wrapper#Wrapper.Start) (Unofficial)
- [x] [State](https://godoc.org/github.com/wlwanpan/minecraft-wrapper#Wrapper.State) - Returns the current state of the Wrapper (Unofficial)
- [x] [Stop](https://minecraft.gamepedia.com/Commands/stop)
- [ ] [StopSound](https://minecraft.gamepedia.com/Commands/stopsound)
- [ ] [Summon](https://minecraft.gamepedia.com/Commands/summon)
- [ ] [Tag](https://minecraft.gamepedia.com/Commands/tag)
- [ ] [Team](https://minecraft.gamepedia.com/Commands/team)
- [ ] [TeamMsg](https://minecraft.gamepedia.com/Commands/teammsg)
- [ ] [Teleport](https://minecraft.gamepedia.com/Commands/teleport)
- [x] [Tell](https://minecraft.gamepedia.com/Commands/tell)
- [ ] [TellRaw](https://minecraft.gamepedia.com/Commands/tellraw)
- [x] [Tick](https://godoc.org/github.com/wlwanpan/minecraft-wrapper#Wrapper.Tick) - Returns the running game tick (Unofficial)
- [ ] [Title](https://minecraft.gamepedia.com/Commands/title)
- [ ] [Trigger](https://minecraft.gamepedia.com/Commands/trigger)
- [ ] [Weather](https://minecraft.gamepedia.com/Commands/weather)
- [x] [Whitelist](https://minecraft.gamepedia.com/Commands/whitelist)
- [ ] [WorldBorder](https://minecraft.gamepedia.com/Commands/worldborder)

Note: this list might be incomplete...

## GameEvents :construction:

List of game events and their respective data...

## Minecraft resources

- [Gamepedia](https://minecraft.gamepedia.com)
- [DigMinecraft](https://www.digminecraft.com/game_commands)

## Help and contributions

Feel free to drop a PR, file an issue or proposal of changes you want to have.
