package main

import (
	"io"
	"log"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const prefix string = ".."

func main() {
	godotenv.Load() // defaults to .env file
	DISCORD_TOKEN := os.Getenv("DISCORD_TOKEN")

	logFile := "bot.log"
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Write to Stdout and the log file
	w := io.MultiWriter(os.Stderr, file)

	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(w, handlerOpts))
	slog.SetDefault(logger)

	// Bot setup
	sess, err := discordgo.New("Bot " + DISCORD_TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		args := strings.Split(m.Content, " ")

		if m.Author.ID == s.State.User.ID {
			return
		}
		if args[0] != prefix {
			return
		}
		if strings.Contains(m.Content, "biden") {
			s.ChannelMessageSend(m.ChannelID, "blast")
		}
		if args[1] == "proverbs" {
			proverbs := []string{
				"Don't communicate by sharing memory, share memory by communicating.",
				"Concurrency is not parallelism.",
				"Channels orchestrate; mutexes serialize.",
				"The bigger the interface, the weaker the abstraction.",
				"Make the zero value useful.",
				"interface{} says nothing.",
				"Gofmt's style is no one's favorite, yet gofmt is everyone's favorite.",
				"A little copying is better than a little dependency.",
				"Syscall must always be guarded with build tags.",
				"Cgo must always be guarded with build tags.",
				"Cgo is not Go.",
				"With the unsafe package there are no guarantees.",
				"Clear is better than clever.",
				"Reflection is never clear.",
				"Errors are values.",
				"Don't just check errors, handle them gracefully.",
				"Design the architecture, name the components, document the details.",
				"Documentation is for users.",
				"Don't panic.",
			}

			selection := rand.Intn(len(proverbs))
			slog.Info("prov", "proverb", proverbs[selection])

			author := discordgo.MessageEmbedAuthor{
				Name: "Rob Pike",
				URL:  "https//go-proverbs.github.io",
			}
			embed := discordgo.MessageEmbed{
				Title:  proverbs[selection],
				Author: &author,
			}
			s.ChannelMessageSendEmbed(m.ChannelID, &embed)
			s.ChannelMessageSend(m.ChannelID, proverbs[selection])
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	slog.Info("App Token = " + DISCORD_TOKEN)
	slog.Info("Bot is online")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	sess.Close()
}
