package voicy

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/XzFrosT/crunch/utils"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/voice"
	"github.com/diamondburned/arikawa/v3/voice/voicegateway"
	"github.com/diamondburned/oggreader"
	"github.com/pkg/errors"
)

const (
	stoppedState = iota
	changingState
	pausedState
	playingState
)

type Session struct {
	Session *voice.Session

	source string
	isOpus bool

	Position time.Duration

	state   int
	channel chan int

	context context.Context
	cancel  context.CancelFunc
}

func New(state *state.State, guildID discord.GuildID, channelID discord.ChannelID) (*Session, error) {
	session, err := voice.NewSession(state)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create a voice session")
	}

	if err := session.JoinChannel(guildID, channelID, false, true); err != nil {
		return nil, errors.Wrap(err, "unable to connect to voice channel")
	}

	return &Session{Session: session}, nil
}

func (vs *Session) PlayURL(source string, isOpus bool) error {
	if vs.state != stoppedState && vs.state != changingState {
		vs.Stop()
	}

	vs.context, vs.cancel = context.WithCancel(context.Background())
	vs.source, vs.isOpus = source, isOpus

	ffmpeg := exec.CommandContext(vs.context, "ffmpeg",
		"-loglevel", "error", "-reconnect", "1", "-reconnect_streamed", "1", "-reconnect_delay_max", "5", "-ss", utils.FormatTime(vs.Position),
		"-i", source, "-vn", "-codec", utils.Is(vs.isOpus, "copy", "libopus"), "-vbr", "off", "-frame_duration", "20", "-f", "opus", "-")

	stdout, err := ffmpeg.StdoutPipe()
	if err != nil {
		vs.stop()
		return errors.Wrapf(err, "failed to get ffmpeg stdout")
	}

	var stderr bytes.Buffer
	ffmpeg.Stderr = &stderr

	if err := ffmpeg.Start(); err != nil {
		vs.stop()
		return errors.Wrapf(err, "failed to start ffmpeg process")
	}

	if err := vs.SendSpeaking(); err != nil {
		vs.stop()
		return errors.Wrapf(err, "failed to send speaking packet to discord")
	}

	vs.setState(playingState)

	if err := oggreader.DecodeBuffered(vs, stdout); err != nil && vs.state != changingState {
		vs.stop()
		return err
	}

	if err, std := ffmpeg.Wait(), stderr.String(); err != nil && std != "" {
		vs.stop()
		return errors.Wrapf(errors.New(strings.ReplaceAll(std, vs.source, "source")), "ffmpeg returned error")
	}

	if vs.state == changingState {
		return vs.PlayURL(vs.source, vs.isOpus)
	}

	vs.stop()
	return nil
}

func (vs *Session) Destroy() {
	vs.Stop()
	vs.Session.Leave()
}

func (vs *Session) Seek(position time.Duration) {
	if vs.state == stoppedState {
		return
	}
	vs.Position = position

	vs.setState(changingState)
	vs.Stop()
}

func (vs *Session) Resume() {
	if vs.state == pausedState {
		vs.setState(playingState)
		vs.SendSpeaking()
	}
}

func (vs *Session) Pause() {
	if vs.state != stoppedState && vs.state != changingState {
		vs.setState(pausedState)
	}
}

func (vs *Session) Stop() {
	if vs.state == stoppedState {
		return
	}

	if vs.cancel != nil {
		vs.cancel()
		vs.waitState(stoppedState, playingState, changingState, pausedState)
	}
}

func (vs *Session) stop() {
	vs.cancel()
	vs.setState(stoppedState)
	vs.Position = 0
}

func (vs *Session) SendSpeaking() error {
	return vs.Session.Speaking(voicegateway.Microphone)
}

func (vs *Session) Write(data []byte) (int, error) {
	if vs.state == pausedState {
		vs.waitState(playingState, stoppedState, changingState)
	}

	vs.Position = vs.Position + (20 * time.Millisecond)
	return vs.Session.WriteCtx(vs.context, data)
}

func (vs *Session) setState(state int) {
	vs.state = state
	if vs.channel != nil {
		vs.channel <- state
	}
}

func (vs *Session) waitState(states ...int) {
	vs.channel = make(chan int)

	for {
		if newState := <-vs.channel; utils.IntegerArrayContains(states, newState) {
			close(vs.channel)
			vs.channel = nil
			break
		}
	}
}
