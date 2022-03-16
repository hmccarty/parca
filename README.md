# PARCA

Replacement for the former ARC assistant.

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
			[
				- channel: [channelid]
				- calendar: [calendarid]
			]

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
