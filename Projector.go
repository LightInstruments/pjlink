package pjlink

import (
	"bufio"
	"errors"
	"net"
	"strings"
	"time"
)

const pjLinkPort = "4352"

type Projector struct {
	Address    string
	Port       string
	Password   string
}

func NewProjector(IP string, password string) *Projector {
	return &Projector{
		Address:  IP,
		Port:     pjLinkPort,
		Password: password,
	}
}

//--------------------------------------------------------------------------------------------------------------------//
//--------------- Functional Calls -----------------------------------------------------------------------------------//
//--------------------------------------------------------------------------------------------------------------------//

//--------------- Power ----------------------------------------------------------------------------------------------//
func (pr *Projector) GetPowerStatus() (*PJResponse, error) {
	req := PJRequest{
		Class:     1,
		Command:   "POWR",
		Parameter: "?",
	}
	return pr.SendRequest(req)
}

func (pr *Projector) TurnOn() (*PJResponse, error) {
	req := PJRequest{
		Class:     1,
		Command:   "POWR",
		Parameter: "1",
	}
	return pr.SendRequest(req)
}

func (pr *Projector) TurnOff() (*PJResponse, error) {
	req := PJRequest{
		Class:     1,
		Command:   "POWR",
		Parameter: "0",
	}
	return pr.SendRequest(req)
}

//--------------------------------------------------------------------------------------------------------------------//
// Low-Level Calls
//--------------------------------------------------------------------------------------------------------------------//
func (pr *Projector) SendRequest(request PJRequest) (*PJResponse, error) {
	if err := request.Validate(); err != nil { //malformed command, don't send
		return nil, err
	} else { //send request and parse response into struct
		response, requestError := pr.sendRawRequest(request)
		if requestError != nil {
			return nil, requestError
		} else {
			return response, nil
		}
	}
}

func (pr *Projector) sendRawRequest(request PJRequest) (*PJResponse, error) {
	//establish TCP connection with PJLink device
	connection, connectionError := pr.connectToPJLink()
	defer connection.Close()

	if connectionError != nil {
		return nil, connectionError
	}

	// Define a split function that separates on carriage return (i.e '\r').
	onCarriageReturn := func(data []byte, atEOF bool) (advance int, token []byte,
		err error) {
		for i := 0; i < len(data); i++ {
			if data[i] == '\r' {
				return i + 1, data[:i], nil
			}
		}
		// There is one final token to be delivered, which may be the empty string.
		// Returning bufio.ErrFinalToken here tells Scan there are no more tokens
		// after this but does not trigger an error to be returned from Scan itself.
		return 0, data, bufio.ErrFinalToken
	}

	//setup scanner
	scanner := bufio.NewScanner(connection)
	scanner.Split(onCarriageReturn)
	scanner.Scan() //grab a line
	challenge := strings.Split(scanner.Text(), " ")

	//verify PJLink and correct class
	if !pr.verifyPJLink(challenge) {
		// TODO: Handle not PJLink class 1
		return nil, errors.New("Not a PJLINK class 1 connection")
	}
	seed := challenge[2]
	stringCommand := request.toRaw(seed, pr.Password)

	//send command
	connection.Write([]byte(stringCommand + "\r"))
	scanner.Scan() //grab response line

	resp := NewPJResponse()
	resp.Parse(scanner.Text())

	return resp, nil
}

//attempts to establish a TCP socket with the specified IP:port
//success: returns populated pjlinkConn struct and nil error
//failure: returns empty pjlinkConn and error
func (pr *Projector) connectToPJLink() (net.Conn, error) {
	protocol := "tcp" //PJLink always uses TCP
	timeout := 10      //represents seconds

	connection, connectionError := net.DialTimeout(protocol, net.JoinHostPort(pr.Address, pr.Port), time.Duration(timeout)*time.Second)
	if connectionError != nil {
		return connection, errors.New("failed to establish a connection with " +
			"pjlink device. error msg: " + connectionError.Error())
	}
	return connection, connectionError
}

//verify we receive a pjlink class 1 challenge
//success: returns true
//failure: returns false
func (pr *Projector) verifyPJLink(response []string) bool {
	if response[0] != "PJLINK" {
		return false
	}

	if response[1] != "1" {
		return false
	}

	return true
}
