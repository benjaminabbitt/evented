# Evented Framework
# Not Dead
This is a work in progress.  Please star, and come back in the future.  Depending on personal and professional workload, tentative release at the end of 2020.

Note:  My family is moving, so work on this is slow.  I am continuing to work slowly, as time allows.  I hope for a more concentrated push in mid 2021.

## What is this thing?

Evented is a framework for implementing CQRS/ES in an enterprise setting.  All the communications and data storage elements that are specific to the architecture have been extracted and implemented once, with high reliability, moderate to high performance, small container images, and cross-language bindings.

A complete architectural diagram (start below, under Architecture to be stepped/walked through the diagram.)

![Full Architectural Diagram](https://github.com/benjaminabbitt/evented/blob/master/Evented.svg)

Evented aims to be compliant with the [Reactive Manifesto](https://www.reactivemanifesto.org/)

## Languages
### First class languages
First class languages have two tiers of abstraction and levels at which the developer can implement the business, projection, and saga logic.  Tier A is the same as the second class languages.  Tier B is a managed, language-specific tier that abstracts away common things like logging, configuration, some GRPC work, and other cross-cutting concerns.

* Go
* [Java (JVM)](https://github.com/benjaminabbitt/evented-url)

Future:

* C# (CLR)
* Typescript (Node)
* Python

### First and a half
First and a half languages are languages that the first class libraries should work with, but will not be tested.  Some efforts may be made to make using these languages easier with the first class language support libraries.

* Kotlin (JVM)
* Javascript (Node)
* F# (CLR)
* Scala (JVM) -- Much synergy with this framework and Akka

### Second class languages
Second class languages will work by implementing GRPC endpoints based on the provided protocol buffer definitions.  
* C++
* Dart
* Objective C
* PHP
* Ruby

## Benefits
### Elasticity (Scalability)
Evented fundamentally structures/architects your application to be scalable.  Whether you're handling a few events a day or a few events per millisecond, the framework is built to support your use case.

Note: this may be overkill if you're handling a few events per day, but if you want to go for it, nothing is going to break just because the load is low.

### Big Data
By using the Evented system, you can project, re-project, re-re-project the data to different views and data warehouses on an as-needed basis easily.

### Reactive

### Full Fidelity
The framework captures, in the event log, everything that happens that mutates the data model.  Every change, creation, or deletion is persisted for all time within the event store.  This makes audit trails straightforward, as the event log *is* the audit trail.

## What Is This?

## Do I have to?
### Use Go?
No, GRPC allows the clients (you, a user of the framework) to write in a myriad of languages and runtimes.

### Use Domain Driven Design?
No, but it will certainly help structure your application.  Sometimes, domain driven design vocabulary is used within the Evented system, but I endeavour to keep its use as minimal as possible.

### Use Event Sourcing?
No, but it will help, and the framework makes it as easy as possible.  The framework will require you to use events, but you do not have to use event sourcing.  The business logic can be passed a serialized snapshot and a single-event list for every execution and return a new snapshot which the framework will store.

## I have questions.  Why did you...?
### Use Go?
A few factors:
* I wanted something that had robust support for GRPC, so that framework users can program in the language and runtime of their choosing.
* Go can and does compile to static executable code, reducing container sizes.  Other options were unattractive (except for Rust, which looks great, but doesn't have first-class support for GRPC)
* I wanted to learn Go, as it is gaining steam in the DevOps community.

### Use the generated protobuf/grpc data models throughout all layers of the application?  Isn't that bad architecture?
Arguably yes, its bad architecture, in the conventional sense of technical architecture.  At one point in the history, the project separated the layers with different data models per layer.  Ultimately, this introduced large changes as the data models changed, and led to a few annoying-to-track-down defects in local testing.

In the end, I decided that using the data models is not bad architecture.
* The layers would be possible to re-introduce at a later time without substantial re-work (copy the generated, use an automated mapper).  Therefor, this selection is fundamentally not architecture.
* It was more difficult (and annoying) to pass the events/commands/projections through as bytestreams.  This framework does not care what the underlying events/commands/projections are, and has no mechanism for deserializing them.  To decouple the data layers from protobuf/grpc would require writing an abstraction around bytestreams that I felt wasn't a good use of my time.  Keeping the data layers as protobuf `Any`s with custom wrappers around everything else doesn't actually achieve decoupling while substantially increasing codebase size. 
