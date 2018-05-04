package pjlink

import (
	"bufio"
	"errors"
	"net"
	"strings"
	"time"
	"log"
)

const pjLinkPort = "4352"

type PJProjector struct {
	Address  string
	Port     string
	Password string
}

func NewProjector(IP string, password string) *PJProjector {
	return &PJProjector{
		Address:  IP,
		Port:     pjLinkPort,
		Password: password,
	}
}

//--------------------------------------------------------------------------------------------------------------------//
//--------------- Functional Calls -----------------------------------------------------------------------------------//
//--------------------------------------------------------------------------------------------------------------------//

//--------------- Power ----------------------------------------------------------------------------------------------//
func (pr *PJProjector) GetPowerStatus() (*PJResponse, error) {
	req := PJRequest{
		Class:     1,
		Command:   "POWR",
		Parameter: "?",
	}
	return pr.SendRequest(req)
}

func (pr *PJProjector) TurnOn() (*PJResponse, error) {
	req := PJRequest{
		Class:     1,
		Command:   "POWR",
		Parameter: "1",
	}
	return pr.SendRequest(req)
}

func (pr *PJProjector) TurnOff() (*PJResponse, error) {
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
func (pr *PJProjector) SendRequest(request PJRequest) (*PJResponse, error) {
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

func (pr *PJProjector) sendRawRequest(request PJRequest) (*PJResponse, error) {
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

	seed := pr.checkAuthentication(challenge)
	stringCommand := request.toRaw(seed, pr.Password)


	//send command
	connection.Write([]byte(stringCommand))
	scanner.Scan() //grab response line

	resp := NewPJResponse()
	err := resp.Parse(scanner.Text())
	if err != nil {
		log.Println("Auth error.")
		return resp, err
	}

	return resp, nil
}

//attempts to establish a TCP socket with the specified IP:port
//success: returns populated pjlinkConn struct and nil error
//failure: returns empty pjlinkConn and error
func (pr *PJProjector) connectToPJLink() (net.Conn, error) {
	protocol := "tcp" //PJLink always uses TCP
	timeout := 10     //represents seconds

	connection, connectionError := net.DialTimeout(protocol, net.JoinHostPort(pr.Address, pr.Port), time.Duration(timeout)*time.Second)
	if connectionError != nil {
		return connection, errors.New("failed to establish a connection with " +
			"pjlink device. error msg: " + connectionError.Error())
	}
	return connection, connectionError
}

// check if this Projector uses authentication. If so return true and the given seed. Otherwise false and an empty string.
func (pr *PJProjector) checkAuthentication(response []string) (seed string) {
	if response[0] != "PJLINK" {
		return ""
	}
	if response[1] == "0" {
		return ""
	} else if response[1] == "1" {
		return response[2]
	}

	return ""
}
