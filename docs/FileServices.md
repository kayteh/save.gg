# File Services

Since we have to deal with files, lots of thought has to go into dealing with them.

## Primary: "The Black Hole"

**System:** Infinit Storage

This side is what we serve actual save packs with. Everything here is an XZipped file or tarball. For backing this, use the best tools for the job, and do not rely on one single service. If S3 and GCS are the best, use them. This system should be smart, and use local SSD caching wherever possible, especially for heavily-downloaded things.

Replication factor should be 3.
