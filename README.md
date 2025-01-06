# retro-sort

`retro-sort` sorts your files into a folder structure suitable for use with retro hardware

## Usage

```
retro-sort [src] [dst] [flags]
```

### Flags:

```
  -n, --dry-run       Dry run. Print the file names and exit
  -g, --glob string   Only include files matching this glob (default "*")
  -h, --help          help for retro-sort
  -q, --quiet         Don't print anything, just do it
  -s, --size int      Maximum number of directory entries (default 100)
  -u, --upper         Make upper-case directory names
```

`retro-sort` scans `src` for files (optionally matching the glob provided in the `--glob` flag).
It then determines an appropriate folder structure where no folder contains more than `size` files.
Finally, `retro-sort` copies all of the files to their new locations (unless you specify `--dry-run`).

retro-sort will exit with an error if `dst` already exists.
If you want to merge the files output from retro-sort, use a new folder first and then copy/move everything after.

For aesthetic reasons, you can specify `--upper` to have `retro-sort` create all the directory names in upper-case.

## Example

To give you an idea, here's a silly and trivial example. Imagine you have a folder of disk images that looks like this:

```
disks/aargh.disk
disks/addams_family.disk
disks/chaos_engine.disk
disks/elite.disk
disks/elite_ii.disk
disks/exolon.disk
```

And you are working with some retro hardware that simply can't cope if a folder contains any more than 2 files.

If you run `retro-sort -s 2 disks out`, then `retro-sort` will create the following directory structure:

```
disks/a/aargh.disk
disks/a/addams_family.disk
disks/c/chaos_engine.disk
disks/e/el/elite.disk
disks/e/el/elite_ii.disk
disks/e/ex/exolon.disk
```

If you picked a more generous size, `retro-sort` would make different decisions and consolidate where possible:

```
> retro-sort -s 4 disks out
disks/a-c/aargh.disk
disks/a-c/addams_family.disk
disks/a-c/chaos_engine.disk
disks/e/elite.disk
disks/e/elite_ii.disk
disks/e/exolon.disk
```
