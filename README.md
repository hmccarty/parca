# PARCA

Replacement for the former ARC assistant.

## Commands

Currency:
- get_balance (X)
- set_balance (X)
- leaderboard (X)
- thanks      (X)
- pay         (X)

Calendar:
- set_calendar    ( )
- remove_calendar ( )
- print_calendars ( )
- today           ( )
- week            ( )
- add_event       ( )

General:
- status           ( )
- poll             ( )
- 8ball            ( )
- create_role_menu ( )
- remind           ( )

Domain verification:
- configuredomain ( )
- verify          ( )

Games:
- arcdle ( )
- daily  ( )
- bounty ( )

## REDIS Structure

```
user:
	<userid>:
		- username: str
		- balance: float

verify:
	<guildid>:
		- domain: str
		- role: roleid
		<userid>:
			- code: int

arcdle:
	<userid>:
		- channel: channelid
		- message: messageid
		- status: int
		- hidden: str
		- visible: str

daily: [userid]

backlog: [str]

calendar:
	<guildid>:
		<channelid>:
			- [calendarid]

bounty:
	<guildid>:
		[
			title: str
			user: userid
			guild: guildid
			channel: channelid
			message: messageid
			amt: float
		]
```
