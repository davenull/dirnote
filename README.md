# dirnote - Notes for dirs

  Why keep notes about directories in the directories? I can't think of a good reason...

  So I created dirnotes, a simple CLI app written in go, powered by sqlite. Simply cd to your dir, and type `dn get` to see the notes you created previously. No need to remeber what dir you are in, dirnotes followr you around your computer. 

```
  dn help
/Users/birdie/.dirnote
An app to keep notes for your dirs

Usage:

    dirnote <command> [arguments...]

The commands are:

    add               adds a new note in the current dir
    del               deletes a note by global ID
    get               gets notes for the dir
    help              shows help message
    version           shows version of the application

Version: the best v0.0.1
```