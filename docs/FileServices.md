# File Services

Since we have to deal with files, lots of thought has to go into dealing with them. The services are in layers, and are in different locations per usage.

## Ingress

Ingress layer is a regional system that only is accessible to the current user. Ingress systems are awaiting processing and are not to be committed to any primary file storage medium like S3/GCS/B3. Things from here can be deleted.
