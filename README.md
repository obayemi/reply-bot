# reply-bot
simple go bot to make automatic predefined replies to certain messages


## usage
clone, build (the usual)
make a discord bot user and get its token
create yaml db file

```sh
reply-bot -t [discord token] -i [db file]
```


## db file format:
```yaml
- trigger: string
  replies:  # replies to pick from
    - reply 1
    - reply 2
    - reply 3
  command: [mute|unmute]  # optional
  frequency: [0-100]  # optional
- ...
```
