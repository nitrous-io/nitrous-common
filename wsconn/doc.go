/*
wsconn implements a wrapper around an websocket connection. The wrapper implements
io.ReadWriteCloser so it can be used in place of functions that requires
io.Reader, or io.Writer such as github.com/nitrous-io/freighter/pkg/shell
*/
package wsconn
