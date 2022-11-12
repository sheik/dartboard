# Fab

## What is it?
Fab is a tool for orchestrating builds and tests. It is built on top of
Go, allowing you to use the power of the Go language to help your development
loop.

## Get Started
First install the fab command:

    go install github.com/sheik/fab/cmd/fab@latest

Next, inside your project root directory, you can instantiate a fab file:

    fab init

This will create a file called **fab.go** which is excluded from your build
process. Once this is created you can edit it, and run it:

    fab

To see all of your actions:

    fab help
   
If you want to update fab, run:

    fab update

