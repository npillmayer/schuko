## Warning

Under **heavy re-construction**.  Use Version v0.1.0 tag for stable releases.

# Schuko – Adapters for Configuration and Logging in Go 

Application configuration and loggin is addressed by quite a lot of Go libraries
out there. We do not intend to re-invent the wheel, but rather place a layer on
top of existing libraries.

## The Problems it Solves

A [blog post by Dave Cheney](https://dave.cheney.net/2017/01/23/the-package-level-logger-anti-pattern)
is an excellent problem statement. It says:

> The first problem with declaring a package level log variable is the tight coupling […]. Package
> `foo` now depends directly on package `mylogger` at *compile time*.
>
> The second problem is the tight coupling between package `foo` and package `mylogger` is
> transitive. Any package that consumes `package` foo is itself dependant on `mylogger` at
> compile time.
>
> This leads to a third problem, Go projects composed of packages using multiple logging libraries, *or* fiefdoms of projects who can only consume packages that use their particular logging
> library.

The “fiefdoms of projects” may not exist in reality, but the general notion is valid.
Dave's post continues with possible solutions, headlining “**Interfaces to the rescue**.”
However, the pattern he proposes may feel cumbersome, motivating this project to solve the
de-coupling challenge in a client-friendly way.

All-in-all we want the libs to be able to express their necessity for a logging
facility, but we want the main application to be in charge of how logging should
be done. **Examples** for how to do this can be found in Schuko's
package documentation.

#### Perspective of Main App

An application pulls in a lot of supporting libraries, all of which may use different
strategies for logging. The main application wants to avoid compiling in a lot of code for additional logging frameworks and it wants the supporting libs not to
pollute its own log with outputs of different formats. And sometimes it wants the
supporting libs to, frankly, shut up!

#### Perspective of Library Code

The supporting libs want to do logging, as any sensible code base will. But
what logging-framework should it use? A lib has no control over the context it
will be integrated in: Will structured logging be used or simple lines of text?
Will there be a demand for super-high-speed logging or is the focus more on 
debugging user interaction for a desktop app?

## Logging and Configuration

Configuration and logging are coupled in most applications: The config may
set up logging differently, depending on config values. And following along
on what the configuration process is doing requires logging.
We introduce slim interfaces for both, avoiding any tight coupling to
configuration frameworks or logging frameworks, and between any of these.
In particular we'll make it easy to re-configure between development and
production configuration+logging.

#### Logging from Tests

You've introduced a great new feature in your lib and now you're going to run tests
for it. But all your smartly placed log statements need a global configuration 
different from a production environment, and it would be nice to somehow synchronize
the packages' log output with `testing.T.Logf(…)`. This requires just a single
line of code with Schuko.

#### Lowest Common Denominator

When Dave Cheney
[talks about logging](https://dave.cheney.net/2015/11/05/lets-talk-about-logging),
he makes the case for a reduced set of functionality for loggers, at least in terms
of log levels. We agree to this perspective. Moreover, a facade which wants
to cover a large variety of logging frameworks needs to restrict the set of
possible operations in some way. If your project is in need of a unique feature
of a certain logging framework, Schuko may not be for you. Also, if high-end
performance is of the essence, Schuko's layer of abstraction may be too costly.

## Name

We use the term *tracing* as opposed to logging for no particular semantic
reason. It's easier this way to have package names not to be confused with
the tons of packages out there with “log” in their names.

“Schuko” is the German name for a system of secure power-plugs in continental
Europe, short for “Schutzkontakt”, which roughly translates to “safe contact”.

<img src="http://npillmayer.github.io/external/Schuko-Stecker.svg" style="max-width:230px" width="230px">
