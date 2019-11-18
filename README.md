Ignition Spec v1-v2.x.0 to v3.0.0 Config Converter
===================================================

## What is this?

This is a tool and library for convert old (v1-v2.x.0) Ignition configs to the
new v3.0.0 format.

## Why is this not part of Ignition?

The old spec versions have bugs that allow specifying configs that don't make
sense. For example, it is valid for a v2.1+ config to specify the same path
should be both a directory and a file. The behavior there is defined by
Ignition's implementation instead of the spec and in certain edge cases, by the
contents of the filesystem Ignition is operating on.

This means Ignition can't be guaranteed to automatically translate an old
config to an equivalent new config; it can fail at conversion. Since Ignition
internally translates old configs to the latest config, this would mean old
Ignition configs could stop working on newer versions of whatever OS included
Ignition. Additionally, due to the change in how filesystems are handled (new
configs require specifying the path relative to the sysroot that Ignition
should mount the filesystem at), some configs require extra information to
convert from the old versions to the new versions.

This tool exists to allow _mechcanical_ translation of old configs to new
configs. If you are also switching operating systems, other changes may be
necessary.

## How can I ensure my old config is translatable?

Most of the problems in old configs stem from specifying illogical things. Make
sure you don't have any duplicate entries. Do not rely on the order in which
files, directories, or links are created. Most configs should be translatable
without problems.
