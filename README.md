# Evented Framework

## What is this thing?

Evented is a framework for implementing CQRS/ES in an enterprise setting. All the communications and data storage
elements that are specific to the architecture have been extracted and implemented once, with high reliability, moderate
to high performance, small container images, and cross-language bindings.

A complete architectural diagram (start below, under Architecture to be stepped/walked through the diagram.)

![Full Architectural Diagram](https://github.com/benjaminabbitt/evented/blob/master/Evented.svg)

Evented aims to be compliant with the [Reactive Manifesto](https://www.reactivemanifesto.org/)

## Languages

* Go (this repository)
* [Java (JVM)](https://github.com/benjaminabbitt/evented-url)

Other languages with GRPC bindings may be used trivially.

Python, .NET languages, other JVM languages (Kotlin), are all perfectly

## Benefits

### Elasticity (Scalability)

Evented fundamentally structures/architects your application to be scalable. Whether you're handling a few events a day
or a few events per millisecond, the framework is built to support your use case.

Note: this may be overkill if you're handling a few events per day, but if you want to go for it, nothing is going to
break just because the load is low.

### Big Data

By using the Evented system, you can project, re-project, re-re-project the data to different views and data warehouses
on an as-needed basis easily.

### Reactive

### Full Fidelity

The framework captures, in the event log, everything that happens that mutates the data model. Every change, creation,
or deletion is persisted for all time within the event store. This makes audit trails straightforward, as the event
log *is* the audit trail.

## Do I have to?

### Use Go?

No, GRPC allows the clients (you, a user of the framework) to write in a myriad of languages and runtimes.

### Use Domain Driven Design?

No, but it will certainly help structure your application. Sometimes, domain driven design vocabulary is used within the
Evented system, but I endeavour to keep its use as minimally as possible.

### Use Event Sourcing?

No, but it will help, and the framework makes it as easy as possible. The framework will require you to use events, but
you do not have to use event sourcing. The business logic can be passed a serialized snapshot of the business state and
a single-event list for every execution and return a new snapshot which the framework will persist.

## I have questions. Why did you...?

### Use Go?

A few factors:

* I wanted something that had robust support for GRPC, so that framework users can program in the language and runtime
  of their choosing.
* Go compiles to static executable code, reducing container sizes. Other options were unattractive (except for Rust,
  which looks great, but doesn't have first-party support for GRPC at the time of this writing).
* I wanted to learn Go, as it is gaining steam in the DevOps community.

### Use the generated protobuf/grpc data models throughout all layers of the application? Isn't that bad architecture?

Arguably yes, its bad architecture, in the conventional sense of technical architecture. At one point in the history,
the project separated the layers with different data models per layer. Ultimately, this introduced large changes as the
data models changed, and led to a few annoying-to-track-down defects in local testing.

In the end, I decided that using the data models is not bad architecture.

* The layers would be possible to re-introduce at a later time without substantial re-work (copy the generated, use an
  automated mapper). Therefor, this selection is fundamentally not architecture.
* It was more difficult (and annoying) to pass the events/commands/projections through as bytestreams. This framework
  does not care what the underlying events/commands/projections are, and has no mechanism for deserializing them. To
  decouple the data layers from protobuf/grpc would require writing an abstraction around bytestreams that I felt wasn't
  a good use of time. Keeping the data layers as protobuf `Any`s with custom wrappers around everything else doesn't
  actually achieve decoupling while substantially increasing codebase size.

## Setup

### Windows

Install:

* Chocolatey

#### Kubernetes

Enable k8s on your Windows Docker installation or install it using other tools. There are many available.

#### Consul

Consul (on host machine. Do not set it up as a server, but we will need the consul binary to load parameters in the k/v
store)

`#choco install consul`

#### Helm

`#choco install kubernetes-helm`

#### Make

`#choco install make`

### All OS

#### Install Python 3.

For Windows, use the python.org installation scripts. Chocolatey installs to weird locations and doesn't always set
the `PATH` appropriately.

#### Install cluster services

Install via helm:

* `devops/helm/consul`
* `devops/helm/mongodb`
* `devops/helm/rabbitmq`

#### Load consul k/v

`make load_all`

## Execution

### Ports

Ports are the default ports for these components

Exposed ports are whether they are exposed by default.

| Name            | Port (TCP) | Exposed            |
|-----------------|------------|--------------------|
| Command Handler | 1313       | :heavy-check-mark: |
| Query Handler   | 1314       | :heavy-check-mark: |
| Projector       | 1315       | :x:                |
| Business Logic  | 1737       | :x:                |
| Projector Logic | 1738       | :x:                |

### Building

Building requires some support python files. Dependencies are managed through virtualenv, pip package requirements are
in requirements.txt