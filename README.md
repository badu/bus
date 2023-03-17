# Bus

In it's most simple form, this event bus can be used as following.

The listener declares it's interest for an event, by registering a handler:

`bus.Sub(OnMyEventOccurred)`

and the handler is having the following signature:

`func OnMyEventOccurred(event InterestingEvent)`

The event producer will simply do:

`bus.Pub(InterestingEvent{})`

Optional, to allow the bus to spin a goroutine for dispatching events, implement the following interface:

`func (e InterestingEvent) Async() bool{ return true }`

or

`func (e *InterestingEvent) Async() bool{ return true }`

By default, the bus is using sync events : waits for listeners to complete their jobs before calling the next listener.

Usage : `go get github.com/badu/bus`

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

1. Fire and Forget.

   Imagine a system / application where we have three services : `users`, `notifications` (email and
   SMS) and `audit`. When a user registers, we want to send welcoming messages via SMS and email, but we also want to
   audit that registration for reporting purposes.

   The [UserRegisteredEvent](https://github.com/badu/bus/blob/main/test_scenarios/fire-and-forget/events/main.go#L10)
   will carry the freshly registered username (which is also the email) and phone to the email and sms services. The
   event is [triggered](https://github.com/badu/bus/blob/main/test_scenarios/fire-and-forget/users/service.go#L21) by
   the user service, which performs the creation of the user account. We're using the `fire and forget` technique here,
   because the operation of registration should not depend on the fact that we've been able to
   send an welcoming email or a sms, or the audit system malfunctions.

   Simulating audit service malfunction is easy. Instead of using `Sub`, we're using `SubUnsub` to register the listener
   and return [`true`](https://github.com/badu/bus/blob/main/test_scenarios/fire-and-forget/audit/service.go#L36) to
   unsubscribe on events of that kind.

2. Factory Request Reply

   Imagine a system / application where we need to communicate with different microservices, but in this case we don't
   want to bring them online, we're just wanting to stub the response as those services were alive.

   This technique is useful when we need to test some complicated flows of business logic and facilitates the
   transformation of an integration test into a classic unit test.

   The `cart` service requires two replies from two other microservices `inventory` and `prices`. In the past, I've been
   using a closure function to provide the service with both real GRPC clients or with mocks and stubs. The service
   signature gets complicated and large as we one service would depend on a lot of GRPC clients to aggregate data.

   As you can see
   the [test here](https://github.com/badu/bus/blob/main/test_scenarios/factory-request-reply/main_test.go) it's much
   more elegant and the service constructor is much slimmer.

   Events are one sync and one async, just to check it works in both scenarios.

   Important to note that because a `WaitGroup` is being used in our event struct, we're forced to pass the events by
   using a pointer, instead of passing them by value.

3. Request Reply with Callback

   In this example, we wanted to achieve two things. First is that the `service` and the `repository` are decoupled by
   events. More than that, we wanted that the events are generic on their own.

   The orders service will dispatch a generic request event, one for placing an order, which will carry an `Order` (
   model) struct with that request and another `OrderStatus` (model) struct using the same generic event.

   We are using a channel inside the generic `RequestEvent` to signal the `reply` to the publisher, which in this case
   is a callback function that returns the result as if the publisher would have called directly the listener.

   I am sure that you will find this technique interesting and having a large number of applications.

4. Request Reply with Cancellation

   Last but, not least, this is an example about providing `context.Context` along the publisher subscriber chain.
   The `repository` is simulating a long database call, longer than the context's cancellation, so the service gets the
   deadline exceeded error.

   Note that this final example is not using pointer to the event's struct, but it contains two properties which have
   pointers, so the `service` can access the altered `reply`.

## Recommendations

1. always place your events inside a separate `events` package, avoiding circular dependencies.
2. in general, in `request-reply` scenarios, the events should be passed as pointers (even if it's somewhat slower),
   because changing properties that represents the `reply` would not be reflected. Also, when using `sync.WaitGroup`
   inside your event struct, always use method receivers and pass the event as pointer, otherwise you will be passing a
   lock by value (which is `sync.Locker`).
3. be careful if you don't want to use pointers for events, but you still need to pass values from the listener to the
   dispatcher. You should still have at least one property of that event that is a pointer (see events
   in `request reply with cancellation` for example). Same technique can be applied when you need `sync.Waitgroup` to be
   passed around with an event that is being sent by value, not by pointer.
4. you can override the event name (which is by default, built using `fmt.Sprintf("%T", yourEvent)`) you need to
   implement `EventID() string` interface.
