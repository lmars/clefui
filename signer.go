package signer

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/url"

	"github.com/ethereum/go-ethereum/log"
	"github.com/zserge/webview"
)

//go:generate go-bindata -o assets.go -pkg signer assets/...

var ErrShutdown = errors.New("signer: Signer stopped")

func New() *Signer {
	return &Signer{}
}

type Signer struct {
}

type App struct {
	w webview.WebView
}

func (a *App) Init() {
	a.w.Dispatch(func() {
		log.Info("Injecting CSS")
		a.w.InjectCSS(string(MustAsset("assets/css/vendor/bootstrap.min.css")))
		a.w.InjectCSS(string(MustAsset("assets/css/style.css")))

		log.Info("Injecting Babel and Preact")
		a.w.Eval(string(MustAsset("assets/js/vendor/babel.min.js")))
		a.w.Eval(string(MustAsset("assets/js/vendor/preact.min.js")))

		log.Info("Injecting app code")
		a.w.Eval(fmt.Sprintf(`(function(){
			var n=document.createElement('script');
			n.setAttribute('type', 'text/babel');
			n.appendChild(document.createTextNode("%s"));
			document.body.appendChild(n);
			Babel.transformScriptTags();
		})()`, template.JSEscapeString(string(MustAsset("assets/js/app.jsx")))))
	})
}

func (s *Signer) Run(ctx context.Context) error {
	log.Info("Starting the signer UI")
	w := webview.New(webview.Settings{
		Title: "Ethereum Signer",
		URL:   "data:text/html," + url.PathEscape(string(MustAsset("assets/index.html"))),
		Debug: true,
	})
	go func() {
		<-ctx.Done()
		log.Info("Stopping the signer UI")
		w.Exit()
	}()

	w.Bind("app", &App{w})

	w.Run()

	log.Info("Signer UI exited")
	return nil
}
