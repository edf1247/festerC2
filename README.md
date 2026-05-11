## README

**FesterC2**: A command and control framework written in go.

### Capabilities:

Currently, Fester supports:
- Multiple Listeners/Sessions
- Implant Generation
- Session Management

### ToDos:

1. Currently, sessions are just shells. I'd like to add metasploit-esque sessions, where a user can execute commands outside of a standard reverse shell, but has the option to drop into a shell on the victim machine.
2. Support for pivoting
3. Better implants -> currently only support linux, and no attempt was made to avoid AV
4. Support for multiple teammembers on the same network