/*
Taken from https://github.com/dotcloud/docker/blob/master/utils/stdcopy.go#d17c92965937bbc3181b73bb8a4e21ecee217e5f

This is helpful when we want to parse log stream from docker as docker
multiplexes stdout/stderr logs into a single stream
*/
package stdcopy
