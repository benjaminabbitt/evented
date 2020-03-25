#Command Handler
The command handler accepts commands and events and processes them.

##Commands
Accepts commands, in the form of a `CommandBook`, calls a business logic function via GRPC, receives back an `EventBook`, and passes that to event handling.
 
##Events
Accepts `EventBook`s, analyzes them, potentially breaks them into multiple `EventBook`s, and then forwards them along either synchronously or asynchronously, depending on what the Event Pages specify.