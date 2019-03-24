# Financial Instruments Broker/Dealer

brokerdealer models basic components of a Financial Instruments (Securities) Broker/Dealer.

The primary focus is simple modelling of different business capabilities within a typical broker/dealer organization, with a POC to illustrate how those capabilities could collaborate together efficiently to operate the business. 

The technical goals are:
- create basic models of different business capabilities (typically owned by a specific organizational silo) - DDD Bounded Contexts, if you will
- have the business capabilities communicate in a decoupled fashion through use of asynchronous messages (choreography style), with NATS.io being chosen initially because of its ease of use

A Japan-based organization is assumed, and thus Japan specific financial instruments market practices and conventions are modelled. The different business capabilities of a typical dealer organization to be modelled are:
- trade booking (at least one)
- trade matching (at least one)
- central clearing / novation
- settlement management
- position management
- financial accounting
- credit risk management

The domain events that occur are:
- trades.booking.booked
- trades.booking.cancelled
- trades.booking.amended
- trades.matching.submitted
- trades.matching.matched
- trades.centralclearing.eligible
- trades.centralclearing.ineligible
- trades.centralclearing.novated
- settlements.instructed
- settlements.settled
- settlements.failed
- position.updated

Run the trade_booking simulator:
```
go run trade_booking.go --url=nats://demo.nats.io:4222
```

