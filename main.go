package main
import (
    "fmt"
    "github.com/bwmarrin/discordgo"
    "github.com/zymtom/argconf"
    //"time"
)
var maps map[string]chan message
var trash_amount int
type message struct {
    ID  string
    ChannelID string   
}
func main() {
    paramMap := map[string][]string{
        "username":[]string{"string", "", "Username for your discord user, which is your email"},
        "password":[]string{"string", "", "Password for your discord user"},
        "trash-amount":[]string{"int", "50", "How many messages to keep"},
    }
    values, err := argconf.HandleParams(paramMap)
    if values["username"].(string) == "" || values["password"].(string) == "" {
        fmt.Println("You need to provide a username and password. Do -h for arguments.")
        return 
    }
    trash_amount = values["trash-amount"].(int)
    discord, err := discordgo.New(values["username"].(string), values["password"].(string))
    if err != nil {
		fmt.Println(err)
		return
	}
    discord.AddHandler(messageCreate)
    discord.Open()
    maps = make(map[string]chan message)
    for {
        for _, v := range maps {
            if len(v) > trash_amount {
                for obj := range v {
                    discord.ChannelMessageDelete(obj.ChannelID, obj.ID)
                    if !(len(v) > trash_amount) {
                        break
                    }
                }
            }
        }
        
    }
    var input string
    fmt.Scanln(&input)
    return
}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    /*fmt.Printf("%s %20s %20s %20s > %s\n", m.Author.ID, m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)
    fmt.Println(len(maps[m.ChannelID]))*/
    if maps[m.ChannelID] == nil {
        maps[m.ChannelID] = make(chan message, trash_amount+1)
    }
    maps[m.ChannelID] <- message{ID:m.ID, ChannelID:m.ChannelID}
}

