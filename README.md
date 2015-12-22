Flagspy provides access to command-line flags before Go's flag.Parse() method
has been called. It does similar parsing as the standard flag package, but
treats all flags as type string and does not require any pre-registering of flag
names.

The intent of this package is to give initialization methods early access to
flags for low-level things like logging.

# Main features

 * Access to flag values during initialization.
 * Parses flag values on first access, caching them for faster use.
 * A `Free` method to delete the cache once initialization is complete.

# Basic usage

The intended use case is as an augmentation to the standard flag package:

    package preinit

    import (
      "flag"

      "github.com/hegh/flagspy"
    )

    func init() {
      // This is for documentation.
      _ = flag.String("logdir", ".", "The directory where log files will be written.")
      logdir, ok = flagspy.Get("logdir")
      if !ok {
        logdir = "."
      }

      // Set up logging to logdir.
      // ...
    }

Then, in package main:

    package main

    import "github.com/hegh/flagspy"

    func main() {
      // Flagspy values are not intended for use outside of initialization.
      // Either save the value during initialization, or use a standard flag.
      // You could do this in init(), instead of the beginning of main.
      flagspy.Free()
    }
