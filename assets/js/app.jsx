/** @jsx preact.h */

class App extends preact.Component {
  constructor() {
    super()

    this.state = {
      req: {},
      input: "",
    }

    window.clef.onRequest(this.handleClefRequest)
  }

  handleClefRequest = req => {
    this.setState({ req: req })
  }

  handleInputChange = e => {
    this.setState({ input: e.target.value })
  }

  handleInputSubmit = e => {
    e.preventDefault()

    clef.sendResponse({
      jsonrpc: "2.0",
      id:      this.state.req.id,
      result:  { text: this.state.input }
    })

    this.setState({
      req: {},
      input: "",
    })
  }

  render() {
    return (
      <div className="container">
        <h1>Ethereum Signer</h1>

        {this.state.req.method == "ui_onInputRequired" && (
          <InputForm
            params={this.state.req.params[0]}
            onChange={this.handleInputChange}
            onSubmit={this.handleInputSubmit}
          />
        )}
      </div>
    )
  }
}

const InputForm = ({ params, onChange, onSubmit }) => (
  <form>
    <div className="form-group">
      <label>
        {params.title}
      </label>

      <input
        type={params.isPassword ? "password" : "text"}
        className="form-control"
        onChange={onChange}
      />

      <small className="form-text">
        {params.prompt}
      </small>
    </div>

    <button type="submit" className="btn btn-primary" onClick={onSubmit}>
      Submit
    </button>
  </form>
)

preact.render(<App/>, document.body);
