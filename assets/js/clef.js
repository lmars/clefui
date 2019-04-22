// the clef object is used to communicate with Clef via the UI backend
window.clef = {

  // onRequest registers a callback that will be called when Clef sends a request
  onRequest: function(cb) {
    this.cb = cb
  },

  // dispatchRequest is called by the backend when Clef sends a request
  dispatchRequest: function(req) {
    if (this.cb) {
      this.cb(req)
    }
  },

  // sendResponse sends a response to Clef via the backend
  sendResponse: function(res) {
    backend.sendClefResponse(res)
  }
}
