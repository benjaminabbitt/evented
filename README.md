# Evented Framework
## What Is This?

## Do I have to?
### Use Go?
No, GRPC allows the clients (you, a user of the framework) to write in a myriad of languages and runtimes.

### Use Domain Driven Design?
No, but it will certainly help.

### Use Event Sourcing?
No, but it will help, and the framework is designed to make it as easy as possible.

## Why did you...?
### Use Go?
A few factors:
* I wanted to learn Go, as it is gaining steam in the DevOps community.
* I wanted something that had robust support for GRPC, so that framework users can program in the language and runtime of their choosing.
* Go can and does compile to static executable code, reducing container sizes.

### Use the generated protobuf/grpc data models throughout all layers of the application?  Isn't that bad architecture?
Arguably yes, it's bad architecture, in the conventional sense of technical architecture.  At one point in the history, the project separated the layers with different data models per layer.  Ultimately, this introduced large changes as the data models changed, and led to a few annoying-to-track-down defects in local testing.

In the end, I decided that using the data models is not bad architecture.
* The layers would be possible to re-introduce at a later time without substantial re-work (copy the generated, use an automated mapper).  Therefor, this selection is fundamentally not architecture.
* It was more difficult (and annoying) to pass the events/commands/projections through as bytestreams.  This framework does not care what the underlying events/commands/projections are, and has no mechanism for deserializing them.  To decouple the data layers from protobuf/grpc would require writing an abstraction around bytestreams that I felt wasn't a good use of my time.  Keeping the data layers as protobuf `Any`s with custom wrappers around everything else doesn't actually achieve decoupling while substantially increasing codebase size. 