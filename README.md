# Description

Terminal interface that applies code snippets to hardpoints in files.

Place [Hardpoints](#hardpoints) in your project files with [Tags](#tags) and
run `bolton` to preview and apply a [Bolton's](#boltons) contents to that
location.

## Installation

You must have [Go][1] installed and configured to install binaries into your
`$PATH`. Then simply run the following to install or ugrade the `bolton`
command:

```bash
env GO111MODULE=off go get -u github.com/Nekroze/bolton
```

## Usage

Simply execute the `bolton` command and it will display a tree view and a text panel.

The tree will show all discovered [Hardpoints](#hardpoints) in any files in the
directory in which the command was run, recursively. Underneath each
[Hardpoint](#hardpoints) is the [Tags](#tags) that can be applied to it. If
there are any [Boltons](#boltons) for any given tag they will be listed
underneath those [Tags](#tags) above them.

On the other side of the screen is a text panel that shows that contents of the
file around the currently selected tree entry. When a specific
[Bolton](#boltons) is selected, the tree view will display a preview of the
file it where applied.

```
Hardpoints / Boltons                                   #!/bin/sh
├──README.md:60
│  └──foo                                              # HARDPOINT: shell_function, global_env
├──README.md:67                                        export FOO=bar
│  └──foo
├──examples/lib.py:2
│  └──python_func
├──examples/lib.py:77
│  └──python_class
└──examples/script.sh:3
   ├──shell_function
   └──global_env
      ├──foo
      └──test
```

To apply the selected [Bolton](#boltons) press Enter, or CTRL+C to cancel.

### Boltons

It all starts with Boltons, snippets of text stored in files whose name and
location indicate what [Hardpoints](#hardpoints) they can be placed at. But
more on that later.

Bolton files are stored under directories beneath `~/.boltons` and the
directory they are in under that defines the [Tag](#tags) that bolton belongs
to. For example a file at `~/.boltons/foo/bar` is a Bolton called `bar` under
the `foo` [Tag](#tags).

The name and contents of these files can be anything and will be used verbatim
when applied to a [Hardpoint](#hardpoints).

### Hardpoints

Any line in a file that fits a specific pattern is a Hardpoint. Usually placed
in a comment, they declare a CSV of [Tags](#tags) to filter what
[Boltons](#boltons) can be placed on the line after the Hardpoint when applied.

To be discovered as a Hardpoint a line must end with a string that looks like
the following.

```yaml
HARDPOINT: foo
```

In a normal file it might look like this where it is placed in a comment:

```bash
#!/bin/sh
# HARDPOINT: foo
```

In the above example a Hardpoint is declared on line 2 with a CSV of
[Tags](#tags), in this case just one tag, `foo`.

When a [Bolton](#boltons) is applied to this Hardpoint, the rest of the file
will remain the same however the contents of the [Bolton](#boltons) will be
placed into the file on new lines on what would become line 3.

### Tags

Tags are a simple filtering mechanism, a [Hardpoint](#hardpoints) is declared
with a CSV of tag values that are used to filter down the [Boltons](#boltons)
that can be applied there by looking only in the directory for that Tag under
`~/.boltons`.
