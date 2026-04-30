package ftp_helper

import (
	"errors"
	"fmt"
	"github.com/jlaffaye/ftp"
	"io"
	"strings"
	"time"
)

var (
	FileOrDirectoryNotExistError = errors.New("no such file or directory")
	RelativePathError            = errors.New("only absolute path is allowed")
	InvalidCredentialError       = errors.New("invalid credential submitted for the FTP server")
	ConnectionError              = errors.New("can not connect to the FTP server")
)

type FTPConfig struct {
	Address  string
	Port     int
	Timeout  time.Duration
	Username string
	Password string
}

/*
FTP represents a connection to an FTP server with additional functionality for interacting with the server, such as storing files, creating directories, and handling connections.

Attributes:
- ServerConn: *ftp.ServerConn - The underlying connection to the FTP server.
- dir: []string - Represents the current directory path on the FTP server.

Methods:
- Connect: Establishes a connection to the FTP server using the provided configuration.
- Disconnect: Closes the FTP session.
- StoreFile: Stores a file on the remote FTP server.
- DirectoryWalk: Creates directories on the FTP server based on the provided path.
- SafeStoreFile: Safely stores a file on the remote FTP server after checking for path existence.

Remember to call Disconnect to close the FTP session properly.
*/
type FTP struct {
	*ftp.ServerConn
	dir []string
}

/*
Connect establishes a connection to an FTP server using the provided FTPConfig.
It creates a new FTP instance, dials the FTP server with the specified address and port,
and attempts to log in using the provided username and password.

Parameters:
- conf: FTPConfig struct containing the FTP server details such as address, port, timeout, username, and password.

Returns:
- *FTP: A pointer to the FTP instance if the connection and login are successful.
- error: An error if any issues occur during the connection process.

Possible error values:
- ConnectionError: Indicates a failure to connect to the FTP server.
- InvalidCredentialError: Indicates that the provided credentials are invalid for the FTP server.

Note: The timeout for the connection can be specified in the FTPConfig. If not provided, a default timeout of 2 seconds is used.
*/
func Connect(conf FTPConfig) (*FTP, error) {
	t := 2 * time.Second
	if conf.Timeout != 0 {
		t = conf.Timeout
	}

	conn, err := ftp.Dial(
		fmt.Sprintf("%s:%d", conf.Address, conf.Port),
		ftp.DialWithTimeout(t),
	)
	if err != nil {
		return nil, ConnectionError
	}

	instance := &FTP{ServerConn: conn}

	if err = instance.login(conf.Username, conf.Password); err != nil {
		return nil, InvalidCredentialError
	}

	return instance, nil
}

/*
Disconnect closes the connection to the FTP server by sending a QUIT command.
It returns an error if there was a problem disconnecting from the server.
*/
func (f *FTP) Disconnect() error {
	return f.Quit()
}

/*
login method authenticates the user with the provided username and password.
It calls the Login method of the embedded ServerConn field to perform the login operation.

Parameters:
- username: a string representing the username of the user
- password: a string representing the password of the user

Returns:
- error: an error if the login operation fails, nil otherwise
*/
func (f *FTP) login(username, password string) error {
	return f.ServerConn.Login(username, password)
}

/*
StoreFile stores a file in the FTP server at the specified path.

Parameters:
- path: a string representing the path where the file will be stored.
- file: an io.Reader interface representing the content of the file to be stored.

Returns:
- error: an error indicating any issues that occurred during the file storage process.
*/
func (f *FTP) StoreFile(path string, file io.Reader) error {
	err := f.ServerConn.Stor(path, file)
	if err != nil && err.Error() == "Can't open that file: No such file or directory" {
		return FileOrDirectoryNotExistError
	}
	return nil
}

/*
DirectoryWalk walks through the directory structure based on the provided path.
It trims the leading '/' from the path and splits it into individual directory names.
Then, it calls the internal walk method to navigate through the directories recursively.

Parameters:
- path: a string representing the path to walk through

Returns:
- error: an error if any occurred during the directory walk process
*/
func (f *FTP) DirectoryWalk(path string) error {
	path = strings.TrimPrefix(path, "/")
	paths := strings.Split(path, "/")

	return f.walk(paths)
}

/*
walk recursively creates directories based on the provided paths.

Parameters:
- paths: a slice of strings representing the directory structure to be created.

Returns:
- error: an error if any occurred during the directory creation process.

Behavior:
- If the first path element is empty, the function returns nil.
- The function checks if the directory specified by the paths exists. If not, it creates the directory.
- It recursively creates directories for each path element in the input paths slice.
*/
func (f *FTP) walk(paths []string) error {
	if paths[0] == "" {
		return nil
	}

	absolutePath := "/"
	if f.dir != nil {
		absolutePath = "/" + strings.Join(f.dir, "/")
	}

	dirs, err := f.ServerConn.List(absolutePath)
	if err != nil && err.Error() != "550 Can't check for file existence" {
		return err
	}

	isPathExist := false
	for _, dir := range dirs {
		if dir.Name == paths[0] {
			isPathExist = true
		}
	}

	if !isPathExist {
		if err := f.ServerConn.MakeDir(fmt.Sprintf("%s/%s", absolutePath, paths[0])); err != nil {
			return err
		}
	}

	if len(paths[1:]) > 0 {
		f.dir = append(f.dir, paths[0])
		return f.walk(paths[1:])
	}

	return nil
}

/*
SafeStoreFile stores a file in the FTP server at the specified path.
It first checks if the path contains any relative path indicators.
Then it splits the path into individual directories and walks through them to ensure they exist.
Finally, it stores the file in the FTP server at the specified path.

Parameters:
- path: a string representing the path where the file should be stored.
- file: an io.Reader interface for reading the file content.

Returns:
- error: an error if any issue occurs during the process, such as invalid path, directory creation failure, or file storage error.
*/
func (f *FTP) SafeStoreFile(path string, file io.Reader) error {
	if strings.Contains("~", path) {
		return RelativePathError
	}

	paths := strings.Split(strings.Trim(path, "/"), "/")

	if err := f.walk(paths[:len(paths)-1]); err != nil {
		return err
	}

	if err := f.ServerConn.Stor(path, file); err != nil {
		return err
	}

	return nil
}
