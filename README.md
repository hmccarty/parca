# PARCA

Replacement for the former ARC assistant.

## TODO

- Refactor responses to be deterministic using types (e.g. don't allow client to make assumptions about interaction response vs. unprompted msg)
  - This is better with ack vs msg, but can be improved for message edit
- Fix rules on bounty and reward
- Cache reminders in case of poweroff
- Refactor error handling and traceback (log and msg all errors)
- Post errors for guilds that add application but not bot
- Add abstraction for react data

## Commands

Currency:

- balance (X)
- setbalance (X)
- leaderboard (X)
- thanks (X)
- pay (X)

Calendar:

- addcalendar (X)
- removecalendar (X)
- printcalendars (X)
- today (X)
- week (X)
- add_event ( )

General:

- status (X)
- poll (X)
- 8ball ( )
- create_role_menu (X)
- remind ( )

Domain verification:

- configuredomain (X)
- verify (X)

Games:

- arcdle ( )
- daily ( )
- bounty (X)

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
