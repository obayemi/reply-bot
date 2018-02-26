package main

import (
    "flag"
    "fmt"
    "math/rand"
    "os"
    "os/signal"
    "syscall"
    "gopkg.in/yaml.v2"
    "io/ioutil"

    "github.com/bwmarrin/discordgo"
)

type Trigger struct {
    Trigger string `yaml:"trigger"`
    Command string `yaml:"command"`
    Frequency uint `yaml:"frequency"`
    Replies []string `yaml:"replies"`
}

var (
    DiscordToken string
    DiscordSession, _ = discordgo.New()
    inFile string
    mute = false
    db = []Trigger{}
)

func stuffizer(s *discordgo.Session, m *discordgo.MessageCreate, trigger Trigger) bool {
    override := false
    switch trigger.Command {
    case "mute":
        mute = true
        override = true
    case "unmute":
        mute = false
    }

    return !mute || override
}

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID {
        return
    }
    for _, trigger := range db {
        if m.Content == trigger.Trigger {
            if stuffizer(s, m, trigger) {
                if trigger.Frequency == 0 || uint(rand.Intn(100)) < trigger.Frequency{
                    s.ChannelMessageSend(
                        m.ChannelID,
                        trigger.Replies[rand.Intn(len(trigger.Replies))],
                    )
                }
            }
        }
    }
}

func getDB(dbFile string) ([]Trigger, error) {
    db := []Trigger{}
    yamlFile, err := ioutil.ReadFile(dbFile)
    if err != nil {
        return db, fmt.Errorf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &db)
    if err != nil {
        return db, fmt.Errorf("Unmarshal: %v", err)
    }
    return db, nil
}

func initDiscord() error {
    var err error
    if DiscordSession.Token == "Bot " {
        return fmt.Errorf("You must provide a Discord authentication token.")
    }
    DiscordSession.State.User, err = DiscordSession.User("@me")
    if err != nil {
        return fmt.Errorf("error fetching user information, %s\n", err)
    }
    DiscordSession.AddHandler(MessageHandler)
    if err := DiscordSession.Open(); err != nil {
        return fmt.Errorf("error opening connection to Discord, %s\n", err)
    }
    return nil
}

func init() {
    var discordToken string
    //discordToken = os.Getenv("DISCORD_TOKEN")
    if discordToken == "" {
        flag.StringVar(&discordToken, "t", "", "Discord Authentication Token")
    }
    discordToken = os.Getenv("DB_FILE")
    flag.StringVar(&inFile, "i", "", "input file")
    flag.Parse()

    DiscordSession.Token = "Bot " + discordToken
}

func main() {
    var err error
    db, err = getDB(inFile)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(0)
    }
    if err := initDiscord(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(0)
    }
    defer DiscordSession.Close()

    fmt.Println(`Now running. Press CTRL-C to exit.`)
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc
}
