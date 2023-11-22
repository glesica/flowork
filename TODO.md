# To Do List

  1. Implement an SSH runner that takes a list of machines and
     multiplexes the task instances over them. This will also
     provide a minimal way to run in cloud environments since
     they all support groups of machines and SSH access.
  2. Solve the general file copy problem where we want to copy
     a file from some arbitrary environment (like S3) to some
     other environment (like an SSH machine) but, when possible,
     we don't want to bring the file through the orchestrator
     machine (local or otherwise). We could do this by running
     a Flowork container on the SSH machine (for the case above)
     and pulling the file down directly. But that's also pretty
     complicated since we'd need to handle permissions. Maybe,
     as a first effort, just refuse to do certain transfers?

Number (2) above is probably going to be the way to go. Turns
out that moving files over SSH is kind of annoying and requires
a particular implementation of scp (that doesn't seem to exist
on Macs?). This way the only requirement remains Docker.

Add a TaskError type that has slots for the task ID, volume,
and whatever other info.
