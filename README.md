## Adapters for Configuration and Logging in Go.

Application configuration is addressed by quite a lot of go libraries out there.
We do not intend to re-invent the wheel here, but rather we need a layer on top of existing libraries.
In particular, we'll integrate logging/tracing-configuration, making it easy to re-configure between
development and production use.

### Name

"Schuko" is the German name for a system of secure power plugs in Europe, short
for "Schutzkontakt".

<img src="http://npillmayer.github.io/img/Schuko-Stecker.svg"
    style="max-width:300px">