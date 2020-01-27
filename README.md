./watch.sh


# TODO: pin rankings
* need randomization
* need top rank results


# TODO
* Store
* Messaging
* Payments
* Validate all user-submitted database writes
* duplicate all writes to a 2nd db instead of moving old data to old_ tables
* create expiration tables for things like events and once an hour wipe expired events

# TODO muid
* require paid accounts for events that last more than 2 or 3 time periods, charge $5 - $10 a month or something
  * this will help prevent spam


# TODO WARNINGS

# Testing
* make 1/10 requests fail and use the app - see how it goes
* make 1/10 requests timeout and use the app - see how it goes


# Training transaction process
* Trainer sets training info:
  * timeslots
  * cancellation policy
* Player selects "request timeslot"
  * ask if player has gear or nees to borrow gear
  * ask for payment information
  * validate payment information
  * notify 
    * that transaction will be billed after the session
    * cancellation policy
  * Player confirms timeslot and price of timeslot
* trainer accepts session
* send reminders
  * 2 hours before..?
* after session:
  * for player:
    * process payment
    * ask player to rate trainer
    * ask player if they want to tip
  * for trainer:
    * payment received
* after tipping:
  * notify trainer tipped
* recalculate popularity?
* TODO: refunds
* TODO: add timeout during paymennt process

## Asking for payment information
* Basics: card creation is between customer and Stripe; we only store customer ID and card id with a date for payments to be processed in the future
* API: Check if user has customer accounts on stripe, etc
  * API: for each account, if one is not found, create them so that all account creation happens in one place
  * Client: show default payment option with a "choose other method" option
    * Client: or load payment options and highlight default
* Client: Show screen to select/add/remove credit card
  * API: perform actions against stripe API. stripe user ID is in cassandra
* Client: User confirms payment method and amount
* API: send user payment details to stripe, with username
* Client: after session, user can add a tip (ask to use same credit card)
* Problem: a user's CC could be cancelled between booking and having the session. Maybe stripe can test a card the day before a session. I guess our company will have to pay the trainer, then freeze the user's account until payments are made
  * Stripe partially supports card changes: https://stripe.com/docs/saving-cards


# Convo
* remove participant
* make admin
* report
* block?
* delete message


# Thinkery
* We could base the code - object creation - around user requests
  * These user req objects could generate other objects which are returned
  * This would make much of the code flow more logical
  * It would give request objects clear logic flow
  * Objects could be compared directly, might also improve logic ruggedness
    * Might require an extra file per route

## Schedule / Planner
### TODO
* minimum lesson length - hide free blocks of time shorter than this
* on available update should we warn on appts that fall outside of new sched?
* should we allow trainers to mark sessions as group vs inidividual?
* should we limit how many active requests a user can have active at once?
* block list
* use schedule tweaks to replace availability on the weekly schedule

### Design Patterns
* Allow scheduling conflicts on our end; let the user sort all of that out. Perhaps they want some overlap or some conflicts.

## Lingo
* cql: as function prefix, sends direct CQL calls
* save: save from database
* load: load from database
* delete: delete from database
* pl: payload
* appt: appointment
* owner: owner of a planner/calendar
* mapsac: google maps places autocomplete
* Command structs are used for clarification to the app. If a success has two possible outcomes on the user's app, then Command is used to clarify.

### TODO
* implement new tables for reverse lookups instead of secondary indices (might not matter if we can run on a single node)
* Batch process some commands

## Geohash
### TODO Optimizations
* use integer constants instead of letters in switch statements to optimize

## Other TODO
* advanced error handling perhaps with sentry
* temporary user table until status is updated?
* restrict trainers being able to change name
* allow players to hide lookup by nearby


# JWT changes
* store which devices user is logged in on
  * a user can logout from one device
  * a password change triggers logout on all devices (a simple foreign key relation table in cassandra)
* solution 1) jwt's can be set to have a short expiration time (5 mins)
  * after expiration, if the user is still logged in on the login table, a new jwt is generated
    * this seems like it could be bug prone if the user fails to save the new jwt. Well, I guess the old jwt would simply be sent again and a new one would be issued, so no error would occur
	* I think the react native code has a universal request hook which can intercept and save these new jwt tokens
* solution 2) jwt's can be signed with data in the json payload or another header
  * i don't think this actually has any useful effect
* solution 3) jwt's can be cleared with push notifications
  * useful for password changes


## Cassandra
* consistency level for quorum set to 1 for dev
* probably 3 servers
  * 1) byte ordered
  * 2) normal
  * 3) questionable: high-consistency dataset
