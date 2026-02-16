// Package trace2go provides a logger/tracer factory in the spirit of log4j, but much simpler.
//
//    tracing.RegisterTraceAdapter("test", getTT, true)
//    tracing.SetTraceSelector(trace2go.Selector()) // install trace2go as global selector
//    conf := testconfig.Conf{                      // lightweight configuration
//        "test.adapter":           "test",         // test.adapter will adapt to testTracer below
//        "LEVEL.my.new.trace":     "Info",         // test tracer should have level info
//    }
//    trace2go.ConfigureRoot(conf, "LEVEL")         // root will spawn 'testTracer' children
//    tracer := tracing.Select("my.new.trace")      // now get tracer from factory
//    buf := &bytes.Buffer{}                        // log destination
//    tracer.SetOutput(buf)
//    msg := "this is a test info"
//    tracer.Infof(msg)        // this should log to my.new.trace at Info level
//    traceout := buf.String() // collect the output
/*
License

Governed by a 3-Clause BSD license. License file may be found in the root
folder of this module.

Copyright © Norbert Pillmayer <norbert@pillmayer.com>

*/
package trace2go
