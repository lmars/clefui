package clefui

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/url"

	"github.com/ethereum/go-ethereum/log"
	"github.com/zserge/webview"
)

//go:generate go-bindata -o assets.go -pkg clefui assets/...

func Run(ctx context.Context) error {
	log.Info("Starting Clef")
	clef, err := StartClef(ctx, "clef")
	if err != nil {
		log.Error("Error starting Clef", "err", err)
		return err
	}
	defer clef.Stop()

	log.Info("Initialising the UI")
	ui, err := NewUI(clef)
	if err != nil {
		return err
	}

	log.Info("Running the UI loop")
	go func() {
		<-ctx.Done()
		log.Info("Stopping the UI loop")
		ui.Exit()
	}()
	ui.Run()

	log.Info("UI loop stopped")
	return nil
}

const uiHTML = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">

    <script type="text/javascript">
      // Notify the backend when the DOM has loaded so it can inject
      // and render the app
      document.addEventListener("DOMContentLoaded", backend.injectApp);
    </script>
  </head>

  <body class="text-center">
  </body>
</html>
`

func NewUI(clef *Clef) (*UI, error) {
	// create a webview object
	w := webview.New(webview.Settings{
		Title: "Clef UI",
		URL:   "data:text/html," + url.PathEscape(uiHTML),
		Debug: true,
	})

	// expose the UI as the 'backend' JavaScript variable
	ui := &UI{WebView: w, clef: clef}
	if _, err := w.Bind("backend", ui); err != nil {
		return nil, err
	}

	// forward Clef requests to the JavaScript app
	go ui.ForwardClefRequests()

	return ui, nil
}

type UI struct {
	webview.WebView

	clef *Clef
}

// InjectApp is called within the initial HTML when the DOMContentLoaded event
// has fired, at which point we can inject and render the app.
func (ui *UI) InjectApp() {
	ui.Dispatch(func() {
		log.Info("Injecting CSS")
		ui.InjectCSS(string(MustAsset("assets/css/vendor/bootstrap.min.css")))
		ui.InjectCSS(string(MustAsset("assets/css/style.css")))

		log.Info("Injecting Babel and Preact")
		ui.Eval(string(MustAsset("assets/js/vendor/babel.min.js")))
		ui.Eval(string(MustAsset("assets/js/vendor/preact.min.js")))

		log.Info("Injecting app code")
		ui.Eval(string(MustAsset("assets/js/clef.js")))
		ui.Eval(fmt.Sprintf(`(function(){
			var n=document.createElement('script');
			n.setAttribute('type', 'text/babel');
			n.appendChild(document.createTextNode("%s"));
			document.body.appendChild(n);
			Babel.transformScriptTags();
		})()`, template.JSEscapeString(string(MustAsset("assets/js/app.jsx")))))
	})
}

// ForwardClefRequests reads JSON-RPC requests from Clef and forwards them to
// the JavaScript app using the 'window.clef' variable (see assets/js/clef.js).
func (ui *UI) ForwardClefRequests() {
	dec := json.NewDecoder(ui.clef.stdout)
	for {
		var msg json.RawMessage
		err := dec.Decode(&msg)
		if err == io.EOF {
			return
		} else if err != nil {
			log.Error("Error decoding Clef JSON-RPC request", "err", err)
			continue
		}
		log.Debug("Forwarding Clef request", "msg", string(msg))
		ui.Dispatch(func() {
			ui.Eval(fmt.Sprintf("window.clef.dispatchRequest(%s)", msg))
		})
	}
}

// SendClefResponse is called from the JavaScript app to respond to a Clef request.
func (ui *UI) SendClefResponse(msg json.RawMessage) {
	log.Debug("Forwarding Clef response", "msg", string(msg))
	ui.clef.stdin.Write(msg)
}
