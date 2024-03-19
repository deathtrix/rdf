# RDF - Rolling hash diff

### Description
The project uses Rabin-Karp rolling hash for file diffing

Output:
```
Operation   sonf:eonf       soof:eoof
O           256:512         256:512
```

- sonf - start offset newfile
- eonf - end offset newfile
- soof - start offset oldfile
- eoof - end offset oldfile

Operation type:
- O - original position block
- M - moved block
- C - changed block

### Usage
```
make build
./rdf <block-size> <oldfilename> <newfilename>
```

### Testing
```
make test
make bench
```

### Assumptions
- "Should be able to recognize changes between chunks. Only the exact differing
locations should be added to the delta." from the Requirements was not clear to me.
If you only save the locations, and not the changed data you won't be able to rebuild
the new file from the old one and delta. As I understood from the above that changing
data should not be saved in delta, so I did not add it but I can implement it if needed.

- "detects chunk removals and detects additions between chunks with shifted original chunks"
was not clear. Was I supposed to detect them like changes or detect them as inserts and deletes?
The current implementation only detects them as changes. I considered it out of scope for this
project, but I can implement it if needed.

- also, I implemented the signature building and delta building in the same cli command
as it was not clear from the specs in they should be separated. To separate them you
only need to save the Signature to a file using encoding/gob and maybe use cobra for cli.

