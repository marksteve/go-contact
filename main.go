package main

import (
  "flag"
  "log"
  "net/http"
  "strings"

  "github.com/marksteve/go-mailgun"
)

var addr, from, to, mailgunKey string
var mgCli *mailgun.Client

type Mail struct {
  from      string
  to        []string
  cc        []string
  bcc       []string
  subject   string
  html      string
  text      string
  headers   map[string]string
  options   map[string]string
  variables map[string]string
}

func (m *Mail) From() string                 { return m.from }
func (m *Mail) To() []string                 { return m.to }
func (m *Mail) Cc() []string                 { return m.cc }
func (m *Mail) Bcc() []string                { return m.bcc }
func (m *Mail) Subject() string              { return m.subject }
func (m *Mail) Html() string                 { return m.html }
func (m *Mail) Text() string                 { return m.text }
func (m *Mail) Headers() map[string]string   { return m.headers }
func (m *Mail) Options() map[string]string   { return m.options }
func (m *Mail) Variables() map[string]string { return m.variables }

func submit(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST" {
    http.Error(w, "", http.StatusMethodNotAllowed)
    return
  }

  if err := r.ParseForm(); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  subject := r.PostFormValue("subject")
  if subject == "" {
    subject = "New message"
  }
  html := "<html><body>"
  text := ""
  for n, v := range r.PostForm {
    if n == "subject" {
      continue
    }
    html += "<p><strong>" + strings.Title(n) + ":</strong><br>" + v[0] + "</p>"
    text += strings.Title(n) + "\n" + v[0] + "\n\n"
  }
  html += "</body></html>"

  mail := &Mail{
    from:    from,
    to:      []string{to},
    subject: subject,
    html:    html,
    text:    text,
  }
  log.Print("sending...")
  log.Printf(" from: %s", from)
  log.Printf(" to: %s", to)
  log.Printf(" subject: %s", subject)
  msgId, err := mgCli.Send(mail)
  if err != nil {
    log.Print(err)
    http.Error(w, "Failed to send message", http.StatusInternalServerError)
    return
  }
  log.Printf("  msgId: %s", msgId)

  w.Write([]byte("Message sent"))

  return
}

func main() {
  log.Print("starting...")
  log.Printf(" addr: %s", addr)
  log.Printf(" from: %s", from)
  log.Printf(" to: %s", to)
  log.Printf(" mailgunKey: %s", mailgunKey)
  http.HandleFunc("/", submit)
  log.Fatal(http.ListenAndServe(addr, nil))
}

func init() {
  flag.StringVar(&addr, "addr", ":8000", "Address to listen to")
  flag.StringVar(&from, "from", "", "Email to send contact messages from")
  flag.StringVar(&to, "to", "", "Email to send contact messages to")
  flag.StringVar(&mailgunKey, "mailgunKey", "", "Mailgun key")
  flag.Parse()
  if from == "" {
    log.Fatal("From email is required.")
  }
  if to == "" {
    log.Fatal("To email is required.")
  }
  if mailgunKey == "" {
    log.Fatal("Mailgun API Key is required.")
  }
  mgCli = mailgun.New(mailgunKey)
}
