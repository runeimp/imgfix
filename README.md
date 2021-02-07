ImageFix
========

A command line tool to handle simple image fixes like renaming with the correct file extension


Features
--------

* [x] Fix the file extension of mislabeled images
* [ ] Specify a path instead of using a list of files
* [ ] Recurse a path to fix all images in a directory tree
* [ ] Specify max depth to recurse
* [ ] Fix the files modification date based on the image meta-data?
* [ ] ...?


Examples
--------

### Basic Help

```text
$ imgfix -help
Usage: imgfix [OPTIONS] [file1 file2 ...]

OPTIONS:
  -dry-run  Do not modify files (dry-run) (default: false)
  -help   Display this help info (default: false)
  -verbose  Display verbose output (default: false)
  -version  Display app version (default: false)

```

### A Dry-Run

```text
$ imgfix -dry-run is-dir/* is-*
Renaming "is-dir/is-png" to "is-dir/is-png.png" (dry-run)
Renaming "is-gif-animated.jpg" to "is-gif-animated.gif" (dry-run)
Not renaming "is-jpeg-2.png" to "is-jpeg-2.jpg" as the later already exists (dry-run)
Renaming "is-png" to "is-png.png" (dry-run)
Renaming "is-png.gif" to "is-png.png" (dry-run) <--- BUG!
```

### Verbose Output

```text
$ imgfix -dry-run -verbose is-dir/* is-*
Renaming "is-dir/is-png" to "is-dir/is-png.png" (dry-run)
Skipping "is-dir" (it is not an image)
Skipping "is-empty-text.jpeg" (it is not an image)
Renaming "is-gif-animated.jpg" to "is-gif-animated.gif" (dry-run)
Skipping "is-jpeg-1.jpeg" (no fixing needed) (dry-run)
Skipping "is-jpeg-2.jpg" (no fixing needed) (dry-run)
Not renaming "is-jpeg-2.png" to "is-jpeg-2.jpg" as the later already exists (dry-run)
Renaming "is-png" to "is-png.png" (dry-run)
Renaming "is-png.gif" to "is-png.png" (dry-run) <--- BUG!
```