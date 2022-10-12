# dirnote - Notes for dirs

  Why keep notes about directories in the directories? I can't think of a good reason...

  So I created dirnotes, a simple CLI app written in go, powered by sqlite. Simply cd to your dir, and type `dn get` to see the notes you created previously. No need to remeber what dir you are in, dirnotes followr you around your computer.

## Install
```
wget https://github.com/davenull/dirnote/releases/download/dev0.0.2/dn-osx-ARM64
mv dn-osx-ARM64 dn
chmod +x
mv dn ~/bin/
```
Add the following to your .bashrc or .zshrc

```
export PATH="$HOME/bin:$PATH"
```

## Usage
To initialize a dir, add a note!!
```
dn add "Hello dirnote!"
```
On first run, dirnote will set itself up and create the database, then you just use it!

## Backing Up dirnote DB

As the DB is only open when you run the command, backing up the db file to a cloud drive is simple.
The only thing you need to save is `~/.dirnote/dirnotes.sqlite` and you can rest assured your notes are safe!



## CLI help output

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