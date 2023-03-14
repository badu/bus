# Bus

In it's most simple form, this event bus can be used as following.

The listener declares it's interest for an event, by registering a handler:

`bus.Sub(OnMyEventOccurred)`

and the handler is having the following signature:

`func OnMyEventOccurred(event InterestingEvent)`

where `InterestingEvent` has to be a struct which has to implement two methods:

`type InterestingEvent struct{}`

`func (e InterestingEvent) EventID() string{ return "MyUniqueName" }`

where the string which represents the name of the event has to be unique across the event system.

`func (e InterestingEvent) Async() bool{ return true }`

where we signal that the event will be passed to the listeners by spinning up a goroutine.

The event producer will simply do:

`bus.Pub(InterestingEvent{})`

## What Problem Does It Solve?

Decoupling of components: publishers and subscribers can operate independently of each other, with no direct knowledge
of each other's existence. This decoupling allows for greater flexibility and scalability, as new publishers and
subscribers can be added to the system without disrupting existing components. Also, this facilitates testing by
triggering or ignoring certain events in some scenarios.

Asynchronous messaging: messages can be sent and received asynchronously (by spinning up goroutines), which means that
publishers and subscribers don't have to wait for each other to consume their messages. This can improve performance and
response times in a system.

Reliability: the message broker acts as a buffer between publishers and subscribers, ensuring that messages are
delivered even if one or more components in the system are temporarily unavailable.

Modularity: the Pub-Sub pattern can be used to break a monolithic application into smaller, more modular components.
Each component can then be developed and tested independently, making the overall system easier to maintain and update.

## Scenarios of Usage

Inside the `test_scenarios` folder, you can find the following scenarios:

1. Fire and Forget

2. Request Reply with Callback

3. Request Reply with Cancellation

4. Factory Request Reply

## Recommendations

1. when using `sync.WaitGroup` inside your event struct, always use method receivers and pass the event as pointer,
   otherwise you will be passing a lock by value (which is `sync.Locker`).
2. be careful if you don't want to use pointers for events, but you still need to pass values from the listener to the
   dispatcher. You should still have at least one property of that event that is a pointer (see events
   in `request reply with cancellation` for example). Same technique can be applied when you need `sync.Waitgroup` to be
   passed around with an event that is being sent by value, not by pointer.

## F.A.Q.

TBD
