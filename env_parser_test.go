package olayc

import (
	"reflect"
	"testing"
)

func TestEnvParser(t *testing.T) {
	psr := &envParser{}
	psr.parse([]string{
		"SHELL=/bin/zsh",
		"SHLVL=2",
		"SSH_AUTH_SOCK=/private/tmp/com.apple.launchd.guACPbxyU5/Listeners",
		"TERM=xterm-256color",
		"TERM_PROGRAM=tmux",
		"TERM_PROGRAM_VERSION=3.2",
		"TERM_SESSION_ID=w0t0p0:53538359-4F0F-4551-9088-091F2DDE2173",
		"TMPDIR=/var/folders/z_/kk3ybcfj4p5gl21bpwb5by8h0000gn/T/",
		"TMUX=/private/tmp/tmux-501/default,48207,0",
		"TMUX_PANE=%1",
		"XPC_FLAGS=0x0",
		"XPC_SERVICE_NAME=0",
		"_P9K_SSH_TTY=/dev/ttys002",
		"_P9K_TTY=/dev/ttys002",
		"__CFBundleIdentifier=com.googlecode.iterm2",
		"__CF_USER_TEXT_ENCODING=0x0:25:52",
		"P9K_SSH=0",
		"_=/usr/bin/env",
	})
	var got = psr.kvs
	var expect = []KV{
		{"shell", "/bin/zsh"},
		{"shlvl", uint64(2)},
		{"ssh.auth.sock", "/private/tmp/com.apple.launchd.guACPbxyU5/Listeners"},
		{"term", "xterm-256color"},
		{"term.program", "tmux"},
		{"term.program.version", float64(3.2)},
		{"term.session.id", "w0t0p0:53538359-4F0F-4551-9088-091F2DDE2173"},
		{"tmpdir", "/var/folders/z_/kk3ybcfj4p5gl21bpwb5by8h0000gn/T/"},
		{"tmux", "/private/tmp/tmux-501/default,48207,0"},
		{"tmux.pane", "%1"},
		{"xpc.flags", "0x0"},
		{"xpc.service.name", uint64(0)},
		{"p9k.ssh.tty", "/dev/ttys002"},
		{"p9k.tty", "/dev/ttys002"},
		{"cfbundleidentifier", "com.googlecode.iterm2"},
		{"cf.user.text.encoding", "0x0:25:52"},
		{"p9k.ssh", uint64(0)},
	}
	if len(expect) != len(got) {
		t.Logf("expect(len %v) != got(len %v)\n", len(expect), len(got))
		t.Logf("expect: %v\n", expect)
		t.Logf("got: %v\n", got)
		t.FailNow()
	}

	for i, kvExpect := range expect {
		kvGot := got[i]
		if kvExpect.key != kvGot.key || kvExpect.value != kvGot.value {
			t.Errorf("[%v] expect(%v:%v) != got(%v:%v)", i, expect[i], reflect.TypeOf(expect[i].value),
				got[i], reflect.TypeOf(got[i].value))
		}
	}
}
